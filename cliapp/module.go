package cliapp

import (
	"github.com/goatcms/goatcli/cliapp/gclicommands/buildc"
	"github.com/goatcms/goatcli/cliapp/gclicommands/cleanc"
	"github.com/goatcms/goatcli/cliapp/gclicommands/clonec"
	"github.com/goatcms/goatcli/cliapp/gclicommands/containerc"
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
	"github.com/goatcms/goatcli/cliapp/gcliservices/template/shtf"
	"github.com/goatcms/goatcli/cliapp/gcliservices/template/simpletf"
	"github.com/goatcms/goatcli/cliapp/gcliservices/vcs"
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
	m.registerCommands(a)
	m.registerHealthCheckers(a)
	return m.registerServices(a)
}

func (m *Module) registerCommands(a app.App) {
	term := a.Terminal()
	term.SetCommand(buildc.Commands()...)
	term.SetCommand(cleanc.Commands()...)
	term.SetCommand(clonec.Commands()...)
	term.SetCommand(containerc.Commands()...)
	term.SetCommand(datac.Commands()...)
	term.SetCommand(depsc.Commands()...)
	term.SetCommand(initc.Commands()...)
	term.SetCommand(propertiesc.Commands()...)
	term.SetCommand(scriptsc.Commands()...)
	term.SetCommand(secretsc.Commands()...)
	term.SetCommand(vcsc.Commands()...)
}

func (m *Module) registerHealthCheckers(a app.App) {
	a.SetHealthChecker("git", GitHealthChecker)
}

func (m *Module) registerServices(a app.App) error {
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
	))
}

// InitDependencies is init callback to inject dependencies inside module
func (m *Module) InitDependencies(a app.App) (err error) {
	dp := a.DependencyProvider()
	return goaterr.ToError(goaterr.AppendError(nil,
		amtf.RegisterFunctions(dp),
		simpletf.RegisterFunctions(dp),
		shtf.RegisterFunctions(dp),
	))
}

// Run is run event callback
func (m *Module) Run(a app.App) error {
	return nil
}
