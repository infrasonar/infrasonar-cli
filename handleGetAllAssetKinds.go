package main

type TGetAllAssetKinds struct {
	project *Project
	output  string
}

func handleGetAllAssetKinds(cmd *TGetAllAssetKinds) {
	assetKinds, err := getAssetKinds(cmd.project.api)
	exitOnErr(err)
	exitOutput(assetKinds, cmd.output)
}
