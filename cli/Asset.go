package cli

type TDisabledChecks struct {
	Collector string `json:"collector" yaml:"collector"`
	Check     string `json:"check" yaml:"check"`
}

type TCollector struct {
	Key    string         `json:"key" yaml:"key"`
	Config map[string]any `json:"config,omitempty" yaml:"config,omitempty"`
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
}

type AssetCli struct {
	Id             int               `json:"id,omitempty" yaml:"id,omitempty"`
	Name           string            `json:"name,omitempty" yaml:"name,omitempty"`
	Zone           *int              `json:"zone,omitempty" yaml:"zone,omitempty"`
	Labels         []string          `json:"labels,omitempty" yaml:"labels,omitempty"`
	Description    string            `json:"description,omitempty" yaml:"description,omitempty"`
	Mode           string            `json:"mode,omitempty" yaml:"mode,omitempty"`
	Kind           string            `json:"kind,omitempty" yaml:"kind,omitempty"`
	DisabledChecks []TDisabledChecks `json:"disabledChecks,omitempty" yaml:"disabledChecks,omitempty"`
	Collectors     []TCollector      `json:"collectors,omitempty" yaml:"collectors,omitempty"`
}
