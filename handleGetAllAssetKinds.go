package main

type TGetAllAssetKinds struct {
	config *Config
	output string
}

func handleGetAllAssetKinds(cmd *TGetAllAssetKinds) {
	assetKinds, err := getAssetKinds(cmd.config.Api)
	exitOnErr(err)
	exitOutput(assetKinds, cmd.output)
}
