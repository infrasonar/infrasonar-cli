package cli

import (
	"encoding/json"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type State struct {
	Container *Container        `json:"container" yaml:"container"`
	Labels    map[string]*Label `json:"labels,omitempty" yaml:"labels,omitempty"`
	Zones     []*Zone           `json:"zones" yaml:"zones"`
	Assets    []*AssetCli       `json:"assets" yaml:"assets"`

	// For internal use only
	labelMap *LabelMap
}

func StateFromFile(fn string) (*State, error) {
	data, err := os.ReadFile(fn)
	if err != nil {
		return nil, fmt.Errorf("failed to read '%s': %s", fn, err)
	}

	ext, err := GetJsonOrYaml(fn)
	if err != nil {
		return nil, err
	}

	var state State

	switch ext {
	case "yaml":
		err = yaml.Unmarshal(data, &state)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal YAML: %s", err)
		}
	case "json":
		err = json.Unmarshal(data, &state)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal JSON: %s", err)
		}
	}
	return &state, nil
}

func (s *State) makeLabelMap() {
	lm := NewLabelMap()
	for key, label := range s.Labels {
		lm.reverse[label.Id] = key
	}
	lm.labels = s.Labels
	s.labelMap = lm
}

func (s *State) LabelById(labelId int) *Label {
	if s.labelMap == nil {
		s.makeLabelMap()
	}
	return s.labelMap.LabelById(labelId)
}

func (s *State) AssetById(assetId int) *AssetCli {
	for _, a := range s.Assets {
		if a.Id == assetId {
			return a
		}
	}
	return nil
}

func (s *State) ZoneById(zoneId int) *Zone {
	for _, z := range s.Zones {
		if z.Zone == zoneId {
			return z
		}
	}
	return nil
}
