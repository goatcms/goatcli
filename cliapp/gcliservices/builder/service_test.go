package builder

import (
	"strings"
	"testing"

	"github.com/goatcms/goatcli/cliapp/common"
	"github.com/goatcms/goatcli/cliapp/common/am"
	"github.com/goatcms/goatcli/cliapp/common/gclivarutil"
	"github.com/goatcms/goatcli/cliapp/gcliservices"
	"github.com/goatcms/goatcli/cliapp/gcliservices/data"
	"github.com/goatcms/goatcli/cliapp/gcliservices/dependencies"
	"github.com/goatcms/goatcli/cliapp/gcliservices/modules"
	"github.com/goatcms/goatcli/cliapp/gcliservices/repositories"
	"github.com/goatcms/goatcli/cliapp/gcliservices/template"
	"github.com/goatcms/goatcli/cliapp/gcliservices/template/simpletf"
	"github.com/goatcms/goatcli/cliapp/gcliservices/vcs"
	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/app/goatapp"
	"github.com/goatcms/goatcore/varutil/goaterr"
)

const (
	testBuilderLayout = `{{- define "out/file.txt" -}}
		File Content
	{{- end -}}`
	testBuilderTemplate = `
	{{$ctx.RenderOnce "out/file.txt" "" "" "out/file.txt" $ctx.DotData}}`
	testBuilderConfig = `[{
	  "from":"ignore",
	  "to":"ignore",
	  "template":"names",
	  "layout":"default"
	}]`
)

func TestBuilder(t *testing.T) {
	var (
		mapp             app.App
		err              error
		context          []byte
		ctx              app.IOContext
		appData          gcliservices.ApplicationData
		emptyElasticData common.ElasticData
	)
	t.Parallel()
	// prepare mockup application & data
	if mapp, err = goatapp.NewMockupApp(goatapp.Params{}); err != nil {
		t.Error(err)
		return
	}
	fs := mapp.Filespaces().CWD()
	if err = goaterr.ToError(goaterr.AppendError(nil,
		fs.WriteFile(".goat/build/layouts/default/main.tmpl", []byte(testBuilderLayout), 0766),
		fs.WriteFile(".goat/build/templates/names/main.ctrl", []byte(testBuilderTemplate), 0766),
		fs.WriteFile(".goat/build.def.json", []byte(testBuilderConfig), 0766))); err != nil {
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
		data.RegisterDependencies(dp),
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
	if ctx, err = newEmptyIOContext(); err != nil {
		t.Error(err)
		return
	}
	if appData, err = am.NewApplicationData(map[string]string{}); err != nil {
		t.Error(err)
		return
	}
	if emptyElasticData, err = gclivarutil.NewElasticData(map[string]string{}); err != nil {
		t.Error(err)
		return
	}
	if err = deps.BuilderService.Build(ctx, fs, appData, emptyElasticData, emptyElasticData); err != nil {
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
	if !fs.IsFile("out/file.txt") {
		t.Errorf("out/file.txt is not exist")
		return
	}
	if context, err = fs.ReadFile("out/file.txt"); err != nil {
		t.Error(err)
		return
	}
	if !strings.Contains(string(context), "File Content") {
		t.Errorf("File must contains 'File Content' and it is '%s'", context)
		return
	}
}
