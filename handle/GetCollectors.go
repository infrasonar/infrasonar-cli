package handle

import (
	"github.com/infrasonar/infrasonar-cli/cli"
	"github.com/infrasonar/infrasonar-cli/handle/util"
	"github.com/infrasonar/infrasonar-cli/req"
)

type TGetCollectors struct {
	Api        string
	Token      string
	Output     string
	Container  int
	Properties []string
}

func GetCollectors(cmd *TGetCollectors) {
	container := util.EnsureContainer(cmd.Api, cmd.Token, cmd.Container)
	collectors, err := req.GetCollectors(cmd.Api, cmd.Token, container.Id, cmd.Properties, false)
	util.ExitOnErr(err)

	type Tout struct {
		Collectors []*cli.Collector `json:"collectors" yaml:"collectors"`
	}
	out := Tout{Collectors: collectors}
	util.ExitOutput(&out, cmd.Output)
}
