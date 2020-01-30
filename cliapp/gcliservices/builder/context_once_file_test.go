package builder

import (
	"strings"
	"testing"

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

const (
	tesOnceFileContent = `{{- $ctx := .}}
		{{- index $ctx.PlainData "datakey" }}
		{{- index $ctx.Properties.Project "propkey" }}
		{{- index $ctx.Properties.Secrets "secretkey" }}`
	tesOnceFileConfig = `[{
		  "template":"names",
		  "layout":"default"
		}]`
)

func TestCtrlFile(t *testing.T) {
	var (
		mapp    app.App
		err     error
		context []byte
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
		fs.WriteFile(".goat/build/templates/names/file.txt.once", []byte(tesOnceFileContent), 0766),
		fs.WriteFile(".goat/build.def.json", []byte(tesOnceFileConfig), 0766))); err != nil {
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
	var deps RenderFileDeps
	if err = mapp.DependencyProvider().InjectTo(&deps); err != nil {
		t.Error(err)
		return
	}
	if err = renderFile(fs, deps, map[string]string{
		"datakey": "Ala",
	}, map[string]string{
		"propkey": " ma",
	}, map[string]string{
		"secretkey": " kota",
	}); err != nil {
		t.Error(err)
		return
	}
	if !fs.IsFile("file.txt") {
		t.Errorf("file.txt is not exist")
		return
	}
	if context, err = fs.ReadFile("file.txt"); err != nil {
		t.Error(err)
		return
	}
	if string(context) != "Ala ma kota" {
		t.Errorf("File content must be equals to 'Ala ma kota' and it is '%s'", context)
		return
	}
	// It must be render once
	if err = renderFile(fs, deps, map[string]string{
		"datakey": "Ala",
	}, map[string]string{
		"propkey": " nie ma",
	}, map[string]string{
		"secretkey": " kota",
	}); err != nil {
		t.Error(err)
		return
	}
	if !fs.IsFile("file.txt") {
		t.Errorf("file.txt is not exist")
		return
	}
	if context, err = fs.ReadFile("file.txt"); err != nil {
		t.Error(err)
		return
	}
	if string(context) != "Ala ma kota" {
		t.Errorf("File content must be equals to 'Ala ma kota' and it is '%s'. The context can not cahnge after re-render/re-build", context)
		return
	}
}
