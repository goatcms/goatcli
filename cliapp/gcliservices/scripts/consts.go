package scripts

import "regexp"

const (
	scriptsBasePath    = ".goat/scripts/"
	historyInstanceKey = "__scriptHistory"
)

var (
	scriptNamePattern = regexp.MustCompile("^[a-zA-Z_]+$")
)
