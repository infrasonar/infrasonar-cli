package req

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/infrasonar/infrasonar-cli/cli"
	"github.com/infrasonar/infrasonar-cli/re"
)

func GetAssetKinds(api string) ([]string, error) {
	if body, err := httpGet(fmt.Sprintf("%s/asset/kinds", api)); err != nil {
		return nil, err
	} else {
		var assetKinds []string
		err := json.Unmarshal(body, &assetKinds)
		return assetKinds, err
	}
}

func GetContainerId(api, token string) (int, error) {
	uri := fmt.Sprintf("%s/container/id", api)
	if body, err := httpGetAuth(uri, token); err != nil {
		return 0, err
	} else {
		type TContainerId struct {
			ContainerId int `json:"containerId"`
		}
		var unpack TContainerId
		err := json.Unmarshal(body, &unpack)
		if err != nil {
			return 0, err
		}
		return unpack.ContainerId, nil
	}
}

func GetContainer(api, token string, containerId int) (*cli.Container, error) {
	uri := fmt.Sprintf("%s/container/%d?fields=id,name", api, containerId)
	if body, err := httpGetAuth(uri, token); err != nil {
		return nil, err
	} else {
		var container cli.Container
		err := json.Unmarshal(body, &container)
		if err != nil {
			return nil, err
		}
		if container.Id != containerId {
			return nil, errors.New("container ID mismatch")
		}
		return &container, nil
	}
}

func GetAssets(api, token string, containerId, assetId int, fields, filters []string, withCollectors bool) ([]*cli.AssetApi, error) {
	if len(fields) == 0 {
		fields = []string{"id", "name"}
	}
	if assetId != 0 {
		fields = append(fields, "container")
	}
	args := strings.Join(fields, ",")
	args = fmt.Sprintf("?fields=%s", args)
	if withCollectors {
		args += ",disabledChecks&collectors=key,config"
	}
	for _, filter := range filters {
		m := re.AssetFilter.FindStringSubmatch(filter)
		if m == nil {
			continue
		}
		switch m[2] {
		case "==", "=":
			args += fmt.Sprintf("&%s=%s", m[1], m[3])
		case "!=":
			args += fmt.Sprintf("&not-%s=%s", m[1], m[3])
		}
	}
	if assetId != 0 {
		if len(filters) != 0 {
			return nil, errors.New("cannot use both filters (-f/--filter) and asset ID (-a/--asset)")
		}
		uri := fmt.Sprintf("%s/asset/%d%s", api, assetId, args)
		if body, err := httpGetAuth(uri, token); err != nil {
			return nil, err
		} else {
			var asset cli.AssetApi
			err := json.Unmarshal(body, &asset)
			if err != nil {
				return nil, err
			}
			if asset.ContainerId != containerId {
				return nil, fmt.Errorf("mismatch between container ID %d and asset ID %d", containerId, assetId)
			}
			// Just reset the container ID as it is no longer needed
			asset.ContainerId = 0
			assets := []*cli.AssetApi{&asset}
			return assets, nil
		}
	}
	uri := fmt.Sprintf("%s/container/%d/assets%s", api, containerId, args)
	if body, err := httpGetAuth(uri, token); err != nil {
		return nil, err
	} else {
		var assets []*cli.AssetApi
		err := json.Unmarshal(body, &assets)
		if err != nil {
			return nil, err
		}
		return assets, nil
	}
}

func GetCollectors(api, token string, containerId int) ([]*cli.Collector, error) {
	uri := fmt.Sprintf("%s/container/%d/collectors?field=key&options=key,type,default", api, containerId)
	if body, err := httpGetAuth(uri, token); err != nil {
		return nil, err
	} else {
		var collectors []*cli.Collector
		err := json.Unmarshal(body, &collectors)
		if err != nil {
			return nil, err
		}
		return collectors, nil
	}
}

func GetLabels(api, token string, labelIds cli.IntSet) (*cli.LabelMap, error) {
	labelMap := cli.NewLabelMap()
	for labelId := range labelIds {
		uri := fmt.Sprintf("%s/label/%d?fields=name", api, labelId)
		if body, err := httpGetAuth(uri, token); err != nil {
			return nil, err
		} else {
			type Tname struct {
				Name string `json:"name"`
			}
			var unpack Tname
			err := json.Unmarshal(body, &unpack)
			if err != nil {
				return nil, err
			}
			labelMap.Append(labelId, unpack.Name)
		}
	}
	return labelMap, nil
}
