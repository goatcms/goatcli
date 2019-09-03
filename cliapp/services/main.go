package services

import (
	"io"

	"github.com/goatcms/goatcli/cliapp/common"
	"github.com/goatcms/goatcli/cliapp/common/config"
	"github.com/goatcms/goatcore/filesystem"
	"github.com/goatcms/goatcore/repositories"
)

const (
	// StremDataSeparator separates data in streams
	StremDataSeparator = ":"
)

// RepositoriesService provide git repository access
type RepositoriesService interface {
	Filespace(repoURL string, version repositories.Version) (filesystem.Filespace, error)
}

// ClonerService clone an repository
type ClonerService interface {
	Clone(repoURL string, verion repositories.Version, destfs filesystem.Filespace, si common.StringInjector) (err error)
	CloneModules(sourcefs, destfs filesystem.Filespace, si common.StringInjector) (err error)
}

// PropertiesService provide project properties data
type PropertiesService interface {
	ReadDefFromFS(fs filesystem.Filespace) ([]*config.Property, error)
	ReadDataFromFS(fs filesystem.Filespace) (map[string]string, error)
	FillData(def []*config.Property, data map[string]string, defaultData map[string]string, interactive bool) (bool, error)
	WriteDataToFS(fs filesystem.Filespace, data map[string]string) error
}

// SecretsService provide project secret properties data (like passwords, credentials etc)
type SecretsService interface {
	ReadDefFromFS(fs filesystem.Filespace) ([]*config.Property, error)
	ReadDataFromFS(fs filesystem.Filespace) (map[string]string, error)
	FillData(def []*config.Property, data map[string]string, defaultData map[string]string, interactive bool) (bool, error)
	WriteDataToFS(fs filesystem.Filespace, data map[string]string) error
}

// ModulesService process and return modules
type ModulesService interface {
	ReadDefFromFS(fs filesystem.Filespace) ([]*config.Module, error)
}

// DependenciesService process and return modules
type DependenciesService interface {
	ReadDefFromFS(fs filesystem.Filespace) ([]*config.Dependency, error)
	WriteDefToFS(fs filesystem.Filespace, dependencies []*config.Dependency) error
	CloneDependencies(fs filesystem.Filespace, deps []*config.Dependency) error
}

// DataService provide data api
type DataService interface {
	HasDataFromFS(fs filesystem.Filespace, prefix string) bool
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
	Execute(layoutName, TemplatePath string, wr io.Writer, data interface{}) (err error)
	ExecuteTemplate(layoutName, viewName, templateName string, wr io.Writer, data interface{}) (err error)
}
