package cli

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
)

var AssetProperties = []string{"id", "name", "kind", "zone", "description", "mode", "labels", "collectors", "properties"}
var CollectorProperties = []string{"key", "name", "kind", "info", "minVersion", "checks"}
var cliPath string

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

func CliPath() (string, error) {
	if cliPath == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to read home path: %s", err)
		}

		cliPath = path.Join(homeDir, ".infrasonar_cli")
		err = os.MkdirAll(cliPath, os.ModePerm)
		if err != nil {
			return "", fmt.Errorf("failed to make directory '%s': %s", cliPath, err)
		}
	}
	return cliPath, nil
}
