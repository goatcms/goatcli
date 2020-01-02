package vcs

import "regexp"

const (
	// DataDirectoryPath is a directory contains data
	DataDirectoryPath = ".goat/vcs"
	// IgnoredFilesPath is path to file contains ignored files list
	IgnoredFilesPath = ".goat/vcs/ignored"
	// GeneratedFilesPath is path to file contains generated files list
	GeneratedFilesPath = ".goat/vcs/generated"
)

var (
	emptyRegexp = regexp.MustCompile("^\\s*$")
)
