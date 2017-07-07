package cliapp

import (
	"github.com/goatcms/goatcli/cliapp/commands"
	"github.com/goatcms/goatcli/cliapp/commands/buildc"
	"github.com/goatcms/goatcli/cliapp/commands/clonec"
	"github.com/goatcms/goatcli/cliapp/commands/datac"
	"github.com/goatcms/goatcli/cliapp/commands/initc"
	"github.com/goatcms/goatcli/cliapp/services/builder"
	"github.com/goatcms/goatcli/cliapp/services/cloner"
	"github.com/goatcms/goatcli/cliapp/services/data"
	"github.com/goatcms/goatcli/cliapp/services/modules"
	"github.com/goatcms/goatcli/cliapp/services/properties"
	"github.com/goatcms/goatcli/cliapp/services/repositiories"
	"github.com/goatcms/goatcli/cliapp/services/template"
	"github.com/goatcms/goatcore/app"
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
	// commands
	if err := m.registerCommands(a); err != nil {
		return err
	}
	// services
	dp := a.DependencyProvider()
	if err := builder.RegisterDependencies(dp); err != nil {
		return err
	}
	if err := cloner.RegisterDependencies(dp); err != nil {
		return err
	}
	if err := data.RegisterDependencies(dp); err != nil {
		return err
	}
	if err := modules.RegisterDependencies(dp); err != nil {
		return err
	}
	if err := properties.RegisterDependencies(dp); err != nil {
		return err
	}
	if err := repositories.RegisterDependencies(dp); err != nil {
		return err
	}
	if err := template.RegisterDependencies(dp); err != nil {
		return err
	}
	return nil
}

func (m *Module) registerCommands(a app.App) error {
	commandScope := a.CommandScope()
	commandScope.Set("help.command.clone", commands.Clone)
	commandScope.Set("command.clone", clonec.Run)
	commandScope.Set("help.command.data:add", commands.DataAdd)
	commandScope.Set("command.data:add", datac.RunAdd)
	commandScope.Set("help.command.build", commands.Build)
	commandScope.Set("command.build", buildc.Run)
	commandScope.Set("help.command.init", commands.Init)
	commandScope.Set("command.init", initc.Run)
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
