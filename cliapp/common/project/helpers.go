package project

import "github.com/goatcms/goatcore/filesystem"

// IsProjectFilespace check if filespace contains a goat project
func IsProjectFilespace(fs filesystem.Filespace) bool {
	return fs.IsDir(".goat")
}
