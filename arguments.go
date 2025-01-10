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

var reToken = regexp.MustCompile(`^[0-9a-f]{32}$`)
var tokenValidation = func(args []string) error {
	if !reToken.MatchString(args[0]) {
		return errors.New("invalid token")
	}
	return nil
}

func askToken() string {
	fmt.Print("Token: ")
	pass, err := gopass.GetPasswdMasked()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	token := string(pass[:])
	if reToken.MatchString(token) {
		return token
	} else {
		fmt.Println("Invalid token, please enter a correct token")
		return askToken()
	}
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
	Required: true,
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
var optionAssetFields = selectorList(
	false,
	[]string{"id", "name", "kind"},
	"Fields to return. If not specified all fields will be returned",
)

func parseArgs() error {
	parser := argparse.NewParser("infrasonar", "InfraSonar Client")

	cmdVersion := parser.NewCommand("version", "Print version and exit")
	cmdGetAssets := parser.NewCommand("get-assets", "Get container assets")

	/* required */
	getAssetsContainerId := cmdGetAssets.Int("i", "container-id", optionContainerId)

	/* optional */
	getAssetsToken := cmdGetAssets.String("t", "token", optionToken)
	getAssetsFields := cmdGetAssets.String("f", "fields", optionAssetFields)

	// Parse input
	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(""))
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// Print version
	if cmdVersion.Happened() {
		fmt.Printf("InfraSonar version %s\n", Version)
		os.Exit(0)
	}

	if cmdGetAssets.Happened() {
		token := ensureToken(*getAssetsToken)
		containerId := *getAssetsContainerId
		fields := strings.Split(*getAssetsFields, ",")

		fmt.Println(token)
		fmt.Println(containerId)
		fmt.Println(fields)
	}

	return nil
}
