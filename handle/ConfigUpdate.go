package handle

import (
	"fmt"
	"os"

	"github.com/infrasonar/infrasonar-cli/conf"
	"github.com/infrasonar/infrasonar-cli/handle/util"
)

type TConfigUpdate struct {
	Name   string
	Token  string
	Api    string
	Output string
}

func ConfigUpdate(cmd *TConfigUpdate) {
	isChanged := false
	if cmd.Name == "" {
		cmd.Name = util.AskConfigName()
	}

	config := conf.EnsureConfig(cmd.Name)
	if cmd.Api != "" && config.Api != cmd.Api {
		config.Api = cmd.Api
		fmt.Println("API updated")
		isChanged = true
	}
	if cmd.Token != "" {
		config.EncToken = cmd.Output
		fmt.Println("Token updated")
		isChanged = true
	}
	if cmd.Output != "" && config.Output != cmd.Output {
		config.Output = cmd.Output
		fmt.Println("Default output updated")
		isChanged = true
	}

	if isChanged {
		if err := conf.Write(); err != nil {
			fmt.Fprintln(os.Stderr, "failed to write changes")
			os.Exit(1)
		}
		fmt.Println("Changes written")
	} else {
		fmt.Println("No changes")
	}
	os.Exit(0)
}
