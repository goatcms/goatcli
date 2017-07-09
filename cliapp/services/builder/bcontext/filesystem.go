package bcontext

import "github.com/goatcms/goatcore/filesystem"

// Filesystem provide filesystem api
type Filesystem struct {
	fs filesystem.Filespace
}

// IsFile check if node exists and is a file
func (f *Filesystem) IsFile(path string) bool {
	return f.fs.IsFile(path)
}

// IsDir check if node  exists and is a dir
func (f *Filesystem) IsDir(path string) bool {
	return f.fs.IsDir(path)
}

// IsExist check if node exists
func (f *Filesystem) IsExist(path string) bool {
	return f.fs.IsExist(path)
}
