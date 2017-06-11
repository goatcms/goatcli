package services

import (
	"github.com/goatcms/goatcli/cliapp/common/config"
	"github.com/goatcms/goatcore/filesystem"
)

const (
	// RepositoriesService provide git repository access
	RepositoriesService = "Repositories"
	// ProjectService provide git repository access
	ProjectService = "Project"
)

// Repositories provide git repository access
type Repositories interface {
	Filespace(repository, rev string) (filesystem.Filespace, error)
}

// Project provide project api
type Project interface {
	Filespace() (filesystem.Filespace, error)
}

// Properties provide project properties data
type Properties interface {
	Get(key string) (string, error)
}

type Modules interface {
	Init() error
	ModulesConfig() ([]*config.Module, error)
}
