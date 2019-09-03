package builder

import (
	"bytes"
	"strings"
	"testing"

	"github.com/goatcms/goatcli/cliapp/common/config"
	"github.com/goatcms/goatcli/cliapp/services"
	"github.com/goatcms/goatcli/cliapp/services/dependencies"
	"github.com/goatcms/goatcli/cliapp/services/modules"
	"github.com/goatcms/goatcli/cliapp/services/repositories"
	"github.com/goatcms/goatcli/cliapp/services/template"
	"github.com/goatcms/goatcli/cliapp/services/vcs"
	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/app/gio"
	"github.com/goatcms/goatcore/app/mockupapp"
	"github.com/goatcms/goatcore/varutil/goaterr"
)

const (
	testBuildDefJSON = `[{"from":"fromv","to":"tov","layout":"layoutv","template":"templatev", "properties":{"key1":"value1"}}]`
)

func TestDataDefFromFile(t *testing.T) {
	var (
		mapp       app.App
		err        error
		builderDef []*config.Build
	)
	t.Parallel()
	// prepare mockup application
	if mapp, err = mockupapp.NewApp(mockupapp.MockupOptions{
		Input:  gio.NewInput(strings.NewReader("")),
		Output: gio.NewOutput(new(bytes.Buffer)),
	}); err != nil {
		t.Error(err)
		return
	}
	if err = mapp.RootFilespace().WriteFile(BuildDefPath, []byte(testBuildDefJSON), 0766); err != nil {
		t.Error(err)
		return
	}
	if err = goaterr.ToErrors(goaterr.AppendError(nil,
		RegisterDependencies(mapp.DependencyProvider()),
		modules.RegisterDependencies(mapp.DependencyProvider()),
		dependencies.RegisterDependencies(mapp.DependencyProvider()),
		repositories.RegisterDependencies(mapp.DependencyProvider()),
		template.RegisterDependencies(mapp.DependencyProvider()),
		vcs.RegisterDependencies(mapp.DependencyProvider()),
		template.InitDependencies(mapp))); err != nil {
		t.Error(err)
		return
	}
	// test
	var deps struct {
		BuilderService services.BuilderService `dependency:"BuilderService"`
	}
	if err = mapp.DependencyProvider().InjectTo(&deps); err != nil {
		t.Error(err)
		return
	}
	if builderDef, err = deps.BuilderService.ReadDefFromFS(mapp.RootFilespace()); err != nil {
		t.Error(err)
		return
	}
	if len(builderDef) != 1 {
		t.Errorf("expected one module and take %d", len(builderDef))
		return
	}
}

func TestDataDefDefaultEmpty(t *testing.T) {
	var (
		err      error
		mapp     app.App
		buildDef []*config.Build
	)
	t.Parallel()
	// prepare mockup application
	if mapp, err = mockupapp.NewApp(mockupapp.MockupOptions{
		Input:  gio.NewInput(strings.NewReader("")),
		Output: gio.NewOutput(new(bytes.Buffer)),
	}); err != nil {
		t.Error(err)
		return
	}
	if err = goaterr.ToErrors(goaterr.AppendError(nil,
		RegisterDependencies(mapp.DependencyProvider()),
		modules.RegisterDependencies(mapp.DependencyProvider()),
		dependencies.RegisterDependencies(mapp.DependencyProvider()),
		repositories.RegisterDependencies(mapp.DependencyProvider()),
		template.RegisterDependencies(mapp.DependencyProvider()),
		vcs.RegisterDependencies(mapp.DependencyProvider()),
		template.InitDependencies(mapp))); err != nil {
		t.Error(err)
		return
	}
	// test
	var deps struct {
		BuilderService services.BuilderService `dependency:"BuilderService"`
	}
	if err = mapp.DependencyProvider().InjectTo(&deps); err != nil {
		t.Error(err)
		return
	}
	if buildDef, err = deps.BuilderService.ReadDefFromFS(mapp.RootFilespace()); err != nil {
		t.Error(err)
		return
	}
	if len(buildDef) != 0 {
		t.Errorf("should return empty array when def file is not exist (result %v)", buildDef)
		return
	}
}
