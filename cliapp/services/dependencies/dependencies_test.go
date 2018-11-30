package dependencies

import (
	"bytes"
	"strings"
	"testing"

	"github.com/goatcms/goatcli/cliapp/common/config"
	"github.com/goatcms/goatcli/cliapp/services"
	"github.com/goatcms/goatcli/cliapp/services/repositories"
	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/app/gio"
	"github.com/goatcms/goatcore/app/mockupapp"
)

func TestDependenciesFromFile(t *testing.T) {
	t.Parallel()
	var (
		err    error
		mapp   app.App
		output = new(bytes.Buffer)
	)
	// prepare mockup application & data
	if mapp, err = mockupapp.NewApp(mockupapp.MockupOptions{
		Input:  gio.NewInput(strings.NewReader("my_insert_value\n")),
		Output: gio.NewOutput(output),
	}); err != nil {
		t.Error(err)
		return
	}
	if err = mapp.RootFilespace().WriteFile(DependenciesDefPath, []byte(`[{
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
		Dependencies services.DependenciesService `dependency:"DependenciesService"`
	}
	if err = mapp.DependencyProvider().InjectTo(&deps); err != nil {
		t.Error(err)
		return
	}
	var dependencies []*config.Dependency
	if dependencies, err = deps.Dependencies.ReadDefFromFS(mapp.RootFilespace()); err != nil {
		t.Error(err)
		return
	}
	if len(dependencies) != 1 {
		t.Errorf("expected one dependency and take %d", len(dependencies))
		return
	}
}

func TestDependenciesDefaultEmpty(t *testing.T) {
	var err error
	t.Parallel()
	// prepare mockup application & data
	output := new(bytes.Buffer)
	mapp, err := mockupapp.NewApp(mockupapp.MockupOptions{
		Input:  gio.NewInput(strings.NewReader("")),
		Output: gio.NewOutput(output),
	})
	if err != nil {
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
		Dependencies services.DependenciesService `dependency:"DependenciesService"`
	}
	if err = mapp.DependencyProvider().InjectTo(&deps); err != nil {
		t.Error(err)
		return
	}
	var dependencies []*config.Dependency
	if dependencies, err = deps.Dependencies.ReadDefFromFS(mapp.RootFilespace()); err != nil {
		t.Error(err)
		return
	}
	if len(dependencies) != 0 {
		t.Errorf("expected no dependencies and take %d", len(dependencies))
		return
	}
}
