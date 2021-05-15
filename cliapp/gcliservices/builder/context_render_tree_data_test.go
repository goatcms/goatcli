package builder

import (
	"testing"

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

func TestRenderDataTreeFile(t *testing.T) {
	var (
		mapp    app.App
		err     error
		context []byte
	)
	t.Parallel()
	// prepare mockup application & data
	if mapp, err = goatapp.NewMockupApp(goatapp.Params{}); err != nil {
		t.Error(err)
		return
	}
	fs := mapp.Filespaces().CWD()
	if err = goaterr.ToError(goaterr.AppendError(nil,
		fs.WriteFile(".goat/build/templates/names/file.txt.render", []byte(`
		{{- $ctx.Data.Tree.datakey }}
		{{- $ctx.Properties.Project.Tree.propkey -}}
		{{- $ctx.Properties.Secrets.Tree.secretkey -}}
		`), 0766),
		fs.WriteFile(".goat/build.def.json", []byte(`[{
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
	// It must be render evry time
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
	if string(context) != "Ala nie ma kota" {
		t.Errorf("File content must be equals to 'Ala nie ma kota' and it is '%s' after re-render/re-build", context)
		return
	}
}
