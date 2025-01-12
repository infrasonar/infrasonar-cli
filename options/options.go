package options

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/akamensky/argparse"
	"github.com/infrasonar/infrasonar-cli/re"
)

var DefaultOutput = &argparse.Options{
	Required: false,
	Validate: func(args []string) error {
		switch args[0] {
		case "json", "yaml", "simple":
			return nil
		}

		return fmt.Errorf("unknown '%s' {yaml,json,simple}", args[0])
	},
	Default: "yaml",
	Help:    "Default output format. {yaml,json,simple}",
}

var Output = &argparse.Options{
	Required: false,
	Validate: func(args []string) error {
		switch args[0] {
		case "", "json", "yaml", "simple":
			return nil
		}

		return fmt.Errorf("unknown '%s' {yaml,json,simple}", args[0])
	},
	Help: "Output format. {yaml,json,simple}",
}

var ConfigName = &argparse.Options{
	Required: false,
	Validate: func(args []string) error {
		if !re.ConfigName.MatchString(args[0]) {
			return errors.New("invalid configuration name")
		}
		return nil
	},
	Help: "Configuration name",
}

var ConfigNewApi = &argparse.Options{
	Required: false,
	Validate: func(args []string) error {
		if !re.IsUrl.MatchString(args[0]) {
			return errors.New("invalid API URL")
		}
		return nil
	},
	Help:    "InfraSonar API url for the project",
	Default: "https://api.infrasonar.com",
}

var ConfigUpdateApi = &argparse.Options{
	Required: false,
	Validate: func(args []string) error {
		if !re.IsUrl.MatchString(args[0]) {
			return errors.New("invalid API URL")
		}
		return nil
	},
	Help: "InfraSonar API url for the project",
}

var Token = &argparse.Options{
	Required: false,
	Validate: func(args []string) error {
		if !re.Token.MatchString(args[0]) {
			return errors.New("invalid token")
		}
		return nil
	},
	Help: "Token for authentication with the InfraSonar API",
}

var Container = &argparse.Options{
	Required: false,
	Validate: func(args []string) error {
		if containerId, err := strconv.Atoi(args[0]); err == nil {
			if containerId <= 0 {
				return errors.New("expecting a value greater than 0")
			}
		}
		return nil
	},
	Help: "Container ID",
}

var AssetFilter = &argparse.Options{
	Required: false,
	Validate: func(args []string) error {
		for _, arg := range args {
			m := re.AssetFilter.FindStringSubmatch(arg)
			if m == nil {
				return fmt.Errorf("invalid '%s'. valid example: -f kind==linux -f collector==snmp -f label!=123 -f zone=0", arg)
			}
			switch m[1] {
			case "collector", "kind":
				continue
			case "label", "zone":
				if !re.Number.MatchString(m[3]) {
					return fmt.Errorf("%s must be compared with a %s ID, for example: %s%s123", m[1], m[1], m[1], m[2])
				}
				continue
			}
			return fmt.Errorf("unknown '%s'. {kind,collector,label}", m[1])
		}
		return nil
	},
	Help: "Filter assets. Multiple filters are allowed, for example: -f kind==linux -f collector==snmp -f label!=123 -f zone=0",
}

var AssetProperties = selectorList(
	false,
	[]string{"id", "name", "kind", "description", "mode", "labels", "collectors", "properties"},
	"Asset properties to return. If not specified all properties will be returned",
)
