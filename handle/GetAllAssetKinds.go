package handle

import (
	"github.com/infrasonar/infrasonar-cli/handle/util"
	"github.com/infrasonar/infrasonar-cli/req"
)

type TGetAllAssetKinds struct {
	Api    string
	Output string
}

func GetAllAssetKinds(cmd *TGetAllAssetKinds) {
	assetKinds, err := req.GetAssetKinds(cmd.Api)
	util.ExitOnErr(err)
	util.ExitOutput(assetKinds, cmd.Output)
}
