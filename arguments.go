package main

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/akamensky/argparse"
)

var reAssetFilter = regexp.MustCompile(`^(\w+)(\=\=|\!\=)(\w+)$`)
var reNumber = regexp.MustCompile(`^[0-9]+$`)
var reIsUrl = regexp.MustCompile(`^https?://\S+$`)
var reToken = regexp.MustCompile(`^[0-9a-f]{32}$`)
var reConfigName = regexp.MustCompile(`^[a-zA-Z_]\w*$`)
var tokenValidation = func(args []string) error {
	if !reToken.MatchString(args[0]) {
		return errors.New("invalid token")
	}
	return nil
}

func getOutput(outputArg string, config *Config) string {
	if outputArg == "" {
		return config.Output
	}
	return outputArg
}

func selectorList(required bool, allowed []string, help string) *argparse.Options {
	return &argparse.Options{
		Required: required,
		Validate: func(args []string) error {
			seen := map[string]*struct{}{}

			for _, choice := range strings.Split(args[0], ",") {
				if !slices.Contains(allowed, choice) {
					return fmt.Errorf("invalid '%s'", choice)
				}
				if _, ok := seen[choice]; ok {
					return fmt.Errorf("double '%s'", choice)
				}
				seen[choice] = nil
			}

			return nil
		},
		Help: fmt.Sprintf("%s. {%s}", help, strings.Join(allowed, ",")),
	}
}

var optionToken = &argparse.Options{
	Required: false,
	Validate: tokenValidation,
	Help:     "Token for authentication with the InfraSonar API",
}

var optionContainer = &argparse.Options{
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
var optionAssetProperties = selectorList(
	false,
	[]string{"id", "name", "kind", "description", "mode", "labels", "collectors", "properties"},
	"Asset properties to return. If not specified all properties will be returned",
)

var optionAssetFilter = &argparse.Options{
	Required: false,
	Validate: func(args []string) error {
		for _, arg := range args {
			m := reAssetFilter.FindStringSubmatch(arg)
			if m == nil {
				return fmt.Errorf("invalid '%s'. valid example: -f kind==linux -f collector==snmp -f label!=123 -f zone=0", arg)
			}
			switch m[1] {
			case "collector", "kind":
				continue
			case "label", "zone":
				if !reNumber.MatchString(m[3]) {
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

var optionDefaultOutput = &argparse.Options{
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

var optionOutput = &argparse.Options{
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

var optionConfigName = &argparse.Options{
	Required: false,
	Validate: func(args []string) error {
		if !reConfigName.MatchString(args[0]) {
			return errors.New("invalid configuration name")
		}
		return nil
	},
	Help: "Configuration name",
}

var optionConfigNewApi = &argparse.Options{
	Required: false,
	Validate: func(args []string) error {
		if !reIsUrl.MatchString(args[0]) {
			return errors.New("invalid API URL")
		}
		return nil
	},
	Help:    "InfraSonar API url for the project",
	Default: "https://api.infrasonar.com",
}

func parseArgs() {
	parser := argparse.NewParser("infrasonar", "InfraSonar Client")

	/*
	 *  CMD: version
	 */
	cmdVersion := parser.NewCommand("version", "Print version and exit")
	/*
	 *  CMD: config
	 */
	cmdConfig := parser.NewCommand("config", "Manage client configurations")
	/*
	 *  CMD: config new
	 */
	cmdConfigNew := cmdConfig.NewCommand("new", "Create a new client configuration")
	cmdConfigNewSetName := cmdConfigNew.String("", "set-name", optionConfigName)
	cmdConfigNewSetToken := cmdConfigNew.String("", "set-token", optionToken)
	cmdConfigNewSetApi := cmdConfigNew.String("", "set-api", optionConfigNewApi)
	cmdConfigNewSetOutput := cmdConfigNew.String("", "set-output", optionDefaultOutput)
	/*
	 *  CMD: get
	 */
	cmdGet := parser.NewCommand("get", "Get InfraSonar data")
	cmdGetConfig := cmdGet.String("", "config", optionConfigName)
	/*
	 *  CMD: get assets
	 */
	cmdGetAssets := cmdGet.NewCommand("assets", "Get container assets")
	cmdGetAssetsContainer := cmdGetAssets.Int("c", "container", optionContainer)
	cmdGetAssetsProperties := cmdGetAssets.String("p", "properties", optionAssetProperties)
	cmdGetAssetsFilter := cmdGetAssets.StringList("f", "filter", optionAssetFilter)
	cmdGetAssetsOutput := cmdGetAssets.String("o", "output", optionOutput)
	/*
	 *  CMD: get kinds
	 */
	cmdGetAllAssetKinds := cmdGet.NewCommand("all-asset-kinds", "Get all available asset kinds")
	cmdGetAllAssetKindsOutput := cmdGetAllAssetKinds.String("o", "output", optionOutput)

	// Parse input
	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(nil))
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// Version does not require configs
	if cmdVersion.Happened() {
		handleVersion()
	}

	// Initialize configurations
	readConfigurations()

	if cmdConfig.Happened() {
		if cmdConfigNew.Happened() {
			handleConfigNew(&TConfigNew{
				name:   *cmdConfigNewSetName,
				token:  *cmdConfigNewSetToken,
				api:    *cmdConfigNewSetApi,
				output: *cmdConfigNewSetOutput,
			})
		}
	}

	if cmdGet.Happened() {
		config := configurations.ensureConfig(*cmdGetConfig)

		fmt.Println(config.GetToken())
		if cmdGetAssets.Happened() {
			handleGetAssets(&TGetAssets{
				config:     config,
				output:     getOutput(*cmdGetAssetsOutput, config),
				container:  *cmdGetAssetsContainer,
				properties: strings.Split(*cmdGetAssetsProperties, ","),
				filters:    *cmdGetAssetsFilter,
			})
		}
		if cmdGetAllAssetKinds.Happened() {
			handleGetAllAssetKinds(&TGetAllAssetKinds{
				config: config,
				output: getOutput(*cmdGetAllAssetKindsOutput, config),
			})
		}
	}
	fmt.Println(parser.Usage(nil))
}
