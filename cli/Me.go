package cli

import (
	"errors"
	"fmt"
	"slices"
	"strings"
)

type Me struct {
	Permissions *[]string `json:"permissions,omitempty" yaml:"permissions,omitempty"`
	TokenType   string    `json:"tokenType,omitempty" yaml:"tokenType,omitempty"`
}

func (m *Me) CheckApplyPermissions() error {
	if m.Permissions == nil {
		return errors.New("permissions missing")
	}
	missing := []string{}
	for _, required := range []string{
		"API",
		"ASSET_MANAGEMENT",
		"CHECK_MANAGEMENT",
		"CONTAINER_ADMIN",
		"CONTAINER_MANAGEMENT",
		"READ",
	} {
		if !slices.Contains(*m.Permissions, required) {
			missing = append(missing, required)
		}
	}
	if len(missing) > 0 {
		return fmt.Errorf("token is missing the following permissions:\n\n- %s", strings.Join(missing, "\n- "))
	}
	return nil
}
