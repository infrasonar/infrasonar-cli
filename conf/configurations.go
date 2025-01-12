package conf

type Configurations struct {
	Configs []*Config `yaml:"configs"`
}

func (c *Configurations) get(name string) *Config {
	for _, config := range c.Configs {
		if config.Name == name {
			return config
		}
	}
	return nil
}
