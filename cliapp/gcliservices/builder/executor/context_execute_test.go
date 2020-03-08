package executor

import (
	"strings"
	"testing"

	"github.com/goatcms/goatcli/cliapp/gcliservices"
	templates "github.com/goatcms/goatcli/cliapp/gcliservices/template"
	"github.com/goatcms/goatcli/cliapp/gcliservices/template/simpletf"
	"github.com/goatcms/goatcli/cliapp/gcliservices/vcs"
	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/app/gio"
	"github.com/goatcms/goatcore/app/mockupapp"
	"github.com/goatcms/goatcore/app/scope"
	"github.com/goatcms/goatcore/filesystem"
	"github.com/goatcms/goatcore/varutil/goaterr"
	"github.com/goatcms/goatcore/workers"
)

const (
	master = `
		{{ $ctx := . }}
		{{$ctx.RenderOnce "dir/result.txt" "" "testTemplate" "testf" $ctx.DotData}}
	`
	testf = `
	{{define "testf"}}
		{{ $ctx := .}}
		Names:{{block "list" $ctx.DotData}}{{"\n"}}{{range .}}{{println "-" .}}{{end}}{{end}}
	{{end}}`
	overlay = `{{define "list"}} {{join . ", "}}{{end}} `
)

func TestContextExecuteHook(t *testing.T) {
	t.Parallel()
	var (
		guardians         = []string{"Gamora", "Groot", "Nebula", "Rocket", "Star-Lord"}
		mapp              app.App
		fs                filesystem.Filespace
		resultBytes       []byte
		result            string
		templatesExecutor gcliservices.TemplatesExecutor
		generatorExecutor *GeneratorExecutor
		generatorScope    = scope.NewScope(scope.Params{})
		err               error
	)
	for ti := 0; ti < workers.AsyncTestReapeat; ti++ {
		// prepare mockup application
		if mapp, err = mockupapp.NewApp(mockupapp.MockupOptions{
			Input: gio.NewInput(strings.NewReader("")),
		}); err != nil {
			t.Error(err)
			return
		}
		fs = mapp.RootFilespace()
		if err = goaterr.ToError(goaterr.AppendError(nil,
			fs.WriteFile(".goat/build/templates/testTemplate/master.tmpl", []byte(master), filesystem.DefaultUnixFileMode),
			fs.WriteFile(".goat/build/templates/testTemplate/testf.tmpl", []byte(testf), filesystem.DefaultUnixFileMode),
			fs.WriteFile(".goat/build/templates/testTemplate/overlay.tmpl", []byte(overlay), filesystem.DefaultUnixFileMode),
			templates.RegisterDependencies(mapp.DependencyProvider()),
			simpletf.RegisterFunctions(mapp.DependencyProvider()),
		)); err != nil {
			t.Error(err)
			return
		}
		// test
		var deps struct {
			TemplateService gcliservices.TemplateService `dependency:"TemplateService"`
		}
		if err = mapp.DependencyProvider().InjectTo(&deps); err != nil {
			t.Error(err)
			return
		}
		if templatesExecutor, err = deps.TemplateService.TemplatesExecutor(); err != nil {
			t.Error(err)
			return
		}

		if generatorExecutor, err = NewGeneratorExecutor(generatorScope, SharedData{
			PlainData: map[string]string{},
			Properties: GlobalProperties{
				Project: map[string]string{},
				Secrets: map[string]string{},
			},
			FS:      fs,
			VCSData: vcs.NewData(vcs.NewGeneratedFiles(true), vcs.NewPersistedFiles(true)),
		}, 10, templatesExecutor); err != nil {
			t.Error(err)
			return
		}
		if err = generatorExecutor.ExecuteTask(Task{
			Template: TemplateHandler{
				Path: "testTemplate",
			},
			DotData:         guardians,
			BuildProperties: map[string]string{},
			FSPath:          "",
		}); err != nil {
			t.Error(err)
			return
		}
		if err = generatorScope.Wait(); err != nil {
			t.Error(err)
			return
		}
		if resultBytes, err = fs.ReadFile("dir/result.txt"); err != nil {
			t.Error(err)
			return
		}
		result = string(resultBytes)
		if !strings.Contains(result, "Gamora,") {
			t.Errorf("after render result file should contains 'Gamora,' and other heroes. It is: %s", result)
			return
		}
	}
}
