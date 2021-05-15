package builder

import (
	"strings"
	"testing"

	"github.com/goatcms/goatcli/cliapp/common"
	"github.com/goatcms/goatcli/cliapp/common/am"
	"github.com/goatcms/goatcli/cliapp/common/gclivarutil"
	"github.com/goatcms/goatcli/cliapp/gcliservices"
	"github.com/goatcms/goatcli/cliapp/gcliservices/dependencies"
	"github.com/goatcms/goatcli/cliapp/gcliservices/modules"
	"github.com/goatcms/goatcli/cliapp/gcliservices/repositories"
	"github.com/goatcms/goatcli/cliapp/gcliservices/template"
	"github.com/goatcms/goatcli/cliapp/gcliservices/template/simpletf"
	"github.com/goatcms/goatcli/cliapp/gcliservices/vcs"
	"github.com/goatcms/goatcore/app/goatapp"
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
		err        error
		context    []byte
		appData    gcliservices.ApplicationData
		properties common.ElasticData
		secrets    common.ElasticData
		mapp       *goatapp.MockupApp
	)
	t.Parallel()
	// prepare mockup application & data
	if mapp, err = goatapp.NewMockupApp(goatapp.Params{}); err != nil {
		t.Error(err)
		return
	}
	fs := mapp.Filespaces().CWD()
	if err = goaterr.ToError(goaterr.AppendError(nil,
		fs.WriteFile(".goat/build/layouts/default/main.tmpl", []byte(testCTXHookLayout), 0766),
		fs.WriteFile(".goat/build/templates/hook/vcs/git/main.ctrl", []byte(testCTXHookTemplate), 0766),
		fs.WriteFile(".goat/build.def.json", []byte(testCTXHookConfig), 0766))); err != nil {
		t.Error(err)
		return
	}
	dp := mapp.DependencyProvider()
	if err = goaterr.ToError(goaterr.AppendError(nil,
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
		BuilderService gcliservices.BuilderService `dependency:"BuilderService"`
	}
	if err = mapp.DependencyProvider().InjectTo(&deps); err != nil {
		t.Error(err)
		return
	}
	ctx := mapp.IOContext()
	if appData, err = am.NewApplicationData(map[string]string{}); err != nil {
		t.Error(err)
		return
	}
	if properties, err = gclivarutil.NewElasticData(map[string]string{}); err != nil {
		t.Error(err)
		return
	}
	if secrets, err = gclivarutil.NewElasticData(map[string]string{}); err != nil {
		t.Error(err)
		return
	}
	if err = deps.BuilderService.Build(ctx, fs, appData, properties, secrets); err != nil {
		t.Error(err)
		return
	}
	if err = ctx.Scope().Wait(); err != nil {
		t.Error(err)
		return
	}
	if err = ctx.Scope().Trigger(gcliservices.BuildCommitevent, nil); err != nil {
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
