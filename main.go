// an extremely fast, zero config server for modern web applications.
package main

import (
	"html/template"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/evanw/esbuild/pkg/api"
	"github.com/evanw/esbuild/pkg/cli"
	"github.com/osdevisnot/sorvor/pkg/livereload"
	"github.com/osdevisnot/sorvor/pkg/logger"
)

type sorvor struct {
	buildOptions api.BuildOptions
	serveOptions api.ServeOptions
	src          string
	port         string
	dev          bool
}

type npm struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

var extensions = map[string]string{
	".js":  ".js",
	".ts":  ".js",
	".jsx": ".js",
	".tsx": ".js",
	".css": ".css",
}

func readOptions(pkg npm) *sorvor {
	var err error
	var esbuildArgs []string

	osArgs := os.Args[1:]
	serv := &sorvor{}

	for _, arg := range osArgs {
		switch {
		case strings.HasPrefix(arg, "--src"):
			serv.src = arg[len("--src="):]
		case strings.HasPrefix(arg, "--port"):
			port, err := strconv.Atoi(arg[len("--port="):])
			logger.Fatal(err, "Invalid Port Value")
			serv.port = ":" + strconv.Itoa(port)
			serv.serveOptions.Port = uint16(port + 1)
		case arg == "--dev":
			serv.dev = true
		default:
			esbuildArgs = append(esbuildArgs, arg)
		}
	}

	serv.buildOptions, err = cli.ParseBuildOptions(esbuildArgs)
	logger.Fatal(err, "Invalid option for esbuild")

	serv.buildOptions.Bundle = true
	serv.buildOptions.Write = true

	if serv.dev == true && serv.port == "" {
		serv.port = ":1234"
		serv.serveOptions.Port = uint16(1235)
	}
	if serv.buildOptions.Outdir == "" {
		serv.buildOptions.Outdir = "dist"
	}
	if serv.src == "" {
		serv.src = "src"
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
	serv.buildOptions.EntryPoints = []string{filepath.Join(serv.src, entry)}
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

	srcIndex := filepath.Join(serv.src, "index.html")
	targetIndex := filepath.Join(serv.buildOptions.Outdir, "index.html")

	var entries []string

	tmpl, err := template.New("index.html").Funcs(template.FuncMap{
		"livereload": func() template.JS {
			if serv.dev == true {
				return template.JS(livereload.Snippet)
			}
			return ""
		},
		"esbuild": func(entry string) string {
			if serv.dev == true {
				entries = append(entries, filepath.Join(serv.src, entry))
				ext := path.Ext(entry)
				outfile := entry[0:len(entry)-len(ext)] + extensions[ext]
				return "http://localhost:" + strconv.Itoa(int(serv.serveOptions.Port)) + "/" + outfile
			}
			return serv.esbuild(entry)
		},
	}).ParseFiles(srcIndex)
	logger.Fatal(err, "Unable to parse index.html")

	file, err := os.Create(targetIndex)
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
	serv.buildOptions.EntryPoints = serv.build(pkg)

	wg := new(sync.WaitGroup)
	wg.Add(2)

	// start esbuild server
	go func() {
		_, err := api.Serve(serv.serveOptions, serv.buildOptions)
		if err != nil {
			logger.Error(err, "Failed to start esbuild server")
			wg.Done()
		}
	}()

	// start our own server
	go func() {
		logger.Info(logger.BlueText("sørvør"), "ready on", logger.BlueText("http://localhost", serv.port))

		liveReload := livereload.New(serv.src)
		liveReload.Start()
		http.Handle("/livereload", liveReload)

		http.Handle("/", serv)

		err := http.ListenAndServe(serv.port, nil)
		logger.Error(err, "Failed to start http server")
		wg.Done()
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
	} else {
		serv.build(pkg)
	}
}
