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
	"github.com/howeyc/gopass"
)

var reAssetFilter = regexp.MustCompile(`^(\w+)(\=\=|\!\=)(\w+)$`)
var reNumber = regexp.MustCompile(`^[0-9]+$`)
var reIsUrl = regexp.MustCompile(`^https?://\S+$`)
var reToken = regexp.MustCompile(`^[0-9a-f]{32}$`)
var reProjectName = regexp.MustCompile(`^[a-zA-Z_]\w*$`)
var tokenValidation = func(args []string) error {
	if !reToken.MatchString(args[0]) {
		return errors.New("invalid token")
	}
	return nil
}



func ensureToken(token string) string {
	if reToken.MatchString(token) {
		return token
	}
	return askToken()
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

var optionContainerId = &argparse.Options{
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

var optionProjectName = &argparse.Options{
	Required: false,
	Validate: func(args []string) error {
		if !reProjectName.MatchString(args[0]) {
			return errors.New("invalid project name")
		}
		return nil
	},
	Help: "Project name",
}

var optionProjectNewApi = &argparse.Options{
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
	 *  CMD: projects
	 */
	cmdProjects := parser.NewCommand("projects", "Manage projects")
	/*
	 *  CMD: projects new
	 */
	cmdProjectsNew := cmdProjects.NewCommand("new", "Create a new project")
	cmdProjectsNewSetName := cmdProjectsNew.String("", "set-name", optionProjectName)
	cmdProjectsNewSetToken := cmdProjectsNew.String("", "set-token", optionToken)
	cmdProjectsNewSetApi := cmdProjectsNew.String("", "set-api", optionProjectNewApi)
	cmdProjectsNewSetOutput := cmdProjectsNew.String("", "set-output", optionDefaultOutput)

	/*
	 *  CMD: get
	 */
	cmdGet := parser.NewCommand("get", "Get InfraSonar data")
	/*
	 *  CMD: get assets
	 */
	cmdGetAssets := cmdGet.NewCommand("assets", "Get container assets")
	cmdGetAssetsContainerId := cmdGetAssets.Int("c", "container-id", optionContainerId)
	cmdGetAssetsToken := cmdGetAssets.String("t", "token", optionToken)
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
		fmt.Print(parser.Usage("HERE!!"))
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if cmdVersion.Happened() {
		handleVersion()
	}
	if cmdGet.Happened() {
		if cmdGetAssets.Happened() {
			handleGetAssets(&TGetAssets{
				output:      *cmdGetAssetsOutput,
				containerId: *cmdGetAssetsContainerId,
				properties:  strings.Split(*cmdGetAssetsProperties, ","),
				filters:     *cmdGetAssetsFilter,
			})
		}
		if cmdGetAllAssetKinds.Happened() {
			handleGetAllAssetKinds(&TGetAllAssetKinds{
				api:    *api,
				output: *cmdGetAllAssetKindsOutput,
			})
		}
	}
	fmt.Println(parser.Usage(nil))
}
