package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/evanw/esbuild/pkg/api"
	"github.com/evanw/esbuild/pkg/cli"
	"github.com/osdevisnot/sorvor/pkg/logger"
	"github.com/osdevisnot/sorvor/pkg/pkgjson"
	"github.com/osdevisnot/sorvor/pkg/sorvor"
)

var version = "development"

func readOptions(pkgJSON *pkgjson.PkgJSON) *sorvor.Sorvor {
	var err error
	var esbuildArgs []string

	osArgs := os.Args[1:]
	serv := &sorvor.Sorvor{}

	for _, arg := range osArgs {
		switch {
		case strings.HasPrefix(arg, "--version"):
			logger.Info("sørvør version", version)
			os.Exit(0)
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
	if serv.BuildOptions.Platform == api.PlatformNode {
		for key, _ := range pkgJSON.Dependencies {
			serv.BuildOptions.External = append(serv.BuildOptions.External, key)
		}
	}
	return serv
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
			serv.RunEntry(serv.Entry)
		} else {
			serv.BuildEntry(serv.Entry)
		}
	} else {
		if serv.Serve == true {
			serv.ServeIndex(pkgJSON)
		} else {
			serv.BuildIndex(pkgJSON)
		}
	}
}
