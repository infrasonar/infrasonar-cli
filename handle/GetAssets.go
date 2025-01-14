package handle

import (
	"reflect"

	"github.com/infrasonar/infrasonar-cli/cli"
	"github.com/infrasonar/infrasonar-cli/handle/util"
	"github.com/infrasonar/infrasonar-cli/req"
)

type TGetAssets struct {
	Api             string
	Token           string
	Output          string
	OutFn           string
	Container       int
	Asset           int
	Properties      []string
	Filters         []string
	IncludeDefaults bool
}

func getCollector(collectors []*cli.Collector, key string) *cli.Collector {
	for _, collector := range collectors {
		if collector.Key == key {
			return collector
		}
	}
	return nil
}

func removeDefaults(assets []*cli.AssetApi, collectors []*cli.Collector) {
	for _, asset := range assets {
		for _, c := range asset.Collectors {
			collector := getCollector(collectors, c.Key)
			if collector == nil {
				continue
			}
			toDelete := []string{}
			for k, v := range c.Config {
				for _, o := range collector.Options {
					if o.Key == k && reflect.DeepEqual(o.Default, v) {
						toDelete = append(toDelete, k)
					}
				}
			}
			for _, k := range toDelete {
				delete(c.Config, k)
			}
		}
	}
}

func replaceUse(assets []*cli.AssetApi) {
	for _, asset := range assets {
		for _, c := range asset.Collectors {
			if v, ok := c.Config["_use"]; ok {
				delete(c.Config, "_use")
				c.Config["use"] = v
			}
		}
	}
}

func getLabelMap(api, token string, assets []*cli.AssetApi) (*cli.LabelMap, error) {
	labelsIds := cli.IntSet{}
	for _, asset := range assets {
		for _, labelId := range asset.Labels {
			labelsIds.Set(labelId)
		}
	}
	return req.GetLabels(api, token, labelsIds)
}

func getAssetsCli(assets []*cli.AssetApi, labelMap *cli.LabelMap) []*cli.AssetCli {
	m := []*cli.AssetCli{}
	for _, a := range assets {
		labels := []string{}
		for _, labelId := range a.Labels {
			labels = append(labels, labelMap.GetName(labelId))
		}
		m = append(m, &cli.AssetCli{
			Id:             a.Id,
			Name:           a.Name,
			Zone:           a.Zone,
			Labels:         labels,
			Description:    a.Description,
			Mode:           a.Mode,
			Kind:           a.Kind,
			DisabledChecks: a.DisabledChecks,
			Collectors:     a.Collectors,
		})
	}
	return m
}

func ensureState(cmd *TGetAssets) *cli.State {
	state := cli.State{}
	state.Info = cli.NewInfo()

	util.Log(cmd.OutFn, "Get container...")
	container := util.EnsureContainer(cmd.Api, cmd.Token, cmd.Container)
	withCollectors := util.Itob(util.RemoveFromSlice(&cmd.Properties, "collectors"))

	util.Log(cmd.OutFn, "Get assets...")
	assets, err := req.GetAssets(cmd.Api, cmd.Token, container.Id, cmd.Asset, cmd.Properties, cmd.Filters, withCollectors)
	util.ExitOnErr(err)

	util.Log(cmd.OutFn, "Get zones...")
	zones, err := req.GetZones(cmd.Api, cmd.Token, container.Id)
	util.ExitOnErr(err)

	if withCollectors && !cmd.IncludeDefaults {
		util.Log(cmd.OutFn, "Get collector info...")
		collectors, err := req.GetCollectors(cmd.Api, cmd.Token, container.Id, []string{"key"}, true)
		util.ExitOnErr(err)
		removeDefaults(assets, collectors)
	}

	replaceUse(assets)

	util.Log(cmd.OutFn, "Get labels...")
	labelMap, err := getLabelMap(cmd.Api, cmd.Token, assets)
	util.ExitOnErr(err)

	state.Container = container
	state.Labels = labelMap.Labels()
	state.Zones = zones
	state.Assets = getAssetsCli(assets, labelMap)

	return &state
}

func GetAssets(cmd *TGetAssets) {
	state := ensureState(cmd)
	util.ExitOutput(state, cmd.Output, cmd.OutFn)
}
