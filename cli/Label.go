package cli

import "fmt"

type Label struct {
	Id          int    `json:"id" yaml:"id"`
	Name        string `json:"name,omitempty" yaml:"name,omitempty"`
	Color       string `json:"color,omitempty" yaml:"color,omitempty"`
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
}

func (label *Label) Str() string {
	if label.Name == "" {
		return fmt.Sprintf("%d", label.Id)
	}
	return label.Name
}

var DefaultColor = "Steel"

var DefaultLabel = Label{
	Id:          0,
	Name:        "",
	Color:       DefaultColor,
	Description: "",
}
