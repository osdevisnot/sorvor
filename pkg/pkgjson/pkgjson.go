package pkgjson

import (
	"encoding/json"
)

// PkgJSON represents a node.js `package.json`
type PkgJSON struct {
	Name             string            `json:"name"`
	Version          string            `json:"version"`
	Author           string            `json:"author"`
	Description      string            `json:"description"`
	Dependencies     map[string]string `json:"dependencies"`
	PeerDependencies map[string]string `json:"peerDependencies"`
	DevDependencies  map[string]string `json:"devDependencies"`
}

// Parse parses `package.json` and returns the structure.
func Parse(payload []byte) (*PkgJSON, error) {
	var packageJSON *PkgJSON
	err := json.Unmarshal(payload, &packageJSON)
	return packageJSON, err
}
