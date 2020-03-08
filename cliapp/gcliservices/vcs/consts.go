package vcs

import "regexp"

const (
	// DataDirectoryPath is a directory contains data
	DataDirectoryPath = ".goat/vcs"
	// PersistedFilesPath is path to file contains persisted files list
	PersistedFilesPath = ".goat/vcs/persisted"
	// GeneratedFilesPath is path to file contains generated files list
	GeneratedFilesPath = ".goat/vcs/generated"
)

var (
	emptyRegexp = regexp.MustCompile("^\\s*$")
)
