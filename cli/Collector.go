package cli

type Collector struct {
	Key     string `json:"key"`
	Options []struct {
		Key     string `json:"key"`
		Default any    `json:"default"`
		Type    string `json:"type"`
	} `json:"options"`
}
