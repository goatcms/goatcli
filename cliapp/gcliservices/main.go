package gcliservices

import (
	"io"
	"text/template"

	"github.com/goatcms/goatcli/cliapp/common"
	"github.com/goatcms/goatcli/cliapp/common/config"
	"github.com/goatcms/goatcore/app"
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
	FillData(ctx app.IOContext, def []*config.Property, data map[string]string, defaultData map[string]string, interactive bool) (bool, error)
	WriteDataToFS(fs filesystem.Filespace, data map[string]string) error
}

// SecretsService provide project secret properties data (like passwords, credentials etc)
type SecretsService interface {
	ReadDefFromFS(fs filesystem.Filespace, properties common.ElasticData, data ApplicationData) ([]*config.Property, error)
	ReadDataFromFS(fs filesystem.Filespace) (map[string]string, error)
	FillData(ctx app.IOContext, def []*config.Property, data map[string]string, defaultData map[string]string, interactive bool) (bool, error)
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
	ConsoleReadData(ctx app.IOContext, def *config.DataSet) (map[string]string, error)
}

// TemplateService provide template api
type TemplateService interface {
	TemplatesExecutor() (TemplatesExecutor, error)
	TemplateExecutor(path string) (TemplateExecutor, error)
}

// TemplatesExecutor render data for view tree
type TemplatesExecutor interface {
	IsEmpty(layoutName, templatePath string) (is bool, err error)
	Templates(layoutName, templatePath string) (list []string, err error)
	Execute(layoutName, templatePath string, wr io.Writer, data interface{}) (err error)
	ExecuteTemplate(layoutName, templatePath, templateName string, wr io.Writer, data interface{}) (err error)
}

// TemplateExecutor render data for single template
type TemplateExecutor interface {
	IsEmpty() (is bool, err error)
	Templates() (list []string, err error)
	Execute(wr io.Writer, data interface{}) (err error)
	ExecuteTemplate(templateName string, wr io.Writer, data interface{}) (err error)
}

// TemplateAssetsProvider return assets templates
type TemplateAssetsProvider interface {
	Layout(name string) (tmpl *template.Template, err error)
	Base() (*template.Template, error)
}

// TemplateConfig return templates configuration
type TemplateConfig interface {
	AddFunc(name string, f interface{}) (err error)
	Func() template.FuncMap
	IsCached() bool
	FS() filesystem.Filespace
}
