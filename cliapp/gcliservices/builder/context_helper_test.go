package builder

import (
	"strings"
	"testing"

	"github.com/goatcms/goatcli/cliapp/common/gclivarutil"

	"github.com/goatcms/goatcli/cliapp/common"

	"github.com/goatcms/goatcli/cliapp/common/am"
	"github.com/goatcms/goatcli/cliapp/gcliservices"
	"github.com/goatcms/goatcli/cliapp/gcliservices/dependencies"
	"github.com/goatcms/goatcli/cliapp/gcliservices/modules"
	"github.com/goatcms/goatcli/cliapp/gcliservices/repositories"
	"github.com/goatcms/goatcli/cliapp/gcliservices/template"
	"github.com/goatcms/goatcli/cliapp/gcliservices/template/simpletf"
	"github.com/goatcms/goatcli/cliapp/gcliservices/vcs"
	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/app/gio"
	"github.com/goatcms/goatcore/app/mockupapp"
	"github.com/goatcms/goatcore/varutil/goaterr"
)

func TestCTXBuilderHelper(t *testing.T) {
	var (
		mapp       app.App
		err        error
		context    []byte
		appData    gcliservices.ApplicationData
		properties common.ElasticData
		secrets    common.ElasticData
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
	if err = goaterr.ToError(goaterr.AppendError(nil,
		fs.WriteFile(".goat/build/helpers/helper.tmpl", []byte(`{{- define "helper" -}}
		Hello helper
	{{- end}}`), 0766),
		fs.WriteFile(".goat/build/layouts/default/main.tmpl", []byte(`{{- define "out/file.txt"}}
		{{- template "helper" . -}}
	{{- end}}`), 0766),
		fs.WriteFile(".goat/build/templates/names/main.ctrl", []byte(`
		{{- $ctx.RenderOnce "out/file.txt" "" "" "out/file.txt" $ctx.DotData -}}
		`), 0766),
		fs.WriteFile(".goat/build.def.json", []byte(`[{
			"from":"ignore",
			"to":"ignore",
			"template":"names",
			"layout":"default"
		  }]`), 0766))); err != nil {
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
	if appData, err = am.NewApplicationData(map[string]string{
		"datakey": "Ala",
	}); err != nil {
		t.Error(err)
		return
	}
	if properties, err = gclivarutil.NewElasticData(map[string]string{
		"propkey": " ma",
	}); err != nil {
		t.Error(err)
		return
	}
	if secrets, err = gclivarutil.NewElasticData(map[string]string{
		"secretkey": " kota",
	}); err != nil {
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
	if !fs.IsFile("out/file.txt") {
		t.Errorf("out/file.txt is not exist")
		return
	}
	if context, err = fs.ReadFile("out/file.txt"); err != nil {
		t.Error(err)
		return
	}
	if string(context) != "Hello helper" {
		t.Errorf("File content must be equals to 'Ala ma kota' and it is '%s'", context)
		return
	}
}
