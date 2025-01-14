package handle

import (
	"github.com/infrasonar/infrasonar-cli/conf"
	"github.com/infrasonar/infrasonar-cli/handle/util"
)

func ConfigDefault(set string) {
	config := conf.EnsureConfig(set)
	if set != "" {
		conf.Default(config)
		if err := conf.Write(); err != nil {
			util.ExitErr("failed to write changes")
		}
		util.ExitOk("Default configuration: %s\n", config.Name)
	}
	util.ExitOk(config.Name)
}
