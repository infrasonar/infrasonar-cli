package main

import (
	"encoding/json"
	"fmt"
)

func getAssetKinds(api string) ([]string, error) {
	if body, err := httpGet(fmt.Sprintf("%s/asset/kinds", api)); err != nil {
		return nil, err
	} else {
		var assetKinds []string
		err := json.Unmarshal(body, &assetKinds)
		return assetKinds, err
	}
}
