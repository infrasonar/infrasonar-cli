package cli

import (
	"errors"
	"path/filepath"
	"strings"
)

var AssetProperties = []string{"id", "name", "kind", "zone", "description", "mode", "labels", "collectors", "properties"}
var CollectorProperties = []string{"key", "name", "kind", "info", "minVersion", "checks"}

type IntSet map[int]struct{}

func (s IntSet) Set(k int) {
	s[k] = struct{}{}
}

func GetJsonOrYaml(fn string) (string, error) {
	extension := strings.ToLower(filepath.Ext(fn))
	switch extension {
	case ".yaml", ".yml":
		return "yaml", nil
	case ".json":
		return "json", nil
	}
	return "", errors.New("expecting a .json or .yml/.yaml file")
}
