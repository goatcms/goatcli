package gcliservices

import (
	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/filesystem"
)

// Project contains single project data
type Project struct {
	Local        ProjectLocal
	PasswordHash []byte
	FS           filesystem.Filespace
	EncryptFS    filesystem.Filespace
	Properties   map[string]string
	Secrets      map[string]string
	Data         ApplicationData
}

// ProjectLocal contains local temporary data
type ProjectLocal struct {
	ID string `json:"id"`
}

// GCLIProjectManager return goat cli application inputs
type GCLIProjectManager interface {
	Project(ctx app.IOContext) (project *Project, err error)
}

// GCLIEnvironment manage environments
type GCLIEnvironment interface {
	LoadEnvs(ctxScope app.Scope, propertiesData, secretsData map[string]string) (err error)
}

// GCLIKeystorageService storage project keys (like saved hashed passowrds)
type GCLIKeystorageService interface {
	Password(key string, prompt string) (passwordHash []byte, err error)
	RePassword(key string, prompt, rePrompt string) (passwordHash []byte, err error)
	Set(key string, value []byte) (err error)
	Sets(kv map[string][]byte) (err error)
	Get(key string) (value []byte, err error)
	Remove(key string) (err error)
}
