package conf

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/infrasonar/infrasonar-cli/cli"
	"gopkg.in/yaml.v3"
)

var configurationsFn string
var conf = Configurations{}

func Def() *Config {
	if len(conf.Configs) == 0 {
		return nil
	}
	return conf.Configs[0]
}

func EnsureConfig(name string) *Config {
	if name == "" {
		config := Def()
		if config == nil {
			fmt.Fprintf(os.Stderr, "It appears no configuration has been set up.\nYou can create a new configuration using this command:\n\n  %s config new\n\n", filepath.Base(os.Args[0]))
			os.Exit(1)
		}
		return config
	}
	config := conf.get(name)
	if config == nil {
		fmt.Fprintf(os.Stderr, "Configuration '%s' not found. Use this command to list the configurations:\n\n  %s config list\n\n", name, filepath.Base(os.Args[0]))
		os.Exit(1)
	}
	return config
}

func Delete(c *Config) {
	n := []*Config{}
	for _, other := range conf.Configs {
		if other != c {
			n = append(n, other)
		}
	}
	conf.Configs = n
}

func Default(c *Config) {
	n := []*Config{}
	for _, other := range conf.Configs {
		if other == c {
			n = append(n, other)
		}
	}
	for _, other := range conf.Configs {
		if other != c {
			n = append(n, other)
		}
	}
	conf.Configs = n
}

func New(name, token, api, output string) (*Config, error) {
	if conf.get(name) != nil {
		return nil, fmt.Errorf("configuration '%s' already exists", name)
	}

	config := Config{
		Name:   name,
		Api:    api,
		Output: output,
	}

	// Set token
	if err := config.SetToken(token); err != nil {
		return nil, err
	}

	// Append the config
	conf.Configs = append(conf.Configs, &config)

	// Write configurations to disk
	err := Write()
	return &config, err
}

func Write() error {
	data, err := yaml.Marshal(conf)
	if err != nil {
		return err
	}
	return os.WriteFile(configurationsFn, data, 0644)
}

func GetConfigs() []*Config {
	return conf.Configs
}

func Initialize() error {
	cliPath, err := cli.CliPath()
	if err != nil {
		return err
	}

	configurationsFn = path.Join(cliPath, "configs.yaml")
	if _, err := os.Stat(configurationsFn); errors.Is(err, os.ErrNotExist) {
		return nil // success, just no configuration yet
	}

	content, err := os.ReadFile(configurationsFn)
	if err != nil {
		return fmt.Errorf("failed to read '%s': %s", configurationsFn, err)
	}
	err = yaml.Unmarshal(content, &conf)
	if err != nil {
		return fmt.Errorf("failed to unpack '%s': %s", configurationsFn, err)
	}
	return nil
}
