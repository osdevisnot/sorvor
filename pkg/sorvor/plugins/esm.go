package plugins

import (
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/evanw/esbuild/pkg/api"
)

var EsmPlugin = api.Plugin{
	Name: "esm",
	Setup: func(build api.PluginBuild) {
		build.OnResolve(api.OnResolveOptions{Filter: `^https?://`},
			func(args api.OnResolveArgs) (api.OnResolveResult, error) {
				return api.OnResolveResult{
					Path:      args.Path,
					Namespace: "http-url",
				}, nil
			})
		build.OnResolve(api.OnResolveOptions{Filter: ".*", Namespace: "http-url"},
			func(args api.OnResolveArgs) (api.OnResolveResult, error) {
				base, err := url.Parse(args.Importer)
				if err != nil {
					return api.OnResolveResult{}, err
				}
				relative, err := url.Parse(args.Path)
				if err != nil {
					return api.OnResolveResult{}, err
				}
				return api.OnResolveResult{
					Path:      base.ResolveReference(relative).String(),
					Namespace: "http-url",
				}, nil
			})

		build.OnLoad(api.OnLoadOptions{Filter: ".*", Namespace: "http-url"},
			func(args api.OnLoadArgs) (api.OnLoadResult, error) {
				res, err := http.Get(args.Path)
				if err != nil {
					return api.OnLoadResult{}, err
				}
				defer res.Body.Close()
				bytes, err := ioutil.ReadAll(res.Body)
				if err != nil {
					return api.OnLoadResult{}, err
				}
				// TODO: write file to filesystem under web-modules/lodash/4.3.0/lodash.js
				contents := string(bytes)
				return api.OnLoadResult{Contents: &contents}, nil
			})
	},
}
