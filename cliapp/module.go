package cmdapp

import (
	"github.com/goatcms/goatcli/cliapp/commands"
	"github.com/goatcms/goatcli/cliapp/commands/clonec"
	"github.com/goatcms/goatcli/cliapp/services/repositiories"
	"github.com/goatcms/goatcore/app"
)

const (
	// CurrentPath is path to current directory
	CurrentPath = "./"
	// CurrentFilespace represent current directory
	CurrentFilespace = "current"
)

//Module is module contains all services
type Module struct {
}

//NewModule create new module instance
func NewModule() app.Module {
	return &Module{}
}

//RegisterDependencies is init callback to register module dependencies
func (m *Module) RegisterDependencies(a app.App) error {
	// filespaces
	if err := m.registerFilesystems(a); err != nil {
		return err
	}
	// commands
	if err := m.registerCommands(a); err != nil {
		return err
	}
	// services
	dp := a.DependencyProvider()
	if err := repositories.RegisterDependencies(dp); err != nil {
		return err
	}
	return nil
}

func (m *Module) registerCommands(a app.App) error {
	commandScope := a.CommandScope()
	// serve
	commandScope.Set("help.command.clone", commands.Clone)
	commandScope.Set("command.clone", clonec.Run)
	return nil
}

func (m *Module) registerFilesystems(a app.App) error {
	root := a.RootFilespace()
	// templates
	currentFS, err := root.Filespace(CurrentPath)
	if err != nil {
		return err
	}
	a.FilespaceScope().Set(CurrentFilespace, currentFS)
	return nil
}

// InitDependencies is init callback to inject dependencies inside module
func (m *Module) InitDependencies(a app.App) error {
	return nil
}

// Run is run event callback
func (m *Module) Run() error {
	return nil
}
