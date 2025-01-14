package cli

import "fmt"

type Zone struct {
	Zone int    `json:"zone" yaml:"zone"`
	Name string `json:"name" yaml:"name"`
}

func (z *Zone) Str() string {
	if z.Name == "" {
		return fmt.Sprintf("%d", z.Zone)
	}
	return z.Name
}
