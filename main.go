// an extremely fast, zero config server for modern web applications.
package main

import (
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"syscall"

	"github.com/evanw/esbuild/pkg/api"
	"github.com/evanw/esbuild/pkg/cli"
	"github.com/osdevisnot/sorvor/pkg/authority"
	"github.com/osdevisnot/sorvor/pkg/livereload"
	"github.com/osdevisnot/sorvor/pkg/logger"
	"github.com/osdevisnot/sorvor/pkg/pkgjson"
)

type sorvor struct {
	BuildOptions api.BuildOptions
	Entry        string
	Host         string
	Port         string
	Serve        bool
	Secure       bool
}

func readOptions(pkg *pkgjson.PkgJSON) *sorvor {
	var err error
	var esbuildArgs []string

	osArgs := os.Args[1:]
	serv := &sorvor{}

	for _, arg := range osArgs {
		switch {
		case strings.HasPrefix(arg, "--host"):
			serv.Host = arg[len("--host="):]
		case strings.HasPrefix(arg, "--port"):
			port, err := strconv.Atoi(arg[len("--port="):])
			logger.Fatal(err, "invalid port value")
			serv.Port = ":" + strconv.Itoa(port)
		case arg == "--serve":
			serv.Serve = true
		case arg == "--secure":
			serv.Secure = true
		case !strings.HasPrefix(arg, "--"):
			serv.Entry = arg
		default:
			esbuildArgs = append(esbuildArgs, arg)
		}
	}

	serv.BuildOptions, err = cli.ParseBuildOptions(esbuildArgs)
	logger.Fatal(err, "Invalid option for esbuild")

	serv.BuildOptions.Bundle = true
	serv.BuildOptions.Write = true

	if serv.Serve == true {
		if serv.Port == "" {
			serv.Port = ":1234"
		}
		serv.BuildOptions.Define = map[string]string{"process.env.NODE_ENV": "'development'"}
	} else {
		serv.BuildOptions.Define = map[string]string{"process.env.NODE_ENV": "'production'"}
	}
	if serv.BuildOptions.Outdir == "" {
		serv.BuildOptions.Outdir = "dist"
	}
	if serv.BuildOptions.Format == api.FormatDefault {
		serv.BuildOptions.Format = api.FormatESModule
	}
	if serv.Entry == "" {
		serv.Entry = "public/index.html"
	}
	if serv.Host == "" {
		serv.Host = "localhost"
	}
	return serv
}

func (serv *sorvor) esbuild(entry string) (string, api.BuildResult) {
	serv.BuildOptions.EntryPoints = []string{entry}
	result := api.Build(serv.BuildOptions)
	for _, err := range result.Errors {
		logger.Warn(err.Text)
	}
	var outfile string
	for _, file := range result.OutputFiles {
		if filepath.Ext(file.Path) != "map" {
			cwd, _ := os.Getwd()
			outfile = strings.TrimPrefix(file.Path, filepath.Join(cwd, serv.BuildOptions.Outdir))
		}
	}
	return outfile, result
}

func (serv *sorvor) run(entry string) {
	var cmd *exec.Cmd
	var outfile string
	var result api.BuildResult
	wg := new(sync.WaitGroup)
	wg.Add(1)

	var onRebuild = func(result api.BuildResult) {
		if cmd != nil {
			err := cmd.Process.Signal(syscall.SIGINT)
			logger.Fatal(err, "failed to stop ", outfile)
		}
		cmd = exec.Command("node", outfile)
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		err := cmd.Start()
		logger.Fatal(err, "failed to start ", outfile)
	}
	// start esbuild in watch mode
	serv.BuildOptions.Watch = &api.WatchMode{OnRebuild: onRebuild}
	outfile, result = serv.esbuild(entry)
	outfile = filepath.Join(serv.BuildOptions.Outdir, outfile)
	onRebuild(result)
	wg.Wait()
}

