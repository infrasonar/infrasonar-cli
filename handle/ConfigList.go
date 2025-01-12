package handle

import (
	"fmt"
	"os"

	"github.com/infrasonar/infrasonar-cli/conf"
)

func ConfigList(more bool) {
	if more {
		mName := 4
		for _, config := range conf.GetConfigs() {
			if len(config.Name) > mName {
				mName = len(config.Name)
			}
		}
		fmt.Printf("%-*s    %-6s    %s\n", mName, "NAME", "OUTPUT", "API")
		for _, config := range conf.GetConfigs() {
			fmt.Printf("%-*s    %-6s    %s\n", mName, config.Name, config.Output, config.Api)
		}

	} else {
		for _, config := range conf.GetConfigs() {
			fmt.Println(config.Name)
		}
	}
	os.Exit(0)
}
