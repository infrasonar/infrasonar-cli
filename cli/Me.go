package cli

type Me struct {
	Permissions *[]string `json:"permissions,omitempty" yaml:"permissions,omitempty"`
	TokenType   string    `json:"tokenType,omitempty" yaml:"tokenType,omitempty"`
}
