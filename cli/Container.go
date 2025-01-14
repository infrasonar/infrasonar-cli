package cli

import "fmt"

type Container struct {
	Id   int    `json:"id" yaml:"id"`
	Name string `json:"name" yaml:"name"`
}

func (c *Container) Str() string {
	if c.Name == "" {
		return fmt.Sprintf("%d", c.Id)
	}
	return c.Name
}
