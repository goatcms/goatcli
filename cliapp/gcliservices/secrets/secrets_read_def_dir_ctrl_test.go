package secrets

import (
	"strings"
	"testing"

	"github.com/goatcms/goatcli/cliapp/common/gclivarutil"

	"github.com/goatcms/goatcli/cliapp/common"

	"github.com/goatcms/goatcli/cliapp/common/am"
	"github.com/goatcms/goatcli/cliapp/common/config"
	"github.com/goatcms/goatcli/cliapp/gcliservices"
	"github.com/goatcms/goatcli/cliapp/gcliservices/template"
	"github.com/goatcms/goatcli/cliapp/gcliservices/template/amtf"
	"github.com/goatcms/goatcli/cliapp/gcliservices/template/simpletf"
	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/app/gio"
	"github.com/goatcms/goatcore/app/mockupapp"
	"github.com/goatcms/goatcore/varutil/goaterr"
)

func TestDataDefFromCtrl(t *testing.T) {
	var (
		err     error
		mapp    app.App
		secrets []*config.Property
		deps    struct {
			Secrets gcliservices.SecretsService `dependency:"SecretsService"`
		}
		appData          gcliservices.ApplicationData
		emptyElasticData common.ElasticData
	)
	t.Parallel()
	// prepare mockup application
	if mapp, err = mockupapp.NewApp(mockupapp.MockupOptions{
		Input: gio.NewInput(strings.NewReader("")),
	}); err != nil {
		t.Error(err)
		return
	}
	if err = mapp.RootFilespace().WriteFile(".goat/secrets.def/script.ctrl", []byte(`
		{{- range $i, $key := (keys $ctx.Data.Plain "app.") -}}
		  {{- $ctx.AddSecret (dict "Key" (print "app." $key ".secret") "Type" "line") -}}
		{{- end -}}
	`), 0766); err != nil {
		t.Error(err)
		return
	}
	dp := mapp.DependencyProvider()
	if err = goaterr.ToError(goaterr.AppendError(nil,
		RegisterDependencies(dp),
		template.RegisterDependencies(mapp.DependencyProvider()),
		simpletf.RegisterFunctions(dp),
		amtf.RegisterFunctions(dp))); err != nil {
		t.Error(err)
		return
	}
	if err = mapp.DependencyProvider().InjectTo(&deps); err != nil {
		t.Error(err)
		return
	}
	// test
	if appData, err = am.NewApplicationData(map[string]string{
		"app.first.name":  "FirstName",
		"app.second.name": "SecondName",
	}); err != nil {
		t.Error(err)
		return
	}
	if emptyElasticData, err = gclivarutil.NewElasticData(map[string]string{}); err != nil {
		t.Error(err)
		return
	}
	if secrets, err = deps.Secrets.ReadDefFromFS(mapp.RootFilespace(), emptyElasticData, appData); err != nil {
		t.Error(err)
		return
	}
	if len(secrets) != 2 {
		t.Errorf("expected two secret definition and take %v", len(secrets))
		return
	}
}
