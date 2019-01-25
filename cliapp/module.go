package cliapp

import (
	"github.com/goatcms/goatcli/cliapp/commands"
	"github.com/goatcms/goatcli/cliapp/commands/buildc"
	"github.com/goatcms/goatcli/cliapp/commands/clonec"
	"github.com/goatcms/goatcli/cliapp/commands/datac"
	"github.com/goatcms/goatcli/cliapp/commands/depsc"
	"github.com/goatcms/goatcli/cliapp/commands/initc"
	"github.com/goatcms/goatcli/cliapp/commands/propertiesc"
	"github.com/goatcms/goatcli/cliapp/commands/secretsc"
	"github.com/goatcms/goatcli/cliapp/services/builder"
	"github.com/goatcms/goatcli/cliapp/services/cloner"
	"github.com/goatcms/goatcli/cliapp/services/data"
	"github.com/goatcms/goatcli/cliapp/services/dependencies"
	"github.com/goatcms/goatcli/cliapp/services/modules"
	"github.com/goatcms/goatcli/cliapp/services/properties"
	"github.com/goatcms/goatcli/cliapp/services/repositories"
	"github.com/goatcms/goatcli/cliapp/services/secrets"
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
func (m *Module) RegisterDependencies(a app.App) (err error) {
	// commands
	if err = m.registerCommands(a); err != nil {
		return err
	}
	// services
	dp := a.DependencyProvider()
	if err = builder.RegisterDependencies(dp); err != nil {
		return err
	}
	if err = cloner.RegisterDependencies(dp); err != nil {
		return err
	}
	if err = data.RegisterDependencies(dp); err != nil {
		return err
	}
	if err = modules.RegisterDependencies(dp); err != nil {
		return err
	}
	if err = dependencies.RegisterDependencies(dp); err != nil {
		return err
	}
	if err = properties.RegisterDependencies(dp); err != nil {
		return err
	}
	if err = secrets.RegisterDependencies(dp); err != nil {
		return err
	}
	if err = repositories.RegisterDependencies(dp); err != nil {
		return err
	}
	return template.RegisterDependencies(dp)
}

func (m *Module) registerCommands(a app.App) error {
	commandScope := a.CommandScope()
	// core commands
	commandScope.Set("help.command.clone", commands.Clone)
	commandScope.Set("command.clone", clonec.Run)
	commandScope.Set("help.command.init", commands.Init)
	commandScope.Set("command.init", initc.RunInit)
	commandScope.Set("help.command.build", commands.Build)
	commandScope.Set("command.build", buildc.Run)
	// data commands
	commandScope.Set("help.command.data:add", commands.DataAdd)
	commandScope.Set("command.data:add", datac.RunAdd)
	// dependencies commands
	commandScope.Set("help.command.deps:add", commands.AddDep)
	commandScope.Set("command.deps:add", depsc.RunAddDep)
	commandScope.Set("help.command.deps:add:go", commands.AddGODep)
	commandScope.Set("command.deps:add:go", depsc.RunAddGODep)
	commandScope.Set("help.command.deps:add:go:import", commands.AddGOImportsDep)
	commandScope.Set("command.deps:add:go:import", depsc.RunAddGOImportsDep)
	// properties command
	commandScope.Set("help.command.properties:set", commands.SetPropertyValueDep)
	commandScope.Set("command.properties:set", propertiesc.RunSetPropertyValue)
	commandScope.Set("help.command.properties:get", commands.GetPropertyValueDep)
	commandScope.Set("command.properties:get", propertiesc.RunGetPropertyValue)
	// secrets command
	commandScope.Set("help.command.secrets:set", commands.SetSecretValueDep)
	commandScope.Set("command.secrets:set", secretsc.RunSetSecretValue)
	commandScope.Set("help.command.secrets:get", commands.GetSecretValueDep)
	commandScope.Set("command.secrets:get", secretsc.RunGetSecretValue)
	// arguments
	commandScope.Set("help.argument.cwd", commands.CWDArg)
	return nil
}

// InitDependencies is init callback to inject dependencies inside module
func (m *Module) InitDependencies(a app.App) error {
	return template.InitDependencies(a)
}

// Run is run event callback
func (m *Module) Run() error {
	return nil
}
