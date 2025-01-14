package handle

import (
	"fmt"
	"os"

	"github.com/infrasonar/infrasonar-cli/cli"
)

func Version() {
	fmt.Printf("InfraSonar version %s\n", cli.Version)
	os.Exit(0)
}
