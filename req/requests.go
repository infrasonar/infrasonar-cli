package req

import (
	"encoding/json"
	"fmt"
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
	if body, err := httpGetAuth(fmt.Sprintf("%s/container/id", api), token); err != nil {
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
