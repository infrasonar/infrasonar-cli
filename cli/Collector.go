package cli

type Collector struct {
	Key  string `json:"key,omitempty" yaml:"key,omitempty"`
	Name string `json:"name,omitempty" yaml:"name,omitempty"`
	Kind string `json:"kind,omitempty" yaml:"kind,omitempty"`
	Info string `json:"info,omitempty" yaml:"info,omitempty"`

	Checks  []string `json:"checks" yaml:"checks"`
	Options []struct {
		Key     string `json:"key"`
		Type    string `json:"type"`
		Default any    `json:"default"`
	} `json:"options,omitempty" yaml:"options,omitempty"`
}
