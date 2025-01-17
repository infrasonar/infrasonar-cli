package handle

import (
	"github.com/infrasonar/infrasonar-cli/handle/util"
	"github.com/infrasonar/infrasonar-cli/req"
)

func GetAllLabelsColors(api, output, outFn string) {
	util.Log(outFn, "Get label colors..")
	assetKinds, err := req.GetLabelColors(api)
	util.ExitOnErr(err)
	util.ExitOutput(assetKinds, output, outFn)
}
