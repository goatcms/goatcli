package dependencies

import (
	"testing"

	"github.com/goatcms/goatcli/cliapp/common/config"
	"github.com/goatcms/goatcli/cliapp/gcliservices"
	"github.com/goatcms/goatcli/cliapp/gcliservices/repositories"
	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/app/goatapp"
)

func TestWriteDefToFS(t *testing.T) {
	var (
		err    error
		mapp   app.App
		list   []*config.Dependency
		result []*config.Dependency
	)
	t.Parallel()
	// prepare mockup application
	if mapp, err = goatapp.NewMockupApp(goatapp.Params{}); err != nil {
		t.Error(err)
		return
	}
	if err = RegisterDependencies(mapp.DependencyProvider()); err != nil {
		t.Error(err)
		return
	}
	if err = repositories.RegisterDependencies(mapp.DependencyProvider()); err != nil {
		t.Error(err)
		return
	}
	// test
	var deps struct {
		Dependencies gcliservices.DependenciesService `dependency:"DependenciesService"`
	}
	if err = mapp.DependencyProvider().InjectTo(&deps); err != nil {
		t.Error(err)
		return
	}
	list = []*config.Dependency{{
		Repo:   "https://github.com/goatcms/goatcore.git",
		Branch: "master",
		Rev:    "3eb26366749bd54f3a871ff9beb5c565195f6233",
		Dest:   "vendor/github.com/goatcms/goatcore",
	}}
	if err = deps.Dependencies.WriteDefToFS(mapp.Filespaces().Root(), list); err != nil {
		t.Error(err)
		return
	}
	if result, err = deps.Dependencies.ReadDefFromFS(mapp.Filespaces().Root()); err != nil {
		t.Error(err)
		return
	}
	if len(result) != 1 {
		t.Errorf("expected one dependency exactly and take %v dependencies", len(result))
		return
	}
	if result[0].Repo != "https://github.com/goatcms/goatcore.git" {
		t.Errorf("expected hash equals to 'https://github.com/goatcms/goatcore.git' and take '%s'", result[0].Repo)
		return
	}
	if result[0].Branch != "master" {
		t.Errorf("expected hash equals to 'master' and take '%s'", result[0].Branch)
		return
	}
	if result[0].Rev != "3eb26366749bd54f3a871ff9beb5c565195f6233" {
		t.Errorf("expected revision equals to '3eb26366749bd54f3a871ff9beb5c565195f6233' and take '%s'", result[0].Rev)
		return
	}
	if result[0].Dest != "vendor/github.com/goatcms/goatcore" {
		t.Errorf("expected destination equals to 'vendor/github.com/goatcms/goatcore' and take '%s'", result[0].Dest)
		return
	}
}
