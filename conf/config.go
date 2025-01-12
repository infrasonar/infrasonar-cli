package conf

import (
	"fmt"
	"os"
)

type Config struct {
	Name     string `yaml:"name"`
	EncToken string `yaml:"token"`
	Api      string `yaml:"api"`
	Output   string `yaml:"output"`
}

func (c *Config) GetToken() (string, error) {
	return decryptAES(c.EncToken)
}

func (c *Config) SetToken(token string) error {
	encToken, err := encryptAES(token)
	if err != nil {
		return err
	}
	c.EncToken = encToken
	return nil
}

func (c *Config) EnsureToken() string {
	token, err := decryptAES(c.EncToken)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read token from configuration '%s' (error: %s)\n", c.Name, err)
		os.Exit(1)
	}
	return token
}
