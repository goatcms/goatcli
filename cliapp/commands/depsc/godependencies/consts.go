package godependencies

import "regexp"

const (
	// MaxImportDepth is default value for max depths during recursive import of dependency
	MaxImportDepth = 404
)

var (
	// AlwaysIgnored is set of ignored strings
	AlwaysIgnored = []*regexp.Regexp{
		regexp.MustCompile("^(.*).golang.org$"),
		regexp.MustCompile("^(.*).golang.org/.*$"),
		regexp.MustCompile("^golang.org$"),
		regexp.MustCompile("^golang.org/.*$"),
	}
)
