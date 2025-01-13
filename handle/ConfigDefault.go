package handle

import (
	"fmt"
	"os"

	"github.com/infrasonar/infrasonar-cli/conf"
)

func ConfigDefault(set string) {
	config := conf.EnsureConfig(set)
	if set != "" {
		conf.Default(config)
		if err := conf.Write(); err != nil {
			fmt.Fprintln(os.Stderr, "failed to write changes")
			os.Exit(1)
		}
		fmt.Fprintf(os.Stderr, "Default configuration: %s\n", config.Name)
		os.Exit(0)
	}
	fmt.Println(config.Name)
	os.Exit(0)
}
