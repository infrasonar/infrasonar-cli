package handle

import (
	"slices"

	"github.com/infrasonar/infrasonar-cli/cli"
	"github.com/infrasonar/infrasonar-cli/handle/util"
	"github.com/infrasonar/infrasonar-cli/req"
)

type TGetMe struct {
	Api        string
	Token      string
	Output     string
	OutFn      string
	Container  int
	Properties []string
}

type TGetMeOut struct {
	Me  *cli.Me `json:"token" yaml:"token"`
	cmd *TGetMe
}

func (o *TGetMeOut) Out() any {
	if len(o.cmd.Properties) == 1 && o.cmd.Properties[0] == "permissions" {
		return *o.Me.Permissions
	}
	if len(o.cmd.Properties) == 1 && o.cmd.Properties[0] == "tokenType" {
		return o.Me.TokenType
	}
	return o.Me
}

func GetMe(cmd *TGetMe) {
	util.Log(cmd.OutFn, "Get container...")
	container := util.EnsureContainer(cmd.Api, cmd.Token, cmd.Container)

	util.Log(cmd.OutFn, "Get token information...")
	me, err := req.GetMe(cmd.Api, cmd.Token, container.Id)
	util.ExitOnErr(err)

	if !slices.Contains(cmd.Properties, "permissions") {
		me.Permissions = nil
	}
	if !slices.Contains(cmd.Properties, "tokenType") {
		me.TokenType = ""
	}

	out := TGetMeOut{Me: me, cmd: cmd}
	util.ExitOutput(&out, cmd.Output, cmd.OutFn)
}
