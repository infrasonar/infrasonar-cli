package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/akamensky/argparse"
	"github.com/infrasonar/infrasonar-cli/cli"
	"github.com/infrasonar/infrasonar-cli/conf"
	"github.com/infrasonar/infrasonar-cli/handle"
	"github.com/infrasonar/infrasonar-cli/handle/util"
	"github.com/infrasonar/infrasonar-cli/install"
	"github.com/infrasonar/infrasonar-cli/options"
)

func getOutput(outputArg string, config *conf.Config) string {
	if outputArg == "" {
		return config.Output
	}
	return outputArg
}

func testOutputFilename(fn string) error {
	if fn != "" {
		fp, err := os.OpenFile(fn, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return fmt.Errorf("failed to create outpuf file: %s", err)
		}
		fp.Close()
	}
	return nil
}

func getAssetProperties(properties string) []string {
	if properties == "" {
		return cli.AssetProperties
	}
	return strings.Split(properties, ",")
}

func getCollectorProperties(properties string) []string {
	if properties == "" {
		return cli.CollectorProperties
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
	cmdGetOutput := cmdGet.String("o", "output", options.Output)
	cmdGetOutputFilename := cmdGet.String("t", "output-filename", options.OutFileName)
	cmdGetUseConfig := cmdGet.String("u", "use-config", options.UseConfig)

	// CMD: get assets
	cmdGetAssets := cmdGet.NewCommand("assets", "Get container assets")
	cmdGetAssetsContainer := cmdGetAssets.Int("c", "container", options.Container)
	cmdGetAssetsAsset := cmdGetAssets.Int("a", "asset", options.Container)
	cmdGetAssetsProperties := cmdGetAssets.String("p", "properties", options.AssetProperties)
	cmdGetAssetsFilter := cmdGetAssets.StringList("f", "filter", options.AssetFilter)
	cmdGetAssetsIncludeDefaults := cmdGetAssets.Flag("i", "include-defaults", options.IncludeDefaults)

	// CMD: get collectors
	cmdGetCollectors := cmdGet.NewCommand("collectors", "Get container collectors")
	cmdGetCollectorsContainer := cmdGetCollectors.Int("c", "container", options.Container)
	cmdGetCollectorsProperties := cmdGetCollectors.String("p", "properties", options.CollectorProperties)
	cmdGetCollectorsCollector := cmdGetCollectors.String("k", "collector", options.Collector)

	// CMD: get all-asset-kinds
	cmdGetAllAssetKinds := cmdGet.NewCommand("all-asset-kinds", "Get all available asset kinds")

	// CMD: apply
	cmdApply := parser.NewCommand("apply", "Apply InfraSonar data from YAML or JSON file")
	cmdApplyFileName := cmdApply.String("f", "filename", options.ApplyFileName)
	cmdApplyDryRun := cmdApply.Flag("d", "dry-run", options.DryRun)
	cmdApplyNoRemove := cmdApply.Flag("n", "no-remove", options.NoRemove)
	cmdApplyUseConfig := cmdApply.String("u", "use-config", options.UseConfig)

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
	if err := conf.Initialize(); err != nil {
		util.ExitOnErr(err)
	}

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
		config := conf.EnsureConfig(*cmdGetUseConfig)
		output := getOutput(*cmdGetOutput, config)
		outFn := *cmdGetOutputFilename
		util.ExitOnErr(testOutputFilename(outFn))

		// CMD: get assets
		if cmdGetAssets.Happened() {
			handle.GetAssets(&handle.TGetAssets{
				Api:             config.Api,
				Token:           config.EnsureToken(),
				Output:          output,
				OutFn:           outFn,
				Container:       *cmdGetAssetsContainer,
				Asset:           *cmdGetAssetsAsset,
				Properties:      getAssetProperties(*cmdGetAssetsProperties),
				Filters:         *cmdGetAssetsFilter,
				IncludeDefaults: *cmdGetAssetsIncludeDefaults,
			})
		}

		// CMD: get collectors
		if cmdGetCollectors.Happened() {
			handle.GetCollectors(&handle.TGetCollectors{
				Api:        config.Api,
				Token:      config.EnsureToken(),
				Output:     output,
				OutFn:      outFn,
				Container:  *cmdGetCollectorsContainer,
				Properties: getCollectorProperties(*cmdGetCollectorsProperties),
				Collector:  *cmdGetCollectorsCollector,
			})
		}

		// CMD: get all-asset-kinds
		if cmdGetAllAssetKinds.Happened() {
			handle.GetAllAssetKinds(&handle.TGetAllAssetKinds{
				Api:    config.Api,
				Output: output,
				OutFn:  outFn,
			})
		}
	}

	// CMD: apply
	if cmdApply.Happened() {
		config := conf.EnsureConfig(*cmdApplyUseConfig)
		api := config.Api
		token := config.EnsureToken()

		handle.Apply(
			api,
			token,
			*cmdApplyFileName,
			*cmdApplyDryRun,
			*cmdApplyNoRemove,
		)
	}
	fmt.Println(parser.Usage(nil))
}
