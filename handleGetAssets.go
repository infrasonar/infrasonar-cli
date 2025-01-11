package main

import (
	"fmt"
	"os"
)

type TGetAssets struct {
	project     *Project
	output      string
	containerId int
	properties  []string
	filters     []string
}

func handleGetAssets(cmd *TGetAssets) {
	fmt.Println(cmd)
	os.Exit(0)
}
