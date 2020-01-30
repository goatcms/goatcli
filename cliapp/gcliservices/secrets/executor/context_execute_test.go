package executor

import (
	"strings"
	"testing"

	"github.com/goatcms/goatcli/cliapp/common/am"
	"github.com/goatcms/goatcli/cliapp/common/config"
	"github.com/goatcms/goatcli/cliapp/gcliservices"
	templates "github.com/goatcms/goatcli/cliapp/gcliservices/template"
	"github.com/goatcms/goatcli/cliapp/gcliservices/template/simpletf"
	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/app/gio"
	"github.com/goatcms/goatcore/app/mockupapp"
	"github.com/goatcms/goatcore/app/scope"
	"github.com/goatcms/goatcore/filesystem"
	"github.com/goatcms/goatcore/varutil/goaterr"
	"github.com/goatcms/goatcore/workers"
)

func TestContextExecute(t *testing.T) {
	t.Parallel()
	var (
		guardians        = []string{"Gamora", "Groot", "Nebula", "Rocket", "Star-Lord"}
		mapp             app.App
		fs               filesystem.Filespace
		templateExecutor gcliservices.TemplateExecutor
		secretsExecutor  *SecretsExecutor
		executorScope    = scope.NewScope(scope.Params{})
		err              error
		secrets          []*config.Property
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
			fs.WriteFile(".goat/secrets.def/testf.ctrl", []byte(`
				{{- $ctx.AddSecret (dict "Key" "Key" "Type" "Type") -}}
			`), filesystem.DefaultUnixFileMode),
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
		if templateExecutor, err = deps.TemplateService.TemplateExecutor(".goat/secrets.def"); err != nil {
			t.Error(err)
			return
		}
		if secretsExecutor, err = NewSecretsExecutor(executorScope, SharedData{
			AppData: am.NewApplicationData(map[string]string{}),
			Properties: GlobalProperties{
				Project: map[string]string{},
			},
			DotData: guardians,
		}, 10, templateExecutor); err != nil {
			t.Error(err)
			return
		}
		if err = secretsExecutor.Execute(); err != nil {
			t.Error(err)
			return
		}
		if err = executorScope.Wait(); err != nil {
			t.Error(err)
			return
		}
		if secrets, err = secretsExecutor.Secrets(); err != nil {
			t.Error(err)
			return
		}
		if len(secrets) != 1 {
			t.Errorf("expected one generated secret definition")
			return
		}
	}
}
