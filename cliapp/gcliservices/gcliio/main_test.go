package gcliio

import (
	"github.com/goatcms/goatcli/cliapp/gcliservices/template"
	"github.com/goatcms/goatcore/varutil/goaterr"

	"github.com/goatcms/goatcli/cliapp/gclimock"
	"github.com/goatcms/goatcli/cliapp/gcliservices/data"
	"github.com/goatcms/goatcli/cliapp/gcliservices/properties"
	"github.com/goatcms/goatcli/cliapp/gcliservices/secrets"
	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/app/mockupapp"
)

func newMockupApp(opt mockupapp.MockupOptions) (mapp app.App, err error) {
	if mapp, err = gclimock.NewApp(opt); err != nil {
		return nil, err
	}
	dp := mapp.DependencyProvider()
	if err = goaterr.ToError(goaterr.AppendError(nil,
		RegisterDependencies(dp),
		secrets.RegisterDependencies(dp),
		properties.RegisterDependencies(dp),
		template.RegisterDependencies(dp),
		data.RegisterDependencies(dp))); err != nil {
		return nil, err
	}
	return mapp, nil
}
