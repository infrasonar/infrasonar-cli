package cli

type Label struct {
	Id    int    `json:"id" yaml:"id"`
	Name  string `json:"name,omitempty" yaml:"name,omitempty"`
	Color string `json:"color,omitempty" yaml:"color,omitempty"`
}
