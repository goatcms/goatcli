package secrets

import (
	"testing"

	"github.com/goatcms/goatcli/cliapp/common/config"
	"github.com/goatcms/goatcli/cliapp/services"
	"github.com/goatcms/goatcore/app"
)

func TestSecretsReadDataFromDataFile(t *testing.T) {
	var (
		err     error
		service services.PropertiesService
		mapp    app.App
		data    map[string]string
	)
	t.Parallel()
	if service, mapp, err = buildMockupApp(""); err != nil {
		t.Error(err)
		return
	}
	if err = mapp.RootFilespace().WriteFile(SecretsDataPath, []byte(testPropFullDataJSON), 0766); err != nil {
		t.Error(err)
		return
	}
	if data, err = service.ReadDataFromFS(mapp.RootFilespace()); err != nil {
		t.Error(err)
		return
	}
	if len(data) != 2 {
		t.Errorf("result data should contains two elements and it have %d", len(data))
		return
	}
	if data["key1"] != "value1" {
		t.Errorf("incorrect key1 value")
		return
	}
	if data["key2"] != "value2" {
		t.Errorf("incorrect key2 value")
		return
	}
}

func TestPropertiesReadDataFromNoFile(t *testing.T) {
	var (
		err     error
		service services.PropertiesService
		mapp    app.App
		data    map[string]string
		def     []*config.Property
	)
	t.Parallel()
	if service, mapp, err = buildMockupApp(""); err != nil {
		t.Error(err)
		return
	}
	if data, err = service.ReadDataFromFS(mapp.RootFilespace()); err != nil {
		t.Error(err)
		return
	}
	if len(data) != 0 {
		t.Errorf("result data should contains zero elements and it have %d", len(def))
		return
	}
}
