package handle

import (
	"fmt"
	"os"

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
		fmt.Fprintln(os.Stderr, "failed to write changes")
		os.Exit(1)
	}
	fmt.Fprintln(os.Stderr, "Configuration removed")
	os.Exit(0)
}
