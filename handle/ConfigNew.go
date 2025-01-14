package handle

import (
	"github.com/infrasonar/infrasonar-cli/conf"
	"github.com/infrasonar/infrasonar-cli/handle/util"
)

type TConfigNew struct {
	Name   string
	Token  string
	Api    string
	Output string
}

func ConfigNew(cmd *TConfigNew) {
	if cmd.Name == "" {
		cmd.Name = util.AskConfigName()
	}
	if cmd.Token == "" {
		cmd.Token = util.AskToken()
	}
	config, err := conf.New(cmd.Name, cmd.Token, cmd.Api, cmd.Output)
	util.ExitOnErr(err)
	util.ExitOk("Configuration '%s' created\n", config.Name)
}
