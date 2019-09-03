package builder

import (
	"bytes"
	"strings"
	"testing"

	"github.com/goatcms/goatcli/cliapp/services"
	"github.com/goatcms/goatcli/cliapp/services/dependencies"
	"github.com/goatcms/goatcli/cliapp/services/modules"
	"github.com/goatcms/goatcli/cliapp/services/repositories"
	"github.com/goatcms/goatcli/cliapp/services/template"
	"github.com/goatcms/goatcli/cliapp/services/vcs"
	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/app/gio"
	"github.com/goatcms/goatcore/app/mockupapp"
	"github.com/goatcms/goatcore/app/scope"
	"github.com/goatcms/goatcore/varutil/goaterr"
)

const (
	testCTXBuilderLayout = `{{- define "out/file.txt"}}
		{{- $ctx := .}}
		{{- index $ctx.PlainData "datakey" }}
		{{- index $ctx.Properties.Project "propkey" }}
		{{- index $ctx.Properties.Secrets "secretkey" }}
	{{- end}}`
	testCTXBuilderTemplate = `
	{{$ctx := .}}
	{{$ctx.RenderOnce "out/file.txt" "" "" "out/file.txt" $ctx.DotData}}
	`
	testCTXBuilderConfig = `[{
	  "from":"ignore",
	  "to":"ignore",
	  "template":"names",
	  "layout":"default"
	}]`
)

func TestCTXBuilderAfterBuild(t *testing.T) {
	var (
		mapp    app.App
		err     error
		context []byte
	)
	t.Parallel()
	// prepare mockup application & data
	if mapp, err = mockupapp.NewApp(mockupapp.MockupOptions{
		Input:  gio.NewInput(strings.NewReader("")),
		Output: gio.NewOutput(new(bytes.Buffer)),
	}); err != nil {
		t.Error(err)
		return
	}
	fs := mapp.RootFilespace()
	if err = goaterr.ToErrors(goaterr.AppendError(nil,
		fs.WriteFile(".goat/build/layouts/default/main.tmpl", []byte(testCTXBuilderLayout), 0766),
		fs.WriteFile(".goat/build/templates/names/main.tmpl", []byte(testCTXBuilderTemplate), 0766),
		fs.WriteFile(".goat/build.def.json", []byte(testCTXBuilderConfig), 0766))); err != nil {
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
	ctxScope := scope.NewScope("test")
	buildContext := deps.BuilderService.NewContext(ctxScope, map[string]string{
		"datakey": "Ala",
	}, map[string]string{
		"propkey": " ma",
	}, map[string]string{
		"secretkey": " kota",
	})
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
	if !fs.IsFile("out/file.txt") {
		t.Errorf("out/file.txt is not exist")
		return
	}
	if context, err = fs.ReadFile("out/file.txt"); err != nil {
		t.Error(err)
		return
	}
	if string(context) != "Ala ma kota" {
		t.Errorf("File content must be equals to 'Ala ma kota' and it is '%s'", context)
		return
	}
}
