package builder

import (
	"bytes"
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
	testCTXHookLayout = `{{- define ".gitignore"}}
		{{- $ctx := .}}
		#ignore these file
		/main.go
	{{- end}}`
	testCTXHookTemplate = `
	{{$ctx.RenderOnce ".gitignore" "" "" ".gitignore" $ctx.DotData}}
	`
	testCTXHookConfig = `[]`
)

func TestCTXHook(t *testing.T) {
	var (
		err     error
		context []byte
	)
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
	fs := mapp.RootFilespace()
	if err = goaterr.ToErrors(goaterr.AppendError(nil,
		fs.WriteFile(".goat/build/layouts/default/main.tmpl", []byte(testCTXHookLayout), 0766),
		fs.WriteFile(".goat/build/templates/hook/vcs/git/main.ctrl", []byte(testCTXHookTemplate), 0766),
		fs.WriteFile(".goat/build.def.json", []byte(testCTXHookConfig), 0766))); err != nil {
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
	appData := am.NewApplicationData(map[string]string{})
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
	if !fs.IsFile(".gitignore") {
		t.Errorf(".gitignore is not exist")
		return
	}
	if context, err = fs.ReadFile(".gitignore"); err != nil {
		t.Error(err)
		return
	}
	if !strings.Contains(string(context), "/main.go") {
		t.Errorf("expected '/main.go' in .gitignore file and take '%s'", context)
		return
	}
}
