package services

import (
	"html/template"

	"github.com/goatcms/goatcli/cliapp/common"
	"github.com/goatcms/goatcli/cliapp/common/config"
	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/filesystem"
)

// Repositories provide git repository access
type Repositories interface {
	Filespace(repository, rev string) (filesystem.Filespace, error)
}

// Cloner clone an repository
type Cloner interface {
	Clone(repository, rev string, destfs filesystem.Filespace, replaces []*config.Replace) (err error)
}

// Project provide project api
type Project interface {
	Filespace() (filesystem.Filespace, error)
}

// Properties provide project properties data
type Properties interface {
	Get(fs filesystem.Filespace) (common.PropertiesResult, error)
}

type Modules interface {
	Init() error
	ModulesConfig() ([]*config.Module, error)
}

type Template interface {
	AddFunc(name string, f interface{}) error
	View(layoutName, viewName string, eventScope app.EventScope) (*template.Template, error)
}

type Compiler interface{}
type FileCompiler interface{}
type TemplateCompiler interface{}
