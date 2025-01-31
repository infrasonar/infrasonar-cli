package options

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/akamensky/argparse"
	"github.com/infrasonar/infrasonar-cli/cli"
	"github.com/infrasonar/infrasonar-cli/handle/util"
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

var IncludeDefaults = &argparse.Options{
	Required: false,
	Help:     "Include default collector configuration values. By default, values equal to the default are not included",
}

var ConfigName = &argparse.Options{
	Required: false,
	Validate: func(args []string) error {
		if !re.MetaKey.MatchString(args[0]) {
			fmt.Println(args[0])
			return errors.New("invalid configuration name")
		}
		return nil
	},
	Help: "Configuration name",
}

var UseConfig = &argparse.Options{
	Required: false,
	Validate: func(args []string) error {
		if !re.MetaKey.MatchString(args[0]) {
			fmt.Println(args[0])
			return errors.New("invalid configuration name")
		}
		return nil
	},
	Help: fmt.Sprintf("Use an alternative configuration. View the default config with: '%s config default'", filepath.Base(os.Args[0])),
}

var ApplyFileName = &argparse.Options{
	Required: true,
	Validate: func(args []string) error {
		fn := args[0]
		if _, err := os.Stat(fn); errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("file does not exist: %s", fn)
		}
		_, err := cli.GetJsonOrYaml(fn)
		return err
	},
	Help: "YAML of JSON input filename to apply",
}

var OutFileName = &argparse.Options{
	Required: false,
	Help:     "Write output to this filename",
}

var Collector = &argparse.Options{
	Required: false,
	Validate: func(args []string) error {
		if !re.MetaKey.MatchString(args[0]) {
			fmt.Println(args[0])
			return errors.New("invalid collector key")
		}
		return nil
	},
	Help: "Collector key",
}

var ConfigSetDefault = &argparse.Options{
	Required: false,
	Help:     "Set as the default configuration",
}

var ConfigListMore = &argparse.Options{
	Required: false,
	Help:     "List with more detailed configuration information",
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

var Asset = &argparse.Options{
	Required: false,
	Validate: func(args []string) error {
		if containerId, err := strconv.Atoi(args[0]); err == nil {
			if containerId <= 0 {
				return errors.New("expecting a value greater than 0")
			}
		}
		return nil
	},
	Help: "Asset ID",
}

var Purge = &argparse.Options{
	Required: false,
	Help:     "Deletes existing labels and collectors if not specified. Without the 'purge' flag, only new labels, collectors, and configuration changes are applied",
}

var DryRun = &argparse.Options{
	Required: false,
	Help:     "Dry run mode. Simulate the changes that would be made without actually applying them. Displays a list of proposed changes",
}

var AssetFilter = &argparse.Options{
	Required: false,
	Validate: func(args []string) error {
		seen := map[string]bool{}
		for _, arg := range args {
			m := re.AssetFilter.FindStringSubmatch(arg)
			// Check for filter syntax
			if m == nil {
				return fmt.Errorf("invalid '%s'. valid example: -f kind==Linux -f collector==snmp -f label!=123 -f zone=0", arg)
			}

			// Check for double filters
			key := fmt.Sprintf("%s%s", m[1], m[2])
			if _, exists := seen[key]; exists {
				return fmt.Errorf("double '%s', each filter may only be applied once", key)
			}
			seen[key] = true

			// Check for valid filters
			switch m[1] {
			case "mode":
				if util.InSlice([]string{"normal", "maintenance", "disabled"}, m[3]) == nil {
					return fmt.Errorf("unknown mode '%s'. {normal,maintenance,disabled}", m[3])
				}
				continue
			case "collector", "kind":
				continue
			case "label", "zone":
				if !re.Number.MatchString(m[3]) {
					return fmt.Errorf("%s must be compared with a %s ID, for example: %s%s123", m[1], m[1], m[1], m[2])
				}
				continue
			}
			return fmt.Errorf("unknown filter '%s'. {collector,kind,label,mode,zone}", m[1])
		}
		return nil
	},
	Help: "Filter assets. Multiple filters are allowed, for example: -f kind==linux -f collector==snmp -f label!=123 -f zone=0",
}

var AssetProperties = selectorList(
	false,
	cli.AssetProperties,
	"Asset properties to return (comma-separated). If omitted, all properties will be returned",
)

var CollectorProperties = selectorList(
	false,
	cli.CollectorProperties,
	"Collector properties to return (comma-separated). If omitted, all properties will be returned",
)

var MeProperties = selectorList(
	false,
	cli.MeProperties,
	"Info properties to return (comma-separated). If omitted, all properties will be returned",
)
