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
	"github.com/goatcms/goatcore/app/gio"
	"github.com/goatcms/goatcore/app/mockupapp"
	"github.com/goatcms/goatcore/app/scope"
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
)

func TestCTXBuilder(t *testing.T) {
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
	if err = mapp.RootFilespace().WriteFile(".goat/build/layouts/default/main.tmpl", []byte(testCTXBuilderLayout), 0766); err != nil {
		t.Error(err)
		return
	}
	if err = mapp.RootFilespace().WriteFile(".goat/build/templates/names/main.tmpl", []byte(testCTXBuilderTemplate), 0766); err != nil {
		t.Error(err)
		return
	}
	if err = RegisterDependencies(mapp.DependencyProvider()); err != nil {
		t.Error(err)
		return
	}
	if err = modules.RegisterDependencies(mapp.DependencyProvider()); err != nil {
		t.Error(err)
		return
	}
	if err = dependencies.RegisterDependencies(mapp.DependencyProvider()); err != nil {
		t.Error(err)
		return
	}
	if err = repositories.RegisterDependencies(mapp.DependencyProvider()); err != nil {
		t.Error(err)
		return
	}
	if err = template.RegisterDependencies(mapp.DependencyProvider()); err != nil {
		t.Error(err)
		return
	}
	if err = template.InitDependencies(mapp); err != nil {
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
	buildConfig := []*config.Build{
		&config.Build{
			From:     "ignore",
			To:       "ignore",
			Template: "names",
			Layout:   "default",
		},
	}
	fs := mapp.RootFilespace()
	ctxScope := scope.NewScope("test")
	if err = deps.BuilderService.Build(ctxScope, fs, buildConfig, map[string]string{
		"datakey": "Ala",
	}, map[string]string{
		"propkey": " ma",
	}, map[string]string{
		"secretkey": " kota",
	}); err != nil {
		t.Error(err)
		return
	}
	if err = ctxScope.Wait(); err != nil {
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
