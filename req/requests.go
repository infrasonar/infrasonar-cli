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
	if body, err := httpAuth("GET", uri, token); err != nil {
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
	if body, err := httpAuth("GET", uri, token); err != nil {
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
		fields = []string{"id"}
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
		if body, err := httpAuth("GET", uri, token); err != nil {
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
	if body, err := httpAuth("GET", uri, token); err != nil {
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

func GetCollectors(api, token string, containerId int, fields []string, withOptions bool) ([]*cli.Collector, error) {
	if len(fields) == 0 {
		fields = []string{"key"}
	}
	if len(fields) == 1 && fields[0] != "key" {
		fields = append(fields, "key")
	}

	args := strings.Join(fields, ",")
	args = fmt.Sprintf("?fields=%s", args)

	if withOptions {
		args += "&options=key,type,default"
	}

	uri := fmt.Sprintf("%s/container/%d/collectors%s", api, containerId, args)
	if body, err := httpAuth("GET", uri, token); err != nil {
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

func GetMe(api, token string, containerId int) (*cli.Me, error) {
	uri := fmt.Sprintf("%s/container/%d/permissions", api, containerId)
	if body, err := httpAuth("GET", uri, token); err != nil {
		return nil, err
	} else {
		var me cli.Me
		err := json.Unmarshal(body, &me)
		if err != nil {
			return nil, err
		}
		return &me, nil
	}
}

func GetLabels(api, token string, labelIds cli.IntSet) (*cli.LabelMap, error) {
	labelMap := cli.NewLabelMap()
	for labelId := range labelIds {
		uri := fmt.Sprintf("%s/label/%d?fields=id,name,color,description", api, labelId)
		if body, err := httpAuth("GET", uri, token); err != nil {
			return nil, fmt.Errorf("failed to retrieve label ID %d (%s)", labelId, err)
		} else {
			var label cli.Label
			err := json.Unmarshal(body, &label)
			if err != nil {
				return nil, fmt.Errorf("failed to unmarshal label ID %d (%s)", labelId, err)
			}
			labelMap.Append(&label)
		}
	}
	return labelMap, nil
}

func GetZones(api, token string, containerId int) ([]*cli.Zone, error) {
	uri := fmt.Sprintf("%s/container/%d/zones", api, containerId)
	if body, err := httpAuth("GET", uri, token); err != nil {
		return nil, err
	} else {
		var zones []*cli.Zone
		err := json.Unmarshal(body, &zones)
		if err != nil {
			return nil, err
		}
		return zones, nil
	}
}

func SetCollectorDisplay(api, token string, containerId int, collectorKey string, display bool) error {
	uri := fmt.Sprintf("%s/container/%d/collector/%s", api, containerId, collectorKey)
	type t struct {
		Display bool `json:"display"`
	}
	data := &t{
		Display: display,
	}
	x := func(display bool) string {
		if display {
			return "on"
		}
		return "off"
	}
	if _, err := httpJson("PATCH", uri, token, &data); err != nil {
		return fmt.Errorf("failed to set collector '%s' %s (%s)", collectorKey, x(display), err)
	}
	return nil
}

func UpsertZone(api, token string, containerId, zone int, name string) error {
	uri := fmt.Sprintf("%s/container/%d/zone", api, containerId)
	type t struct {
		Zone int    `json:"zone"`
		Name string `json:"name"`
	}
	data := &t{
		Zone: zone,
		Name: name,
	}
	if _, err := httpJson("POST", uri, token, &data); err != nil {
		return fmt.Errorf("failed to upsert zone ID %d (%s)", zone, err)
	}
	return nil
}

func CreateAsset(api, token string, containerId int, name string) (int, error) {
	uri := fmt.Sprintf("%s/container/%d/asset", api, containerId)
	type t struct {
		Name string `json:"name"`
	}
	data := &t{
		Name: name,
	}
	if body, err := httpJson("POST", uri, token, &data); err != nil {
		return 0, fmt.Errorf("failed to upsert asset '%s' (%s)", name, err)
	} else {
		type t struct {
			AssetId int `json:"assetId"`
		}
		var data t
		err := json.Unmarshal(body, &data)
		if err != nil {
			return 0, err
		}
		if data.AssetId == 0 {
			return 0, fmt.Errorf("unexpected asset ID 0 for asset '%s'", name)
		}
		return data.AssetId, nil
	}
}

func SetAssetKind(api, token string, assetId int, kind string) error {
	uri := fmt.Sprintf("%s/asset/%d/kind", api, assetId)
	type t struct {
		Kind string `json:"kind"`
	}
	data := &t{
		Kind: kind,
	}
	if _, err := httpJson("PATCH", uri, token, &data); err != nil {
		return fmt.Errorf("failed to set kind '%s' for asset ID %d' (%s)", kind, assetId, err)
	}
	return nil
}

func SetAssetName(api, token string, assetId int, name string) error {
	uri := fmt.Sprintf("%s/asset/%d/name", api, assetId)
	type t struct {
		Name string `json:"name"`
	}
	data := &t{
		Name: name,
	}
	if _, err := httpJson("PATCH", uri, token, &data); err != nil {
		return fmt.Errorf("failed to set name '%s' for asset ID %d' (%s)", name, assetId, err)
	}
	return nil
}

func SetAssetMode(api, token string, assetId int, mode string, duration *int) error {
	uri := fmt.Sprintf("%s/asset/%d/mode", api, assetId)
	type t struct {
		Mode     string `json:"mode"`
		Duration *int   `json:"duration,omitempty"`
	}
	data := &t{
		Mode:     mode,
		Duration: duration,
	}
	if _, err := httpJson("PATCH", uri, token, &data); err != nil {
		return fmt.Errorf("failed to set mode '%s' for asset ID %d' (%s)", mode, assetId, err)
	}
	return nil
}

func SetAssetZone(api, token string, assetId int, zoneId int) error {
	uri := fmt.Sprintf("%s/asset/%d/kind", api, assetId)
	type t struct {
		Zone int `json:"zone"`
	}
	data := &t{
		Zone: zoneId,
	}
	if _, err := httpJson("PATCH", uri, token, &data); err != nil {
		return fmt.Errorf("failed to set zone ID %d for asset ID %d' (%s)", zoneId, assetId, err)
	}
	return nil
}

func SetAssetDescription(api, token string, assetId int, description string) error {
	uri := fmt.Sprintf("%s/asset/%d/description", api, assetId)
	type t struct {
		Description string `json:"description"`
	}
	data := &t{
		Description: description,
	}
	if _, err := httpJson("PATCH", uri, token, &data); err != nil {
		return fmt.Errorf("failed to change the description for asset ID %d' (%s)", assetId, err)
	}
	return nil
}

func AddLabelToAsset(api, token string, assetId, labelId int) error {
	uri := fmt.Sprintf("%s/asset/%d/label/%d", api, assetId, labelId)
	if _, err := httpAuth("PUT", uri, token); err != nil {
		return fmt.Errorf("failed to add label ID %d to asset ID %d (%s)", labelId, assetId, err)
	}
	return nil
}

func DeleteLabelFromAsset(api, token string, assetId, labelId int) error {
	uri := fmt.Sprintf("%s/asset/%d/label/%d", api, assetId, labelId)
	if _, err := httpAuth("DELETE", uri, token); err != nil {
		return fmt.Errorf("failed to remove label ID %d from asset ID %d (%s)", labelId, assetId, err)
	}
	return nil
}

func EnableAssetCheck(api, token string, assetId int, collectorKey, checkKey string) error {
	uri := fmt.Sprintf("%s/asset/%d/collector/%s/check/%s", api, assetId, collectorKey, checkKey)
	if _, err := httpAuth("PUT", uri, token); err != nil {
		return fmt.Errorf("failed to enable check %s/%s on asset ID %d (%s)", collectorKey, checkKey, assetId, err)
	}
	return nil
}

func DisableAssetCheck(api, token string, assetId int, collectorKey, checkKey string) error {
	uri := fmt.Sprintf("%s/asset/%d/collector/%s/check/%s", api, assetId, collectorKey, checkKey)
	if _, err := httpAuth("PUT", uri, token); err != nil {
		return fmt.Errorf("failed to enable check %s/%s on asset ID %d (%s)", collectorKey, checkKey, assetId, err)
	}
	return nil
}

func UpsertCollectorToAsset(api, token string, assetId int, collectorKey string, config map[string]any) error {
	uri := fmt.Sprintf("%s/asset/%d/collector/%s", api, assetId, collectorKey)
	type t struct {
		Config map[string]any `json:"config,omitempty"`
	}
	data := &t{
		Config: config,
	}
	if _, err := httpJson("POST", uri, token, &data); err != nil {
		return fmt.Errorf("failed to upsert collector '%s' to asset ID %d (%s)", collectorKey, assetId, err)
	}
	return nil
}

func RemoveCollectorFromAsset(api, token string, assetId int, collectorKey string) error {
	uri := fmt.Sprintf("%s/asset/%d/collector/%s", api, assetId, collectorKey)
	if _, err := httpAuth("DELETE", uri, token); err != nil {
		return fmt.Errorf("failed to remove collector '%s' from asset ID %d (%s)", collectorKey, assetId, err)
	}
	return nil
}
