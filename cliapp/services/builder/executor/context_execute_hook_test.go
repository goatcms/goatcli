package executor

import (
	"bytes"
	"strings"
	"testing"

	"github.com/goatcms/goatcli/cliapp/services"
	templates "github.com/goatcms/goatcli/cliapp/services/template"
	"github.com/goatcms/goatcli/cliapp/services/template/simpletf"
	"github.com/goatcms/goatcli/cliapp/services/vcs"
	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/app/gio"
	"github.com/goatcms/goatcore/app/mockupapp"
	"github.com/goatcms/goatcore/app/scope"
	"github.com/goatcms/goatcore/filesystem"
	"github.com/goatcms/goatcore/varutil/goaterr"
	"github.com/goatcms/goatcore/workers"
)

const (
	masterHook = `
		{{ $ctx := . }}
		{{$ctx.RenderOnce "dir/result.txt" "" "" "testf" $ctx.DotData}}
	`
	testfHook = `
	{{define "testf"}}
		{{ $ctx := .}}
		Names:{{block "list" $ctx.DotData}}{{"\n"}}{{range .}}{{println "-" .}}{{end}}{{end}}
	{{end}}`
	overlayHook = `{{define "list"}} {{join . ", "}}{{end}} `
)

func TestContextExecute(t *testing.T) {
	t.Parallel()
	var (
		guardians         = []string{"Gamora", "Groot", "Nebula", "Rocket", "Star-Lord"}
		mapp              app.App
		fs                filesystem.Filespace
		resultBytes       []byte
		result            string
		templateExecutor  services.TemplateExecutor
		generatorExecutor *GeneratorExecutor
		generatorScope    = scope.NewScope("generator")
		err               error
	)
	for ti := 0; ti < workers.AsyncTestReapeat; ti++ {
		// prepare mockup application
		if mapp, err = mockupapp.NewApp(mockupapp.MockupOptions{
			Input:  gio.NewInput(strings.NewReader("")),
			Output: gio.NewOutput(new(bytes.Buffer)),
		}); err != nil {
			t.Error(err)
			return
		}
		fs = mapp.RootFilespace()
		if err = goaterr.ToErrors(goaterr.AppendError(nil,
			fs.WriteFile(".goat/build/templates/hook/testHook/git/master.tmpl", []byte(masterHook), filesystem.DefaultUnixFileMode),
			fs.WriteFile(".goat/build/templates/hook/testHook/git/testf.tmpl", []byte(testfHook), filesystem.DefaultUnixFileMode),
			fs.WriteFile(".goat/build/templates/hook/testHook/git/overlay.tmpl", []byte(overlayHook), filesystem.DefaultUnixFileMode),
			templates.RegisterDependencies(mapp.DependencyProvider()),
			simpletf.RegisterFunctions(mapp.DependencyProvider()),
		)); err != nil {
			t.Error(err)
			return
		}
		// test
		var deps struct {
			TemplateService services.TemplateService `dependency:"TemplateService"`
		}
		if err = mapp.DependencyProvider().InjectTo(&deps); err != nil {
			t.Error(err)
			return
		}
		if templateExecutor, err = deps.TemplateService.Build(fs); err != nil {
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
			VCSData: vcs.NewData(vcs.NewGeneratedFiles(true), vcs.NewIgnoredFiles(true)),
		}, 10, templateExecutor); err != nil {
			t.Error(err)
			return
		}
		if err = generatorExecutor.ExecuteHook("testHook", guardians); err != nil {
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
