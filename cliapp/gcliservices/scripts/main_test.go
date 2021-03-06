package scripts

import (
	"github.com/goatcms/goatcli/cliapp/gcliservices/template"
	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/app/bootstrap"
	"github.com/goatcms/goatcore/app/goatapp"
	"github.com/goatcms/goatcore/app/modules/commonm"
	"github.com/goatcms/goatcore/app/modules/ocm"
	"github.com/goatcms/goatcore/app/modules/pipelinem"
	"github.com/goatcms/goatcore/app/modules/terminalm"
	"github.com/goatcms/goatcore/varutil/goaterr"
)

func newApp(params goatapp.Params) (mapp *goatapp.MockupApp, bootstraper app.Bootstrap, err error) {
	if mapp, err = goatapp.NewMockupApp(params); err != nil {
		return nil, nil, err
	}
	dp := mapp.DependencyProvider()
	bootstraper = bootstrap.NewBootstrap(mapp)
	if err = goaterr.ToError(goaterr.AppendError(nil,
		bootstraper.Register(terminalm.NewModule()),
		bootstraper.Register(commonm.NewModule()),
		bootstraper.Register(pipelinem.NewModule()),
		bootstraper.Register(ocm.NewModule()),
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
