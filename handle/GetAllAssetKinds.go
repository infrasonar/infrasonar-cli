package handle

import (
	"github.com/infrasonar/infrasonar-cli/handle/util"
	"github.com/infrasonar/infrasonar-cli/req"
)

type TGetAllAssetKinds struct {
	Api    string
	Output string
	OutFn  string
}

func GetAllAssetKinds(cmd *TGetAllAssetKinds) {
	util.Log(cmd.OutFn, "Get asset kinds..")
	assetKinds, err := req.GetAssetKinds(cmd.Api)
	util.ExitOnErr(err)
	util.ExitOutput(assetKinds, cmd.Output, cmd.OutFn)
}
