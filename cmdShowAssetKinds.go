package main

func showAssetKinds(api string, output string) {
	assetKinds, err := getAssetKinds(api)
	exitOnErr(err)
	exitOutput(assetKinds, output)
}
