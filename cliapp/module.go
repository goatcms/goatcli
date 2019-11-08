package cliapp

import (
	"github.com/goatcms/goatcli/cliapp/commands"
	"github.com/goatcms/goatcli/cliapp/commands/buildc"
	"github.com/goatcms/goatcli/cliapp/commands/cleanc"
	"github.com/goatcms/goatcli/cliapp/commands/clonec"
	"github.com/goatcms/goatcli/cliapp/commands/datac"
	"github.com/goatcms/goatcli/cliapp/commands/depsc"
	"github.com/goatcms/goatcli/cliapp/commands/initc"
	"github.com/goatcms/goatcli/cliapp/commands/propertiesc"
	"github.com/goatcms/goatcli/cliapp/commands/secretsc"
	"github.com/goatcms/goatcli/cliapp/commands/vcsc"
	"github.com/goatcms/goatcli/cliapp/services/builder"
	"github.com/goatcms/goatcli/cliapp/services/cloner"
	"github.com/goatcms/goatcli/cliapp/services/data"
	"github.com/goatcms/goatcli/cliapp/services/dependencies"
	"github.com/goatcms/goatcli/cliapp/services/modules"
	"github.com/goatcms/goatcli/cliapp/services/properties"
	"github.com/goatcms/goatcli/cliapp/services/repositories"
	"github.com/goatcms/goatcli/cliapp/services/secrets"
	"github.com/goatcms/goatcli/cliapp/services/template"
	"github.com/goatcms/goatcli/cliapp/services/template/amtf"
	"github.com/goatcms/goatcli/cliapp/services/template/simpletf"
	"github.com/goatcms/goatcli/cliapp/services/vcs"
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
	return goaterr.ToErrors(goaterr.AppendError(nil,
		builder.RegisterDependencies(dp),
		cloner.RegisterDependencies(dp),
		data.RegisterDependencies(dp),
		modules.RegisterDependencies(dp),
		dependencies.RegisterDependencies(dp),
		properties.RegisterDependencies(dp),
		secrets.RegisterDependencies(dp),
		repositories.RegisterDependencies(dp),
		template.RegisterDependencies(dp),
		vcs.RegisterDependencies(dp)))
}

func (m *Module) registerCommands(a app.App) error {
	return goaterr.ToErrors(goaterr.AppendError(nil,
		app.RegisterCommand(a, "clone", clonec.Run, commands.Clone),
		app.RegisterCommand(a, "init", initc.RunInit, commands.Init),
		app.RegisterCommand(a, "build", buildc.RunBuild, commands.Build),
		app.RegisterCommand(a, "rebuild", buildc.RunRebuild, commands.Rebuild),
		app.RegisterCommand(a, "clean", cleanc.RunClean, commands.Clean),
		app.RegisterCommand(a, "clean:dependencies", cleanc.RunCleanDependencies, commands.CleanDependencies),
		app.RegisterCommand(a, "clean:build", cleanc.RunCleanBuild, commands.CleanBuild),
		app.RegisterCommand(a, "data:add", datac.RunAdd, commands.DataAdd),
		app.RegisterCommand(a, "deps:add", depsc.RunAddDep, commands.AddDep),
		app.RegisterCommand(a, "deps:add:go", depsc.RunAddGODep, commands.AddGODep),
		app.RegisterCommand(a, "deps:add:go:import", depsc.RunAddGOImportsDep, commands.AddGOImportsDep),
		app.RegisterCommand(a, "properties:set", propertiesc.RunSetPropertyValue, commands.SetPropertyValueDep),
		app.RegisterCommand(a, "properties:get", propertiesc.RunGetPropertyValue, commands.GetPropertyValueDep),
		app.RegisterCommand(a, "secrets:set", secretsc.RunSetSecretValue, commands.SetSecretValueDep),
		app.RegisterCommand(a, "secrets:get", secretsc.RunGetSecretValue, commands.GetSecretValueDep),
		app.RegisterCommand(a, "vcs:clean", vcsc.RunClean, commands.VCSClean),
		app.RegisterCommand(a, "vcs:scan", vcsc.RunScan, commands.VCSScan),
		app.RegisterCommand(a, "vcs:ignored:add", vcsc.RunIgnoredAdd, commands.VCSIgnoredAdd),
		app.RegisterCommand(a, "vcs:ignored:remove", vcsc.RunIgnoredRemove, commands.VCSIgnoredRemove),
		app.RegisterCommand(a, "vcs:ignored:list", vcsc.RunIgnoredList, commands.VCSIgnoredList),
		app.RegisterCommand(a, "vcs:generated:list", vcsc.RunGeneratedList, commands.VCSGeneratedList),
		app.RegisterArgument(a, "cwd", commands.CWDArg),
	))
}

// InitDependencies is init callback to inject dependencies inside module
func (m *Module) InitDependencies(a app.App) (err error) {
	dp := a.DependencyProvider()
	return goaterr.ToErrors(goaterr.AppendError(nil,
		amtf.RegisterFunctions(dp),
		simpletf.RegisterFunctions(dp),
	))
}

// Run is run event callback
func (m *Module) Run(a app.App) error {
	return nil
}
