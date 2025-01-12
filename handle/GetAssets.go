package handle

import (
	"fmt"
	"os"
)

type TGetAssets struct {
	Api        string
	Token      string
	Output     string
	Container  int
	Properties []string
	Filters    []string
}

func GetAssets(cmd *TGetAssets) {
	
	fmt.Println(cmd)
	os.Exit(0)
}
