package gcliservices

import (
	"io"
	"time"

	"github.com/goatcms/goatcore/filesystem"
)

// GeneratedFile contains single generated file data
type GeneratedFile struct {
	Path    string
	ModTime time.Time
}

// VCSGeneratedFiles contains generated files
type VCSGeneratedFiles interface {
	All() []*GeneratedFile
	New() map[string]*GeneratedFile
	Get(path string) *GeneratedFile
	ContainsPath(path string) bool
	Add(row *GeneratedFile)
	WriteAll(writer io.Writer) (err error)
	Modified() bool
}

// VCSPersistedFiles contains persisted files
type VCSPersistedFiles interface {
	All() []string
	ContainsPath(path string) bool
	AddPath(path string)
	RemovePath(path string)
	WriteAll(writer io.Writer) (err error)
	Modified() bool
}

// VCSData contains vcs data like generated and persisted files list
type VCSData interface {
	VCSGeneratedFiles() VCSGeneratedFiles
	VCSPersistedFiles() VCSPersistedFiles
}

// VCSService contains changes controll
type VCSService interface {
	ReadDataFromFS(fs filesystem.Filespace) (data VCSData, err error)
	WriteDataToFS(fs filesystem.Filespace, data VCSData) (err error)
}
