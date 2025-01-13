package cli

var AssetProperties = []string{"id", "name", "kind", "zone", "description", "mode", "labels", "collectors", "properties"}
var CollectorProperties = []string{"key", "name", "kind", "info", "minVersion", "checks"}

type IntSet map[int]struct{}

func (s IntSet) Set(k int) {
	s[k] = struct{}{}
}
