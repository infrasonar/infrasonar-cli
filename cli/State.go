package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"time"

	"gopkg.in/yaml.v3"
)

type State struct {
	Info      *Info             `json:"info,omitempty" yaml:"info,omitempty"`
	Container *Container        `json:"container" yaml:"container"`
	Zones     []*Zone           `json:"zones,omitempty" yaml:"zones,omitempty"`
	Labels    map[string]*Label `json:"labels,omitempty" yaml:"labels,omitempty"`
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

func StateFromCache(containerId int) *State {
	cliPath, err := CliPath()
	if err != nil {
		return nil
	}
	fn := path.Join(cliPath, fmt.Sprintf("cache_%09d.json", containerId))
	state, _ := StateFromFile(fn)
	return state
}

func (s *State) ClearCache() {
	cliPath, err := CliPath()
	if err != nil {
		return
	}
	fn := path.Join(cliPath, fmt.Sprintf("cache_%09d.json", s.Container.Id))
	os.Remove(fn)
}

func (s *State) makeLabelMap() {
	lm := NewLabelMap()
	for key, label := range s.Labels {
		lm.reverse[label.Id] = key
	}
	lm.labels = s.Labels
	s.labelMap = lm
}

func (s *State) GetLabelMap() *LabelMap {
	if s.labelMap == nil {
		s.makeLabelMap()
	}
	return s.labelMap
}

func (s *State) LabelById(labelId int) *Label {
	if s.labelMap == nil {
		s.makeLabelMap()
	}
	return s.labelMap.LabelById(labelId)
}

func (s *State) LabelByKey(key string) *Label {
	if s.labelMap == nil {
		s.makeLabelMap()
	}
	return s.labelMap.LabelByKey(key)
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

func (s *State) GetAge() (*time.Duration, error) {
	if s.Info == nil {
		return nil, errors.New("no time info for state")
	}
	return s.Info.GetAge()
}

func (s *State) WriteCache() {
	cliPath, err := CliPath()
	if err != nil {
		return
	}
	fn := path.Join(cliPath, fmt.Sprintf("cache_%09d.json", s.Container.Id))
	fp, err := os.OpenFile(fn, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer fp.Close()

	if out, err := json.Marshal(s); err == nil {
		fmt.Fprintln(fp, string(out[:]))
	}
}

func (s *State) HasCollector() bool {
	for _, a := range s.Assets {
		if a.Collectors != nil && len(*a.Collectors) > 0 {
			return true
		}
	}
	return false
}

func (s *State) HasAssetKind() bool {
	for _, a := range s.Assets {
		if a.Kind != "" {
			return true
		}
	}
	return false
}
