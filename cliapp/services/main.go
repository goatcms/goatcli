package services

import (
	"github.com/goatcms/goatcli/cliapp/common"
	"github.com/goatcms/goatcli/cliapp/common/config"
	"github.com/goatcms/goatcore/filesystem"
)

// Repositories provide git repository access
type Repositories interface {
	Filespace(repository, rev string) (filesystem.Filespace, error)
}

// Cloner clone an repository
type Cloner interface {
	Clone(repository, rev string, destfs filesystem.Filespace, si common.StringInjector) (err error)
}

// Properties provide project properties data
type Properties interface {
	ReadDefFromFS(fs filesystem.Filespace) ([]*config.Property, error)
	ReadDataFromFS(fs filesystem.Filespace) (map[string]string, error)
	FillData(def []*config.Property, data map[string]string, defaultData map[string]string) (bool, error)
	WriteDataToFS(fs filesystem.Filespace, data map[string]string) error
}

// Modules proccess and return modules
type Modules interface {
	ReadDefFromFS(fs filesystem.Filespace) ([]*config.Module, error)
}
