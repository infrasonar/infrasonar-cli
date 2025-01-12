package handle

import (
	"fmt"
	"os"
)

// Update README.md when upgrading to a new release
const version = "1.0.0"

func Version() {
	fmt.Printf("InfraSonar version %s\n", version)
	os.Exit(0)
}
