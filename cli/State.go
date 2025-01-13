package cli

type State struct {
	Container *Container     `json:"container" yaml:"container"`
	Labels    map[string]int `json:"labels,omitempty" yaml:"labels,omitempty"`
	Assets    []*AssetCli    `json:"assets" yaml:"assets"`
}
