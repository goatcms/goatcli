package services

import (
	"io"

	"github.com/goatcms/goatcli/cliapp/common"
	"github.com/goatcms/goatcli/cliapp/common/config"
	"github.com/goatcms/goatcore/filesystem"
)

const (
	// StremDataSeparator separates data in streams
	StremDataSeparator = ":"
)

// RepositoriesService provide git repository access
type RepositoriesService interface {
	Filespace(repository, rev string) (filesystem.Filespace, error)
}

// ClonerService clone an repository
type ClonerService interface {
	Clone(repository, rev string, destfs filesystem.Filespace, si common.StringInjector) (err error)
}

// PropertiesService provide project properties data
type PropertiesService interface {
	ReadDefFromFS(fs filesystem.Filespace) ([]*config.Property, error)
	ReadDataFromFS(fs filesystem.Filespace) (map[string]string, error)
	FillData(def []*config.Property, data map[string]string, defaultData map[string]string) (bool, error)
	WriteDataToFS(fs filesystem.Filespace, data map[string]string) error
}

// ModulesService proccess and return modules
type ModulesService interface {
	ReadDefFromFS(fs filesystem.Filespace) ([]*config.Module, error)
}

// DataService provide data api
type DataService interface {
	ReadDefFromFS(fs filesystem.Filespace) ([]*config.DataSet, error)
	ReadDataFromFS(fs filesystem.Filespace) (map[string]string, error)
	WriteDataToFS(fs filesystem.Filespace, prefix string, data map[string]string) error
	ConsoleReadData(def *config.DataSet) (map[string]string, error)
}

// TemplateService provide template api
type TemplateService interface {
	AddFunc(name string, f interface{}) error
	Build(fs filesystem.Filespace) (TemplateExecutor, error)
}

// TemplateExecutor render data
type TemplateExecutor interface {
	Execute(layoutName, TemplatePath string, wr io.Writer, data interface{}) error
}

// BuilderService build project structure
type BuilderService interface {
	ReadDefFromFS(fs filesystem.Filespace) ([]*config.Build, error)
	Build(fs filesystem.Filespace, buildConfigs []*config.Build, data, properties map[string]string) error
}
