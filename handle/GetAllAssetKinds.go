package handle

import (
	"github.com/infrasonar/infrasonar-cli/handle/util"
	"github.com/infrasonar/infrasonar-cli/req"
)

func GetAllAssetKinds(api, output, outFn string) {
	util.Log(outFn, "Get asset kinds..")
	assetKinds, err := req.GetAssetKinds(api)
	util.ExitOnErr(err)
	util.ExitOutput(assetKinds, output, outFn)
}
