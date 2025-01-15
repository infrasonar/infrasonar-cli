package handle

import (
	"fmt"
	"os"

	"github.com/infrasonar/infrasonar-cli/cli"
	"github.com/infrasonar/infrasonar-cli/handle/util"
	"github.com/infrasonar/infrasonar-cli/req"
)

type TGetCollectors struct {
	Api        string
	Token      string
	Output     string
	OutFn      string
	Container  int
	Properties []string
	Collector  string
}

type TGetCollectorsOut struct {
	Collectors []*cli.Collector `json:"collectors" yaml:"collectors"`
	cmd        *TGetCollectors
}

func (o *TGetCollectorsOut) Out() any {
	if len(o.cmd.Properties) == 1 && o.cmd.Properties[0] == "key" {
		out := []string{}
		for _, c := range o.Collectors {
			out = append(out, c.Key)
		}
		return out
	}
	if len(o.cmd.Properties) == 1 && o.cmd.Properties[0] == "checks" && len(o.Collectors) == 1 {
		out := []string{}
		out = append(out, o.Collectors[0].Checks...)
		return out
	}
	return o.Collectors
}

func GetCollectors(cmd *TGetCollectors) {
	util.Log(cmd.OutFn, "Get container...")
	container := util.EnsureContainer(cmd.Api, cmd.Token, cmd.Container)

	util.Log(cmd.OutFn, "Get collectors...")
	collectors, err := req.GetCollectors(cmd.Api, cmd.Token, container.Id, cmd.Properties, false)
	util.ExitOnErr(err)

	if cmd.Collector != "" {
		out := []*cli.Collector{}
		for _, c := range collectors {
			if c.Key == cmd.Collector {
				out = append(out, c)
			}
		}
		if len(out) == 0 {
			fmt.Fprintf(os.Stderr, "collector '%s' not found\n", cmd.Collector)
			os.Exit(1)
		}
		collectors = out
	}

	out := TGetCollectorsOut{Collectors: collectors, cmd: cmd}
	util.ExitOutput(&out, cmd.Output, cmd.OutFn)
}
