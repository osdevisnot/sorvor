// Package sorvor is an extremely fast, zero config ServeIndex for modern web applications.
package sorvor

import (
	"html/template"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"syscall"

	"github.com/evanw/esbuild/pkg/api"
	"github.com/osdevisnot/sorvor/pkg/cert"
	"github.com/osdevisnot/sorvor/pkg/livereload"
	"github.com/osdevisnot/sorvor/pkg/logger"
	"github.com/osdevisnot/sorvor/pkg/pkgjson"
	"github.com/osdevisnot/sorvor/pkg/sorvor/plugins"
)

// Sorvor struct
type Sorvor struct {
	BuildOptions api.BuildOptions
	Entry        string
	Host         string
	Port         string
	Serve        bool
	Secure       bool
}

// BuildEntry builds a given entrypoint using esbuild
func (serv *Sorvor) BuildEntry(entry string) (string, api.BuildResult) {
	serv.BuildOptions.EntryPoints = []string{entry}
	serv.BuildOptions.Plugins = []api.Plugin{plugins.EnvPlugin, plugins.EsmPlugin}
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

// RunEntry builds an entrypoint and launches the resulting built file using node.js
func (serv *Sorvor) RunEntry(entry string) {
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
	outfile, result = serv.BuildEntry(entry)
	outfile = filepath.Join(serv.BuildOptions.Outdir, outfile)
	onRebuild(result)
	wg.Wait()
}

// BuildIndex walks the index.html, collect all the entries from <script...></script> and <link .../> tags
// it then runs it through esbuild and replaces the references in index.html with new paths
func (serv *Sorvor) BuildIndex(pkg *pkgjson.PkgJSON) []string {

	target := filepath.Join(serv.BuildOptions.Outdir, "index.html")

	var entries []string
	if _, err := os.Stat(serv.Entry); err != nil {
		logger.Fatal(err, "Entry file does not exist. ", serv.Entry)
	}

	tmpl, err := template.New("index.html").Funcs(template.FuncMap{
		"livereload": func() template.HTML {
			if serv.Serve == true {
				return template.HTML(livereload.JsSnippet)
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
			outfile, _ := serv.BuildEntry(entry)
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

// ServeHTTP is an http server handler for sorvor
func (serv *Sorvor) ServeHTTP(res http.ResponseWriter, request *http.Request) {
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

// ServeIndex launches esbuild in watch mode and live reloads all connected browsers
func (serv *Sorvor) ServeIndex(pkg *pkgjson.PkgJSON) {
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
		serv.BuildIndex(pkg)
	}()

	// start our own ServeIndex
	go func() {
		http.Handle("/livereload", liveReload)
		http.Handle("/", serv)

		if serv.Secure {
			// generate self signed certs
			if _, err := os.Stat("key.pem"); os.IsNotExist(err) {
				cert.GenerateKeyPair(serv.Host)
			}
			logger.Info(logger.BlueText("sørvør"), "ready on", logger.BlueText("https://", serv.Host, serv.Port))
			err := http.ListenAndServeTLS(serv.Port, "cert.pem", "key.pem", nil)
			logger.Error(err, "Failed to start https ServeIndex")
		} else {
			logger.Info(logger.BlueText("sørvør"), "ready on", logger.BlueText("http://", serv.Host, serv.Port))
			err := http.ListenAndServe(serv.Port, nil)
			logger.Error(err, "Failed to start http ServeIndex")
		}
	}()

	wg.Wait()
}
