package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"os"
	"path"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Name   string `yaml:"name"`
	Token  string `yaml:"token"`
	Api    string `yaml:"api"`
	Output string `yaml:"output"`
}

func (c *Config) GetToken() (string, error) {
	return DecryptAES(c.Token)
}

func (c *Config) EnsureToken() string {
	token, err := DecryptAES(c.Token)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read token from configuration '%s' (error: %s)\n", c.Name, err)
		os.Exit(1)
	}
	return token
}

type Configurations struct {
	Configs []Config `yaml:"configs"`
}

func (c *Configurations) get(name string) *Config {
	for _, config := range configurations.Configs {
		if config.Name == name {
			return &config
		}
	}
	return nil
}

func (c *Configurations) def() *Config {
	if len(c.Configs) == 0 {
		return nil
	}
	return &c.Configs[0]
}

func (c *Configurations) ensureConfig(name string) *Config {
	if name == "" {
		config := c.def()
		if config == nil {
			fmt.Fprintf(os.Stderr, "It appears no configuration has been set up.\nYou can create a new configuration using this command:\n\n  %s config new\n\n", os.Args[0])
			os.Exit(1)
		}
		return config
	}
	config := c.get(name)
	if config == nil {
		fmt.Fprintf(os.Stderr, "Configuration '%s' not found. Use this command to list the configurations:\n\n  %s config list\n\n", name, os.Args[0])
		os.Exit(1)
	}
	return config
}

func (c *Configurations) new(name, token, api, output string) (*Config, error) {
	if c.get(name) != nil {
		return nil, fmt.Errorf("configuration '%s' already exists", name)
	}

	token, err := EncryptAES(token)
	if err != nil {
		return nil, err
	}
	config := Config{
		Name:   name,
		Token:  token,
		Api:    api,
		Output: output,
	}
	c.Configs = append(c.Configs, config)
	err = c.write()
	return &config, err
}

func (c *Configurations) write() error {
	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	return os.WriteFile(configurationsFn, data, 0644)
}

var configurations = Configurations{}
var configurationsFn string

var cipherKey = []byte("AmH0lGOt7S07N0QrUwgMKjXNC0dxcJPZ")

func EncryptAES(text string) (string, error) {
	c, err := aes.NewCipher(cipherKey)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	out := gcm.Seal(nonce, nonce, []byte(text), nil)
	str := base64.StdEncoding.EncodeToString(out)
	return str, nil
}

func DecryptAES(b64 string) (string, error) {
	ct, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return "", err
	}

	c, err := aes.NewCipher(cipherKey)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(ct) < nonceSize {
		fmt.Println(err)
	}

	nonce, ct := ct[:nonceSize], ct[nonceSize:]
	text, err := gcm.Open(nil, nonce, ct, nil)
	if err != nil {
		return "", err
	}
	return string(text), nil
}

func readConfigurations() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read home path: %s\n", err)
		return
	}
	configurationsFn = path.Join(homeDir, ".infrasonar_cli_configs.yaml")
	if _, err := os.Stat(configurationsFn); errors.Is(err, os.ErrNotExist) {
		return
	}

	content, err := os.ReadFile(configurationsFn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read '%s': %s\n", configurationsFn, err)
		return
	}
	err = yaml.Unmarshal(content, &configurations)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to unpack '%s': %s\n", configurationsFn, err)
	}
}
