package modules

import (
	"strings"
	"testing"

	"github.com/goatcms/goatcli/cliapp/common/config"
	"github.com/goatcms/goatcli/cliapp/gcliservices"
	"github.com/goatcms/goatcore/app/gio"
	"github.com/goatcms/goatcore/app/goatapp"
)

const (
	testModulesDefJSON = `[{"srcClone":"srcCloneValue", "srcRev":"srcRevValue", "srcDir":"srcDirValue", "testClone":"testCloneValue", "testRev":"testRevValue", "testDir":"testDirValue"}]`
)

func TestModulesFromFile(t *testing.T) {
	var (
		err  error
		mapp *goatapp.MockupApp
	)
	t.Parallel()
	// prepare mockup application & data
	if mapp, err = goatapp.NewMockupApp(goatapp.Params{
		IO: goatapp.IO{
			In: gio.NewAppInput(strings.NewReader("my_insert_value\n")),
		},
	}); err != nil {
		t.Error(err)
		return
	}
	if err = mapp.Filespaces().Root().WriteFile(ModulesDefPath, []byte(testModulesDefJSON), 0766); err != nil {
		t.Error(err)
		return
	}
	if err = RegisterDependencies(mapp.DependencyProvider()); err != nil {
		t.Error(err)
		return
	}
	// test
	var deps struct {
		Modules gcliservices.ModulesService `dependency:"ModulesService"`
	}
	if err = mapp.DependencyProvider().InjectTo(&deps); err != nil {
		t.Error(err)
		return
	}
	var modules []*config.Module
	if modules, err = deps.Modules.ReadDefFromFS(mapp.Filespaces().Root()); err != nil {
		t.Error(err)
		return
	}
	if len(modules) != 1 {
		t.Errorf("expected one module and take %d", len(modules))
		return
	}
}

func TestModulesDefaultEmpty(t *testing.T) {
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
	// test
	var deps struct {
		Modules gcliservices.ModulesService `dependency:"ModulesService"`
	}
	if err = mapp.DependencyProvider().InjectTo(&deps); err != nil {
		t.Error(err)
		return
	}
	var modules []*config.Module
	if modules, err = deps.Modules.ReadDefFromFS(mapp.Filespaces().Root()); err != nil {
		t.Error(err)
		return
	}
	if len(modules) != 0 {
		t.Errorf("expected no modules and take %d", len(modules))
		return
	}
}
