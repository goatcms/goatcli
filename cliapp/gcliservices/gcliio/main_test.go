package gcliio

import (
	"github.com/goatcms/goatcore/filesystem"
	"github.com/goatcms/goatcore/varutil/goaterr"

	"github.com/goatcms/goatcli/cliapp/gcliservices/data"
	"github.com/goatcms/goatcli/cliapp/gcliservices/properties"
	"github.com/goatcms/goatcli/cliapp/gcliservices/secrets"
	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/app/mockupapp"
)

func newMockupApp(opt mockupapp.MockupOptions) (mapp app.App, err error) {
	var gefs filesystem.Filespace
	if mapp, err = mockupapp.NewApp(opt); err != nil {
		return nil, err
	}
	dp := mapp.DependencyProvider()
	if err = goaterr.ToError(goaterr.AppendError(nil,
		RegisterDependencies(dp),
		secrets.RegisterDependencies(dp),
		properties.RegisterDependencies(dp),
		data.RegisterDependencies(dp))); err != nil {
		return nil, err
	}
	if mapp.RootFilespace().MkdirAll(".goatcli/efs", filesystem.DefaultUnixDirMode); err != nil {
		return nil, err
	}
	if gefs, err = mapp.RootFilespace().Filespace(".goatcli/efs"); err != nil {
		return nil, err
	}
	if err = mapp.FilespaceScope().Set("gefs", gefs); err != nil {
		return nil, err
	}
	return mapp, nil
}
