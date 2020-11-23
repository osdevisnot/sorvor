package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/evanw/esbuild/pkg/api"
	"github.com/evanw/esbuild/pkg/cli"
)

func handleError(message string, err error, shouldExit bool) {
	if err != nil {
		if shouldExit {
			log.Fatalf("%s : %v\n", message, err)
		} else {
			log.Printf("%s : %v\n", message, err)
		}
	}
}

type sorvor struct {
	buildOptions api.BuildOptions
	serveOptions api.ServeOptions
	src          string
	port         string
	dev          bool
}

func readOptions(pkg npmPackage) sorvor {
	var err error
	var esbuildArgs []string

	osArgs := os.Args[1:]
	serv := sorvor{}

	for _, arg := range osArgs {
		switch {
		case strings.HasPrefix(arg, "--src"):
			serv.src = arg[len("--src="):]
		case strings.HasPrefix(arg, "--port"):
			port, err := strconv.Atoi(arg[len("--port="):])
			handleError("Invalid Port Value", err, true)
			serv.port = ":" + strconv.Itoa(port)
			serv.serveOptions.Port = uint16(port + 1)
		case arg == "--dev":
			serv.dev = true
		default:
			esbuildArgs = append(esbuildArgs, arg)
		}
	}

	serv.buildOptions, err = cli.ParseBuildOptions(esbuildArgs)
	handleError("Invalid option for esbuild", err, true)

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

type npmPackage struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

func readNpmPackage() npmPackage {
	pkg := npmPackage{}

	file, err := ioutil.ReadFile("package.json")
	handleError("Unable to read package.json", err, false)

	err = json.Unmarshal(file, &pkg)
	handleError("Unable to parse package.json", err, false)

	return pkg
}

func (serv sorvor) esbuild(entry string) string {
	serv.buildOptions.EntryPoints = []string{filepath.Join(serv.src, entry)}
	result := api.Build(serv.buildOptions)
	var outfile string
	for _, file := range result.OutputFiles {
		if filepath.Ext(file.Path) != "map" {
			cwd, _ := os.Getwd()
			outfile = strings.TrimPrefix(file.Path, filepath.Join(cwd, serv.buildOptions.Outdir))
			fmt.Printf("created file %s from %s\n", outfile, entry)
		}
	}
	return outfile
}

func (serv sorvor) build(pkg npmPackage) []string {

	srcIndex := filepath.Join(serv.src, "index.html")
	targetIndex := filepath.Join(serv.buildOptions.Outdir, "index.html")

	var entries []string

	tmpl, err := template.New("index.html").Funcs(template.FuncMap{
		"esbuild": func(entry string) string {
			if serv.dev == true {
				entries = append(entries, filepath.Join(serv.buildOptions.Outdir, entry))
				return "http://localhost:8000/" + entry
			} else {
				return serv.esbuild(entry)
			}
		},
	}).ParseFiles(srcIndex)
	handleError("Unable to parse index.html", err, true)

	file, err := os.Create(targetIndex)
	handleError("Unable to create index.html in out dir", err, true)
	defer file.Close()

	err = tmpl.Execute(file, pkg)
	handleError("Unable to execute index.html", err, true)

	return entries
}

func (serv sorvor) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	path := filepath.Join(serv.buildOptions.Outdir, filepath.Clean(request.URL.Path))

	if stat, err := os.Stat(path); err != nil {
		// serve a root index when path is not found
		http.ServeFile(writer, request, filepath.Join(serv.buildOptions.Outdir, "index.html"))
		return
	} else if stat.IsDir() {
		// serve root index when requested path is a directory
		http.ServeFile(writer, request, filepath.Join(serv.buildOptions.Outdir, "index.html"))
		return
	}

	// else just serve the file normally...
	http.ServeFile(writer, request, path)
	return
}

func (serv sorvor) serve(pkg npmPackage) {
	serv.buildOptions.EntryPoints = serv.build(pkg)

	wg := new(sync.WaitGroup)
	wg.Add(2)

	// start esbuild server
	go func() {
		_, err := api.Serve(api.ServeOptions{}, serv.buildOptions)
		if err != nil {
			handleError("Failed to start esbuild server", err, false)
			wg.Done()
		}
	}()

	// start our own server
	go func() {
		log.Printf("Sorvor Ready on http://localhost%s\n", serv.port)
		err := http.ListenAndServe(serv.port, &serv)
		handleError("Unable to start http server", err, false)
		wg.Done()
	}()

	wg.Wait()
}

func main() {
	pkg := readNpmPackage()
	serv := readOptions(pkg)

	err := os.MkdirAll(serv.buildOptions.Outdir, 0775)
	handleError("Unable to create output directory", err, true)

	if serv.dev == true {
		serv.serve(pkg)
	} else {
		serv.build(pkg)
	}
}
