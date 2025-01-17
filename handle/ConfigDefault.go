package handle

import (
	"github.com/infrasonar/infrasonar-cli/conf"
	"github.com/infrasonar/infrasonar-cli/handle/util"
)

func ConfigDefault() {
	config := conf.EnsureConfig("")
	util.ExitOk(config.Name)
}
