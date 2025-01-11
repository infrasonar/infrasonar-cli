package main

import (
	"fmt"
	"os"
)

type TGetAssets struct {
	config     *Config
	output     string
	container  int
	properties []string
	filters    []string
}

func handleGetAssets(cmd *TGetAssets) {
	fmt.Println(cmd)
	os.Exit(0)
}
