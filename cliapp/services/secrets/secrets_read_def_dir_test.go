package secrets

import (
	"bytes"
	"strings"
	"testing"

	"github.com/goatcms/goatcli/cliapp/common/config"
	"github.com/goatcms/goatcli/cliapp/services"
	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/app/gio"
	"github.com/goatcms/goatcore/app/mockupapp"
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
			Secrets services.PropertiesService `dependency:"SecretsService"`
		}
	)
	t.Parallel()
	// prepare mockup application
	if mapp, err = mockupapp.NewApp(mockupapp.MockupOptions{
		Input:  gio.NewInput(strings.NewReader("")),
		Output: gio.NewOutput(new(bytes.Buffer)),
	}); err != nil {
		t.Error(err)
		return
	}
	if err = mapp.RootFilespace().WriteFile(".goat/secrets.def/01_first.json", []byte(testSecretsDefFirstJSON), 0766); err != nil {
		t.Error(err)
		return
	}
	if err = mapp.RootFilespace().WriteFile(".goat/secrets.def/02_second.json", []byte(testSecretsDefSecondJSON), 0766); err != nil {
		t.Error(err)
		return
	}
	if err = mapp.RootFilespace().WriteFile(".goat/data.def/wrong.ex", []byte("WRONG_FILE"), 0766); err != nil {
		t.Error(err)
		return
	}
	if err = RegisterDependencies(mapp.DependencyProvider()); err != nil {
		t.Error(err)
		return
	}
	if err = mapp.DependencyProvider().InjectTo(&deps); err != nil {
		t.Error(err)
		return
	}
	// test
	if secrets, err = deps.Secrets.ReadDefFromFS(mapp.RootFilespace()); err != nil {
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