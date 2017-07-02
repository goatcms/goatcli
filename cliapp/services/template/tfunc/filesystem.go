package tfunc

import "github.com/goatcms/goatcore/filesystem"

// IsFile is file check if node exist and is a file
func IsFile(fs filesystem.Filespace, path string) bool {
	return fs.IsFile(path)
}

// IsDir is file check if node exist and is a directory
func IsDir(fs filesystem.Filespace, path string) bool {
	return fs.IsDir(path)
}

// IsExist is file check if node exist
func IsExist(fs filesystem.Filespace, path string) bool {
	return fs.IsExist(path)
}
