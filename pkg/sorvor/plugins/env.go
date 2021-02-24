package plugins

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/evanw/esbuild/pkg/api"
)

// EnvPlugin reads environment variable
var EnvPlugin = api.Plugin{
	Name: "env",
	Setup: func(build api.PluginBuild) {
		build.OnResolve(api.OnResolveOptions{Filter: `^env$`},
			func(args api.OnResolveArgs) (api.OnResolveResult, error) {
				return api.OnResolveResult{
					Path:      args.Path,
					Namespace: "env-ns",
				}, nil
			})
		build.OnLoad(api.OnLoadOptions{Filter: `.*`, Namespace: "env-ns"},
			func(args api.OnLoadArgs) (api.OnLoadResult, error) {
				mappings := make(map[string]string)
				for _, item := range os.Environ() {
					if equals := strings.IndexByte(item, '='); equals != -1 {
						mappings[item[:equals]] = item[equals+1:]
					}
				}
				bytes, err := json.Marshal(mappings)
				if err != nil {
					return api.OnLoadResult{}, err
				}
				contents := string(bytes)
				return api.OnLoadResult{
					Contents: &contents,
					Loader:   api.LoaderJSON,
				}, nil
			})
	},
}
