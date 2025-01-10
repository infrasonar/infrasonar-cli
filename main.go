package main

import (
	"fmt"
	"os"
)

func main() {
	// Read the arguments
	if err := parseArgs(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
