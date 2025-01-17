package cli

import "fmt"

type TDisabledChecks struct {
	Collector string `json:"collector" yaml:"collector"`
	Check     string `json:"check" yaml:"check"`
}

type TCollector struct {
	Key    string         `json:"key" yaml:"key"`
	Config map[string]any `json:"config,omitempty" yaml:"config,omitempty"`
}

type TProperty struct {
	Key   string `json:"key" yaml:"key"`
	Value any    `json:"value" yaml:"value"`
}

type AssetApi struct {
	Id             int               `json:"id"`
	ContainerId    int               `json:"container"`
	Name           string            `json:"name"`
	Zone           *int              `json:"zone"`
	Labels         []int             `json:"labels"`
	Description    string            `json:"description"`
	Mode           string            `json:"mode"`
	Kind           string            `json:"kind"`
	DisabledChecks []TDisabledChecks `json:"disabledChecks"`
	Collectors     []TCollector      `json:"collectors"`
	Properties     []TProperty       `json:"properties"`
}

type AssetCli struct {
	Id             int                `json:"id,omitempty" yaml:"id,omitempty"`
	Name           string             `json:"name,omitempty" yaml:"name,omitempty"`
	Zone           *int               `json:"zone,omitempty" yaml:"zone,omitempty"`
	Labels         *[]string          `json:"labels,omitempty" yaml:"labels,omitempty"`
	Description    string             `json:"description,omitempty" yaml:"description,omitempty"`
	Mode           string             `json:"mode,omitempty" yaml:"mode,omitempty"`
	Kind           string             `json:"kind,omitempty" yaml:"kind,omitempty"`
	DisabledChecks *[]TDisabledChecks `json:"disabledChecks,omitempty" yaml:"disabledChecks,omitempty"`
	Collectors     *[]TCollector      `json:"collectors,omitempty" yaml:"collectors,omitempty"`
	Properties     *[]TProperty       `json:"properties,omitempty" yaml:"properties,omitempty"`
}

func (a *AssetCli) Str() string {
	if a.Name == "" {
		return fmt.Sprintf("%d", a.Id)
	}
	return a.Name
}

func (a *AssetCli) HasLabelId(labelId int, lm *LabelMap) bool {
	for _, key := range *a.Labels {
		if label := lm.LabelByKey(key); label != nil {
			if label.Id == labelId {
				return true
			}
		}
	}
	return false
}

var DefaultZone = 0

var DefaultAsset = AssetCli{
	Id:             0,
	Name:           "",
	Zone:           &DefaultZone,
	Labels:         &[]string{},
	Description:    "",
	Mode:           "normal",
	Kind:           "Asset",
	DisabledChecks: &[]TDisabledChecks{},
	Collectors:     &[]TCollector{},
}
