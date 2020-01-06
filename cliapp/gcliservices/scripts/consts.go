package scripts

import "regexp"

const (
	scriptsBasePath = ".goat/scripts/"
)

var (
	scriptNamePattern = regexp.MustCompile("^[a-zA-Z_]+$")
)
