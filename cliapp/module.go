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
	"github.com/goatcms/goatcore/varutil/goaterr"
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
	return goaterr.ToErrors(goaterr.AppendError(nil,
		app.RegisterComand(a, "clone", clonec.Run, commands.Clone),
		app.RegisterComand(a, "init", initc.RunInit, commands.Init),
		app.RegisterComand(a, "build", buildc.Run, commands.Build),
		app.RegisterComand(a, "data:add", datac.RunAdd, commands.DataAdd),
		app.RegisterComand(a, "deps:add", depsc.RunAddDep, commands.AddDep),
		app.RegisterComand(a, "deps:add:go", depsc.RunAddGODep, commands.AddGODep),
		app.RegisterComand(a, "deps:add:go:import", depsc.RunAddGOImportsDep, commands.AddGOImportsDep),
		app.RegisterComand(a, "properties:set", propertiesc.RunSetPropertyValue, commands.SetPropertyValueDep),
		app.RegisterComand(a, "properties:get", propertiesc.RunGetPropertyValue, commands.GetPropertyValueDep),
		app.RegisterComand(a, "secrets:set", secretsc.RunSetSecretValue, commands.SetSecretValueDep),
		app.RegisterComand(a, "secrets:get", secretsc.RunGetSecretValue, commands.GetSecretValueDep),
		app.RegisterArgument(a, "cwd", commands.CWDArg),
	))
}

// InitDependencies is init callback to inject dependencies inside module
func (m *Module) InitDependencies(a app.App) error {
	return template.InitDependencies(a)
}

// Run is run event callback
func (m *Module) Run(a app.App) error {
	return nil
}
