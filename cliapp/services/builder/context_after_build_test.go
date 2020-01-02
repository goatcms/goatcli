package builder

import (
	"strings"
	"testing"

	"github.com/goatcms/goatcli/cliapp/common/am"
	"github.com/goatcms/goatcli/cliapp/services"
	"github.com/goatcms/goatcli/cliapp/services/dependencies"
	"github.com/goatcms/goatcli/cliapp/services/modules"
	"github.com/goatcms/goatcli/cliapp/services/repositories"
	"github.com/goatcms/goatcli/cliapp/services/template"
	"github.com/goatcms/goatcli/cliapp/services/template/simpletf"
	"github.com/goatcms/goatcli/cliapp/services/vcs"
	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/app/gio"
	"github.com/goatcms/goatcore/app/mockupapp"
	"github.com/goatcms/goatcore/app/scope"
	"github.com/goatcms/goatcore/varutil/goaterr"
)

const (
	testCTXBuilderAfterBuildTemplate = `
	{{$ctx := .}}
	`
	testCTXBuilderAfterBuildConfig = `[{
	  "from":"ignore",
	  "to":"ignore",
	  "template":"names",
	  "layout":"default",
	  "afterBuild": "echo \"TestOK\""
	}]`
)

func TestCTXBuilder(t *testing.T) {
	var (
		mapp *mockupapp.App
		err  error
	)
	t.Parallel()
	// prepare mockup application & data
	if mapp, err = mockupapp.NewApp(mockupapp.MockupOptions{
		Input: gio.NewInput(strings.NewReader("")),
	}); err != nil {
		t.Error(err)
		return
	}
	fs := mapp.RootFilespace()
	if err = goaterr.ToErrors(goaterr.AppendError(nil,
		fs.WriteFile(".goat/build/templates/names/main.tmpl", []byte(testCTXBuilderAfterBuildTemplate), 0766),
		fs.WriteFile(".goat/build.def.json", []byte(testCTXBuilderAfterBuildConfig), 0766))); err != nil {
		t.Error(err)
		return
	}
	dp := mapp.DependencyProvider()
	if err = goaterr.ToErrors(goaterr.AppendError(nil,
		RegisterDependencies(dp),
		modules.RegisterDependencies(dp),
		dependencies.RegisterDependencies(dp),
		repositories.RegisterDependencies(dp),
		template.RegisterDependencies(dp),
		vcs.RegisterDependencies(dp),
		simpletf.RegisterFunctions(dp))); err != nil {
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
	ctxScope := scope.NewScope("test")
	appModel := am.NewApplicationModel(map[string]string{})
	appData := services.ApplicationData{
		AM:    appModel,
		Plain: map[string]string{},
	}
	buildContext := deps.BuilderService.NewContext(ctxScope, appData, map[string]string{}, map[string]string{})
	if err = buildContext.Build(fs); err != nil {
		t.Error(err)
		return
	}
	if err = ctxScope.Wait(); err != nil {
		t.Error(err)
		return
	}
	if err = ctxScope.Trigger(app.CommitEvent, nil); err != nil {
		t.Error(err)
		return
	}
	outBuffer := mapp.OutputBuffer()
	if !strings.Contains(outBuffer.String(), "TestOK") {
		t.Errorf("expected TestOK in afterBuild command output and take '%s'", outBuffer.String())
		return
	}
}
