package handle

import (
	"github.com/infrasonar/infrasonar-cli/conf"
	"github.com/infrasonar/infrasonar-cli/handle/util"
)

func ConfigDelete(name string) {
	if name == "" {
		name = util.AskConfigName()
	}
	config := conf.EnsureConfig(name)
	conf.Delete(config)
	if err := conf.Write(); err != nil {
		util.ExitErr("failed to write changes")
	}
	util.ExitOk("Configuration removed")
}