func (serv *sorvor) build(pkg *pkgjson.PkgJSON) []string {

	target := filepath.Join(serv.BuildOptions.Outdir, "index.html")

	var entries []string
	if _, err := os.Stat(serv.Entry); err != nil {
		logger.Fatal(err, "Entry file does not exist. ", serv.Entry)
	}

	tmpl, err := template.New("index.html").Funcs(template.FuncMap{
		"livereload": func() template.HTML {
			if serv.Serve == true {
				return template.HTML(livereload.JsSnippeet)
			}
			return ""
		},
		"esbuild": func(entry string) string {
			if serv.Serve == true {
				entry = filepath.Join(filepath.Dir(serv.Entry), entry)
				entries = append(entries, entry)
			} else {
				entry = filepath.Join(filepath.Dir(serv.Entry), entry)
			}
			outfile, _ := serv.esbuild(entry)
			return outfile
		},
	}).ParseFiles(serv.Entry)
	logger.Fatal(err, "Unable to parse index.html")

	file, err := os.Create(target)
	logger.Fatal(err, "Unable to create index.html in outdir")
	defer file.Close()

	err = tmpl.Execute(file, pkg)
	logger.Fatal(err, "Unable to execute index.html")

	return entries
}

func (serv *sorvor) ServeHTTP(res http.ResponseWriter, request *http.Request) {
	res.Header().Set("access-control-allow-origin", "*")
	root := filepath.Join(serv.BuildOptions.Outdir, filepath.Clean(request.URL.Path))

	if stat, err := os.Stat(root); err != nil {
		// Serve a root index when root is not found
		http.ServeFile(res, request, filepath.Join(serv.BuildOptions.Outdir, "index.html"))
		return
	} else if stat.IsDir() {
		// Serve root index when requested root is a directory
		http.ServeFile(res, request, filepath.Join(serv.BuildOptions.Outdir, "index.html"))
		return
	}

	// else just Serve the file normally...
	http.ServeFile(res, request, root)
	return
}

func (serv *sorvor) server(pkg *pkgjson.PkgJSON) {
	liveReload := livereload.New()
	liveReload.Start()
	wg := new(sync.WaitGroup)
	wg.Add(2)

	// start esbuild in watch mode
	go func() {
		serv.BuildOptions.Watch = &api.WatchMode{
			OnRebuild: func(result api.BuildResult) {
				for _, err := range result.Errors {
					logger.Warn(err.Text)
				}
				// send livereload message to connected clients
				liveReload.Reload()
			},
		}
		serv.build(pkg)
	}()

	// start our own server
	go func() {
		http.Handle("/livereload", liveReload)
		http.Handle("/", serv)

		if serv.Secure {
			// generate self signed certs
			if _, err := os.Stat("key.pem"); os.IsNotExist(err) {
				authority.GenerateKeyPair(serv.Host)
			}
			logger.Info(logger.BlueText("sørvør"), "ready on", logger.BlueText("https://", serv.Host, serv.Port))
			err := http.ListenAndServeTLS(serv.Port, "cert.pem", "key.pem", nil)
			logger.Error(err, "Failed to start https server")
		} else {
			logger.Info(logger.BlueText("sørvør"), "ready on", logger.BlueText("http://", serv.Host, serv.Port))
			err := http.ListenAndServe(serv.Port, nil)
			logger.Error(err, "Failed to start http server")
		}
	}()

	wg.Wait()
}

func main() {
	var pkgJSON *pkgjson.PkgJSON
	pkg, err := ioutil.ReadFile("package.json")
	if err == nil {
		pkgJSON, err = pkgjson.Parse(pkg)
	}

	serv := readOptions(pkgJSON)

	err = os.MkdirAll(serv.BuildOptions.Outdir, 0775)
	logger.Fatal(err, "Failed to create output directory")

	if filepath.Ext(serv.Entry) != ".html" {
		if serv.Serve == true {
			serv.run(serv.Entry)
		} else {
			serv.esbuild(serv.Entry)
		}
	} else {
		if serv.Serve == true {
			serv.server(pkgJSON)
		} else {
			serv.build(pkgJSON)
		}
	}
}
