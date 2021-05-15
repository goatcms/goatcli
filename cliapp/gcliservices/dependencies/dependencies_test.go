package dependencies

import (
	"strings"
	"testing"

	"github.com/goatcms/goatcli/cliapp/common/config"
	"github.com/goatcms/goatcli/cliapp/gcliservices"
	"github.com/goatcms/goatcli/cliapp/gcliservices/repositories"
	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/app/gio"
	"github.com/goatcms/goatcore/app/goatapp"
)

func TestDependenciesFromFile(t *testing.T) {
	t.Parallel()
	var (
		err  error
		mapp app.App
	)
	// prepare mockup application & data
	if mapp, err = goatapp.NewMockupApp(goatapp.Params{
		IO: goatapp.IO{
			In: gio.NewAppInput(strings.NewReader("my_insert_value\n")),
		},
	}); err != nil {
		t.Error(err)
		return
	}
	if err = mapp.Filespaces().Root().WriteFile(DependenciesDefPath, []byte(`[{
		"repo": "RepoValue1",
		"branch": "BranchValue1",
		"rev": "RevValue1",
		"dest": "DestValue1",
	}]`), 0766); err != nil {
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
	var dependencies []*config.Dependency
	if dependencies, err = deps.Dependencies.ReadDefFromFS(mapp.Filespaces().Root()); err != nil {
		t.Error(err)
		return
	}
	if len(dependencies) != 1 {
		t.Errorf("expected one dependency and take %d", len(dependencies))
		return
	}
}

func TestDependenciesDefaultEmpty(t *testing.T) {
	var (
		err  error
		mapp *goatapp.MockupApp
	)
	t.Parallel()
	// prepare mockup application & data
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
	var dependencies []*config.Dependency
	if dependencies, err = deps.Dependencies.ReadDefFromFS(mapp.Filespaces().Root()); err != nil {
		t.Error(err)
		return
	}
	if len(dependencies) != 0 {
		t.Errorf("expected no dependencies and take %d", len(dependencies))
		return
	}
}
