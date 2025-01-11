package main

import (
	"errors"
	"fmt"
	"os"
	"path"

	"gopkg.in/yaml.v3"
)

type Project struct {
	name   string `yaml:"name"`
	token  string `yaml:"token"`
	api    string `yaml:"api"`
	output string `yaml:"output"`
}

type Projects struct {
	projects []Project `yaml:"projects"`
}

var projectList Projects = Projects{}

func readProjects() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to read home path: %s", err)
		return
	}
	fn := path.Join(homeDir, ".infrasonar_cli_projects.yaml")
	if _, err := os.Stat(fn); errors.Is(err, os.ErrNotExist) {
		return
	}

	content, err := os.ReadFile(fn)
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to read '%s': %s", fn, err)
		return
	}
	err = yaml.Unmarshal(content, &projectList)
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to unpack '%s': %s", fn, err)
	}
}
