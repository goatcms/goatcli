package scriptsc

import (
	"github.com/goatcms/goatcli/cliapp/gclicommands"
	"github.com/goatcms/goatcli/cliapp/gcliservices/data"
	"github.com/goatcms/goatcli/cliapp/gcliservices/gcliio"
	"github.com/goatcms/goatcli/cliapp/gcliservices/properties"
	"github.com/goatcms/goatcli/cliapp/gcliservices/scripts"
	"github.com/goatcms/goatcli/cliapp/gcliservices/secrets"
	"github.com/goatcms/goatcli/cliapp/gcliservices/template"
	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/app/bootstrap"
	"github.com/goatcms/goatcore/app/mockupapp"
	"github.com/goatcms/goatcore/app/modules/commonm"
	"github.com/goatcms/goatcore/app/modules/pipelinem"
	"github.com/goatcms/goatcore/app/modules/terminalm"
	"github.com/goatcms/goatcore/varutil/goaterr"
)

func newApp(options mockupapp.MockupOptions) (mapp *mockupapp.App, bootstraper app.Bootstrap, err error) {
	if mapp, err = mockupapp.NewApp(options); err != nil {
		return nil, nil, err
	}
	bootstraper = bootstrap.NewBootstrap(mapp)
	dp := mapp.DependencyProvider()
	if err = goaterr.ToError(goaterr.AppendError(nil,
		bootstraper.Register(terminalm.NewModule()),
		bootstraper.Register(commonm.NewModule()),
		bootstraper.Register(pipelinem.NewModule()),
		gcliio.RegisterDependencies(dp),
		data.RegisterDependencies(dp),
		properties.RegisterDependencies(dp),
		secrets.RegisterDependencies(dp),
		scripts.RegisterDependencies(dp),
		template.RegisterDependencies(dp),
		app.RegisterCommand(mapp, "scripts:run", RunScript, gclicommands.ScriptsRun),
	)); err != nil {
		return nil, nil, err
	}
	if err = bootstraper.Init(); err != nil {
		return nil, nil, err
	}
	return mapp, bootstraper, nil
}
