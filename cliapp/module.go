package cliapp

import (
	"github.com/goatcms/goatcli/cliapp/common/uefs"
	"github.com/goatcms/goatcli/cliapp/gclicommands"
	"github.com/goatcms/goatcli/cliapp/gclicommands/buildc"
	"github.com/goatcms/goatcli/cliapp/gclicommands/cleanc"
	"github.com/goatcms/goatcli/cliapp/gclicommands/clonec"
	"github.com/goatcms/goatcli/cliapp/gclicommands/datac"
	"github.com/goatcms/goatcli/cliapp/gclicommands/depsc"
	"github.com/goatcms/goatcli/cliapp/gclicommands/initc"
	"github.com/goatcms/goatcli/cliapp/gclicommands/propertiesc"
	"github.com/goatcms/goatcli/cliapp/gclicommands/scriptsc"
	"github.com/goatcms/goatcli/cliapp/gclicommands/secretsc"
	"github.com/goatcms/goatcli/cliapp/gclicommands/vcsc"
	"github.com/goatcms/goatcli/cliapp/gcliservices/builder"
	"github.com/goatcms/goatcli/cliapp/gcliservices/cloner"
	"github.com/goatcms/goatcli/cliapp/gcliservices/data"
	"github.com/goatcms/goatcli/cliapp/gcliservices/dependencies"
	"github.com/goatcms/goatcli/cliapp/gcliservices/gcliio"
	"github.com/goatcms/goatcli/cliapp/gcliservices/modules"
	"github.com/goatcms/goatcli/cliapp/gcliservices/properties"
	"github.com/goatcms/goatcli/cliapp/gcliservices/repositories"
	"github.com/goatcms/goatcli/cliapp/gcliservices/scripts"
	"github.com/goatcms/goatcli/cliapp/gcliservices/secrets"
	"github.com/goatcms/goatcli/cliapp/gcliservices/template"
	"github.com/goatcms/goatcli/cliapp/gcliservices/template/amtf"
	"github.com/goatcms/goatcli/cliapp/gcliservices/template/simpletf"
	"github.com/goatcms/goatcli/cliapp/gcliservices/vcs"
	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/filesystem"
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
	// filespaces
	if err = m.registerFilespaces(a); err != nil {
		return err
	}
	// services
	dp := a.DependencyProvider()
	return goaterr.ToError(goaterr.AppendError(nil,
		builder.RegisterDependencies(dp),
		cloner.RegisterDependencies(dp),
		data.RegisterDependencies(dp),
		modules.RegisterDependencies(dp),
		dependencies.RegisterDependencies(dp),
		properties.RegisterDependencies(dp),
		secrets.RegisterDependencies(dp),
		repositories.RegisterDependencies(dp),
		template.RegisterDependencies(dp),
		vcs.RegisterDependencies(dp),
		scripts.RegisterDependencies(dp),
		gcliio.RegisterDependencies(dp),
		app.RegisterHealthChecker(a, "git", GitHealthChecker)))
}

func (m *Module) registerFilespaces(a app.App) (err error) {
	var (
		gefs filesystem.Filespace
		fs   filesystem.Filespace
	)
	if err = a.HomeFilespace().MkdirAll(".goatcli/efs"); err != nil {
		return err
	}
	if fs, err = a.HomeFilespace().Filespace(".goatcli/efs"); err != nil {
		return err
	}
	if gefs, err = uefs.BuildEFS(fs, true, true, []byte{}); err != nil {
		return err
	}
	return a.FilespaceScope().Set("gefs", gefs)
}

func (m *Module) registerCommands(a app.App) error {
	return goaterr.ToError(goaterr.AppendError(nil,
		app.RegisterCommand(a, "clone", clonec.Run, gclicommands.Clone),
		app.RegisterCommand(a, "init", initc.RunInit, gclicommands.Init),
		app.RegisterCommand(a, "build", buildc.RunBuild, gclicommands.Build),
		app.RegisterCommand(a, "rebuild", buildc.RunRebuild, gclicommands.Rebuild),
		app.RegisterCommand(a, "clean", cleanc.RunClean, gclicommands.Clean),
		app.RegisterCommand(a, "clean:dependencies", cleanc.RunCleanDependencies, gclicommands.CleanDependencies),
		app.RegisterCommand(a, "clean:build", cleanc.RunCleanBuild, gclicommands.CleanBuild),
		app.RegisterCommand(a, "data:add", datac.RunAdd, gclicommands.DataAdd),
		app.RegisterCommand(a, "deps:add", depsc.RunAddDep, gclicommands.AddDep),
		app.RegisterCommand(a, "deps:add:go", depsc.RunAddGODep, gclicommands.AddGODep),
		app.RegisterCommand(a, "deps:add:go:import", depsc.RunAddGOImportsDep, gclicommands.AddGOImportsDep),
		app.RegisterCommand(a, "properties:set", propertiesc.RunSetPropertyValue, gclicommands.SetPropertyValueDep),
		app.RegisterCommand(a, "properties:get", propertiesc.RunGetPropertyValue, gclicommands.GetPropertyValueDep),
		app.RegisterCommand(a, "secrets:set", secretsc.RunSetSecretValue, gclicommands.SetSecretValueDep),
		app.RegisterCommand(a, "secrets:get", secretsc.RunGetSecretValue, gclicommands.GetSecretValueDep),
		app.RegisterCommand(a, "vcs:clean", vcsc.RunClean, gclicommands.VCSClean),
		app.RegisterCommand(a, "vcs:scan", vcsc.RunScan, gclicommands.VCSScan),
		app.RegisterCommand(a, "vcs:persisted:add", vcsc.RunPersistedAdd, gclicommands.VCSPersistedAdd),
		app.RegisterCommand(a, "vcs:persisted:remove", vcsc.RunPersistedRemove, gclicommands.VCSPersistedRemove),
		app.RegisterCommand(a, "vcs:persisted:list", vcsc.RunPersistedList, gclicommands.VCSPersistedList),
		app.RegisterCommand(a, "vcs:generated:list", vcsc.RunGeneratedList, gclicommands.VCSGeneratedList),
		app.RegisterCommand(a, "scripts:run", scriptsc.RunScript, gclicommands.ScriptsRun),
		app.RegisterCommand(a, "scripts:envs", scriptsc.RunScriptsEnvs, gclicommands.ScriptsEnvs),
		app.RegisterArgument(a, "cwd", gclicommands.CWDArg),
	))
}

// InitDependencies is init callback to inject dependencies inside module
func (m *Module) InitDependencies(a app.App) (err error) {
	dp := a.DependencyProvider()
	return goaterr.ToError(goaterr.AppendError(nil,
		amtf.RegisterFunctions(dp),
		simpletf.RegisterFunctions(dp),
	))
}

// Run is run event callback
func (m *Module) Run(a app.App) error {
	return nil
}
