package re

import "regexp"

var AssetFilter = regexp.MustCompile(`^(\w+)(\=\=|\!\=|\=)(\w+)$`)
var Number = regexp.MustCompile(`^[0-9]+$`)
var IsUrl = regexp.MustCompile(`^https?://\S+$`)
var Token = regexp.MustCompile(`^[0-9a-f]{32}$`)
var ConfigName = regexp.MustCompile(`^[a-zA-Z_]\w*$`)
