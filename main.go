// an extremely fast, zero config server for modern web applications.
package main

import (
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/evanw/esbuild/pkg/api"
	"github.com/evanw/esbuild/pkg/cli"
	"github.com/osdevisnot/sorvor/pkg/authority"
	"github.com/osdevisnot/sorvor/pkg/livereload"
	"github.com/osdevisnot/sorvor/pkg/logger"
)

type sorvor struct {
	buildOptions api.BuildOptions
	entry        string
	host         string
	port         string
	dev          bool
	secure       bool
}

type npm struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

func readOptions(pkg npm) *sorvor {
	var err error
	var esbuildArgs []string

	osArgs := os.Args[1:]
	serv := &sorvor{}

	for _, arg := range osArgs {
		switch {
		case strings.HasPrefix(arg, "--host"):
			serv.host = arg[len("--host="):]
		case strings.HasPrefix(arg, "--port"):
			port, err := strconv.Atoi(arg[len("--port="):])
			logger.Fatal(err, "Invalid Port Value")
			serv.port = ":" + strconv.Itoa(port)
		case arg == "--dev":
			serv.dev = true
		case arg == "--secure":
			serv.secure = true
		case !strings.HasPrefix(arg, "--"):
			serv.entry = arg
		default:
			esbuildArgs = append(esbuildArgs, arg)
		}
	}

	serv.buildOptions, err = cli.ParseBuildOptions(esbuildArgs)
	logger.Fatal(err, "Invalid option for esbuild")

	serv.buildOptions.Bundle = true
	serv.buildOptions.Write = true

	if serv.dev == true {
		if serv.port == "" {
			serv.port = ":1234"
		}
		serv.buildOptions.Define = map[string]string{"process.env.NODE_ENV": "'development'"}
	} else {
		serv.buildOptions.Define = map[string]string{"process.env.NODE_ENV": "'production'"}
	}
	if serv.buildOptions.Outdir == "" {
		serv.buildOptions.Outdir = "dist"
	}
	if serv.buildOptions.Format == api.FormatDefault {
		serv.buildOptions.Format = api.FormatESModule
	}
	if serv.entry == "" {
		serv.entry = "public/index.html"
	}
	if serv.host == "" {
		serv.host = "localhost"
	}
	return serv
}

func readPkg() npm {
	pkg := npm{}

	// file, err := ioutil.ReadFile("package.json")
	// handleError("Unable to read package.json", err, false)

	// err = json.Unmarshal(file, &pkg)
	// handleError("Unable to parse package.json", err, false)

	return pkg
}

func (serv *sorvor) esbuild(entry string) string {
	serv.buildOptions.EntryPoints = []string{entry}
	result := api.Build(serv.buildOptions)
	for _, err := range result.Errors {
		logger.Warn(err.Text)
	}
	var outfile string
	for _, file := range result.OutputFiles {
		if filepath.Ext(file.Path) != "map" {
			cwd, _ := os.Getwd()
			outfile = strings.TrimPrefix(file.Path, filepath.Join(cwd, serv.buildOptions.Outdir))
		}
	}
	return outfile
}

func (serv *sorvor) build(pkg npm) []string {

	target := filepath.Join(serv.buildOptions.Outdir, "index.html")

	var entries []string
	if _, err := os.Stat(serv.entry); err != nil {
		logger.Fatal(err, "Entry file does not exist. ", serv.entry)
	}

	tmpl, err := template.New("index.html").Funcs(template.FuncMap{
		"livereload": func() template.HTML {
			if serv.dev == true {
				return template.HTML(livereload.Snippet)
			}
			return ""
		},
		"esbuild": func(entry string) string {
			if serv.dev == true {
				entry = filepath.Join(filepath.Dir(serv.entry), entry)
				entries = append(entries, entry)
			} else {
				entry = filepath.Join(filepath.Dir(serv.entry), entry)
			}
			return serv.esbuild(entry)
		},
	}).ParseFiles(serv.entry)
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
	root := filepath.Join(serv.buildOptions.Outdir, filepath.Clean(request.URL.Path))

	if stat, err := os.Stat(root); err != nil {
		// serve a root index when root is not found
		http.ServeFile(res, request, filepath.Join(serv.buildOptions.Outdir, "index.html"))
		return
	} else if stat.IsDir() {
		// serve root index when requested root is a directory
		http.ServeFile(res, request, filepath.Join(serv.buildOptions.Outdir, "index.html"))
		return
	}

	// else just serve the file normally...
	http.ServeFile(res, request, root)
	return
}

func (serv *sorvor) serve(pkg npm) {
	liveReload := livereload.New()
	liveReload.Start()
	wg := new(sync.WaitGroup)
	wg.Add(2)

	// start esbuild in watch mode
	go func() {
		serv.buildOptions.Watch = &api.WatchMode{
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

		if serv.secure {
			// generate self signed certs
			if _, err := os.Stat("key.pem"); os.IsNotExist(err) {
				authority.GenerateKeyPair(serv.host)
			}
			logger.Info(logger.BlueText("sørvør"), "ready on", logger.BlueText("https://", serv.host, serv.port))
			err := http.ListenAndServeTLS(serv.port, "cert.pem", "key.pem", nil)
			logger.Error(err, "Failed to start https server")
		} else {
			logger.Info(logger.BlueText("sørvør"), "ready on", logger.BlueText("http://", serv.host, serv.port))
			err := http.ListenAndServe(serv.port, nil)
			logger.Error(err, "Failed to start http server")
		}
	}()

	wg.Wait()
}

func main() {
	pkg := readPkg()
	serv := readOptions(pkg)

	err := os.MkdirAll(serv.buildOptions.Outdir, 0775)
	logger.Fatal(err, "Failed to create output directory")

	if serv.dev == true {
		serv.serve(pkg)
	} else if filepath.Ext(serv.entry) != ".html" {
		serv.esbuild(serv.entry)
	} else {
		serv.build(pkg)
	}
}
