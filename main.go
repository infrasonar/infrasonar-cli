package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/akamensky/argparse"
	"github.com/infrasonar/infrasonar-cli/cli"
	"github.com/infrasonar/infrasonar-cli/conf"
	"github.com/infrasonar/infrasonar-cli/handle"
	"github.com/infrasonar/infrasonar-cli/install"
	"github.com/infrasonar/infrasonar-cli/options"
)

func getOutput(outputArg string, config *conf.Config) string {
	if outputArg == "" {
		return config.Output
	}
	return outputArg
}

func getProperties(properties string) []string {
	if properties == "" {
		return cli.AssetProperties
	}
	return strings.Split(properties, ",")
}

func main() {
	parser := argparse.NewParser("infrasonar", "InfraSonar Client")

	// CMD: version
	cmdVersion := parser.NewCommand("version", "Print version and exit")

	// CMD: install
	cmdInstall := parser.NewCommand("install", "Install the infrasonar client")

	// CMD: config
	cmdConfig := parser.NewCommand("config", "Manage client configurations")

	// CMD: config list
	cmdConfigList := cmdConfig.NewCommand("list", "List all configurations")
	cmdConfigListMore := cmdConfigList.Flag("m", "more", options.ConfigListMore)

	// CMD: config new
	cmdConfigNew := cmdConfig.NewCommand("new", "Create a new client configuration")
	cmdConfigNewSetName := cmdConfigNew.String("", "set-name", options.ConfigName)
	cmdConfigNewSetToken := cmdConfigNew.String("", "set-token", options.Token)
	cmdConfigNewSetApi := cmdConfigNew.String("", "set-api", options.ConfigNewApi)
	cmdConfigNewSetOutput := cmdConfigNew.String("", "set-output", options.DefaultOutput)

	// CMD: config update
	cmdConfigUpdate := cmdConfig.NewCommand("update", "Update a client configuration")
	cmdConfigUpdateName := cmdConfigUpdate.String("", "config", options.ConfigName)
	cmdConfigUpdateSetToken := cmdConfigUpdate.String("", "set-token", options.Token)
	cmdConfigUpdateSetApi := cmdConfigUpdate.String("", "set-api", options.ConfigUpdateApi)
	cmdConfigUpdateSetOutput := cmdConfigUpdate.String("", "set-output", options.Output)

	// CMD: config default
	cmdConfigDefault := cmdConfig.NewCommand("default", "Show the default client configuration")
	cmdConfigDefaultSet := cmdConfigDefault.String("s", "set", options.ConfigSetDefault)

	// CMD: config delete
	cmdConfigDelete := cmdConfig.NewCommand("delete", "Delete a client configuration")
	cmdConfigDeleteName := cmdConfigDelete.String("", "config", options.ConfigName)

	// CMD: get
	cmdGet := parser.NewCommand("get", "Get InfraSonar data")
	cmdGetConfig := cmdGet.String("", "config", options.ConfigName)

	// CMD: get assets
	cmdGetAssets := cmdGet.NewCommand("assets", "Get container assets")
	cmdGetAssetsContainer := cmdGetAssets.Int("c", "container", options.Container)
	cmdGetAssetsAsset := cmdGetAssets.Int("a", "asset", options.Container)
	cmdGetAssetsProperties := cmdGetAssets.String("p", "properties", options.AssetProperties)
	cmdGetAssetsFilter := cmdGetAssets.StringList("f", "filter", options.AssetFilter)
	cmdGetAssetsOutput := cmdGetAssets.String("o", "output", options.Output)
	cmdGetAssetsIncludeDefaults := cmdGetAssets.Flag("i", "include-defaults", options.IncludeDefaults)

	// CMD: get all-asset-kinds
	cmdGetAllAssetKinds := cmdGet.NewCommand("all-asset-kinds", "Get all available asset kinds")
	cmdGetAllAssetKindsOutput := cmdGetAllAssetKinds.String("o", "output", options.Output)

	// Parse input
	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(nil))
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// CMD: version
	if cmdVersion.Happened() {
		handle.Version()
	}

	// CMD: install
	if cmdInstall.Happened() {
		install.Install()
	}

	// Initialize configurations
	conf.Initialize()

	// CMD: config
	if cmdConfig.Happened() {
		// CMD: config list
		if cmdConfigList.Happened() {
			handle.ConfigList(*cmdConfigListMore)
		}

		// CMD: config new
		if cmdConfigNew.Happened() {
			handle.ConfigNew(&handle.TConfigNew{
				Name:   *cmdConfigNewSetName,
				Token:  *cmdConfigNewSetToken,
				Api:    *cmdConfigNewSetApi,
				Output: *cmdConfigNewSetOutput,
			})
		}

		// CMD: config update
		if cmdConfigUpdate.Happened() {
			handle.ConfigUpdate(&handle.TConfigUpdate{
				Name:   *cmdConfigUpdateName,
				Token:  *cmdConfigUpdateSetToken,
				Api:    *cmdConfigUpdateSetApi,
				Output: *cmdConfigUpdateSetOutput,
			})
		}

		// CMD: config default
		if cmdConfigDefault.Happened() {
			handle.ConfigDefault(*cmdConfigDefaultSet)
		}

		// CMD: config delete
		if cmdConfigDelete.Happened() {
			handle.ConfigDelete(*cmdConfigDeleteName)
		}
	}

	// CMD: get
	if cmdGet.Happened() {
		config := conf.EnsureConfig(*cmdGetConfig)

		// CMD: get assets
		if cmdGetAssets.Happened() {
			handle.GetAssets(&handle.TGetAssets{
				Api:             config.Api,
				Token:           config.EnsureToken(),
				Output:          getOutput(*cmdGetAssetsOutput, config),
				Container:       *cmdGetAssetsContainer,
				Asset:           *cmdGetAssetsAsset,
				Properties:      getProperties(*cmdGetAssetsProperties),
				Filters:         *cmdGetAssetsFilter,
				IncludeDefaults: *cmdGetAssetsIncludeDefaults,
			})
		}

		// CMD: get all-asset-kinds
		if cmdGetAllAssetKinds.Happened() {
			handle.GetAllAssetKinds(&handle.TGetAllAssetKinds{
				Api:    config.Api,
				Output: getOutput(*cmdGetAllAssetKindsOutput, config),
			})
		}
	}
	fmt.Println(parser.Usage(nil))
}
