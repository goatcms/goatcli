package secrets

import (
	"testing"

	"github.com/goatcms/goatcli/cliapp/common"

	"github.com/goatcms/goatcli/cliapp/common/am"
	"github.com/goatcms/goatcli/cliapp/common/config"
	"github.com/goatcms/goatcli/cliapp/common/gclivarutil"
	"github.com/goatcms/goatcli/cliapp/gcliservices"
	"github.com/goatcms/goatcli/cliapp/gcliservices/template"
	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/app/goatapp"
)

const (
	testSecretsDefFirstJSON  = `[{"key":"key1", "type":"alnum", "min":1, "max":22},{"key":"key2", "type":"alnum", "min":1, "max":22}]`
	testSecretsDefSecondJSON = `[{"key":"key3", "type":"alnum", "min":1, "max":22},{"key":"key4", "type":"alnum", "min":1, "max":22}]`
)

func TestDataDefFromDirectory(t *testing.T) {
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
	if mapp, err = goatapp.NewMockupApp(goatapp.Params{}); err != nil {
		t.Error(err)
		return
	}
	if err = mapp.Filespaces().Root().WriteFile(".goat/secrets.def/01_first.json", []byte(testSecretsDefFirstJSON), 0766); err != nil {
		t.Error(err)
		return
	}
	if err = mapp.Filespaces().Root().WriteFile(".goat/secrets.def/02_second.json", []byte(testSecretsDefSecondJSON), 0766); err != nil {
		t.Error(err)
		return
	}
	if err = mapp.Filespaces().Root().WriteFile(".goat/data.def/wrong.ex", []byte("WRONG_FILE"), 0766); err != nil {
		t.Error(err)
		return
	}
	if err = RegisterDependencies(mapp.DependencyProvider()); err != nil {
		t.Error(err)
		return
	}
	if err = template.RegisterDependencies(mapp.DependencyProvider()); err != nil {
		t.Error(err)
		return
	}
	if err = mapp.DependencyProvider().InjectTo(&deps); err != nil {
		t.Error(err)
		return
	}
	// test
	if appData, err = am.NewApplicationData(map[string]string{}); err != nil {
		t.Error(err)
		return
	}
	if emptyElasticData, err = gclivarutil.NewElasticData(map[string]string{}); err != nil {
		t.Error(err)
		return
	}
	if secrets, err = deps.Secrets.ReadDefFromFS(mapp.Filespaces().Root(), emptyElasticData, appData); err != nil {
		t.Error(err)
		return
	}
	if len(secrets) != 4 {
		t.Errorf("expected four secrets definitions and take %d", len(secrets))
		return
	}
	// check order
	if secrets[0].Key != "key1" {
		t.Errorf("secrets[0].Key should be equals to key1 and take %v", secrets[0].Key)
		return
	}
	if secrets[1].Key != "key2" {
		t.Errorf("secrets[1].Key should be equals to key1 and take %v", secrets[1].Key)
		return
	}
	if secrets[2].Key != "key3" {
		t.Errorf("secrets[2].Key should be equals to key1 and take %v", secrets[2].Key)
		return
	}
	if secrets[3].Key != "key4" {
		t.Errorf("secrets[3].Keyshould be equals to key1 and take %v", secrets[3].Key)
		return
	}
}
