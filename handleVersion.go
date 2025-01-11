package main

import (
	"fmt"
	"os"
)

func handleVersion () {
	fmt.Printf("InfraSonar version %s\n", Version)
	os.Exit(0)
}