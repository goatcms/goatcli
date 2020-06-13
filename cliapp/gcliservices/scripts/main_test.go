package scripts

import (
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
	if err = app.RegisterCommand(mapp, "testCommand", func(a app.App, ctx app.IOContext) (err error) {
		return ctx.IO().Out().Printf("output")
	}, "description"); err != nil {
		return nil, nil, err
	}
	dp := mapp.DependencyProvider()
	bootstraper = bootstrap.NewBootstrap(mapp)
	if err = goaterr.ToError(goaterr.AppendError(nil,
		bootstraper.Register(terminalm.NewModule()),
		bootstraper.Register(commonm.NewModule()),
		bootstraper.Register(pipelinem.NewModule()),
		RegisterDependencies(dp),
		template.RegisterDependencies(dp),
	)); err != nil {
		return nil, nil, err
	}
	if err = bootstraper.Init(); err != nil {
		return nil, nil, err
	}
	return mapp, bootstraper, nil
}
