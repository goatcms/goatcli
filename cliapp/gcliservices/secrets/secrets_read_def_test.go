package secrets

import (
	"testing"

	"github.com/goatcms/goatcli/cliapp/common/gclivarutil"

	"github.com/goatcms/goatcli/cliapp/common"

	"github.com/goatcms/goatcli/cliapp/common/am"
	"github.com/goatcms/goatcli/cliapp/common/config"
	"github.com/goatcms/goatcli/cliapp/gcliservices"
	"github.com/goatcms/goatcore/app"
)

func TestSecretsReadDefFromDataFile(t *testing.T) {
	var (
		err        error
		service    gcliservices.SecretsService
		mapp       app.App
		def        []*config.Property
		appData    gcliservices.ApplicationData
		properties common.ElasticData
	)
	t.Parallel()
	if service, mapp, err = buildMockupApp(""); err != nil {
		t.Error(err)
		return
	}
	if err = mapp.Filespaces().Root().WriteFile(SecretsDefPath, []byte(testPropDefJSON), 0766); err != nil {
		t.Error(err)
		return
	}
	if appData, err = am.NewApplicationData(map[string]string{}); err != nil {
		t.Error(err)
		return
	}
	if properties, err = gclivarutil.NewElasticData(map[string]string{}); err != nil {
		t.Error(err)
		return
	}
	if def, err = service.ReadDefFromFS(mapp.Filespaces().Root(), properties, appData); err != nil {
		t.Error(err)
		return
	}
	if len(def) != 2 {
		t.Errorf("result def should contains two elements and it have %d", len(def))
		return
	}
	if def[0].Key != "key1" {
		t.Errorf("incorrect first element key")
		return
	}
	if def[1].Key != "key2" {
		t.Errorf("incorrect second element key")
		return
	}
}

func TestPropertiesReaDefFromNoFile(t *testing.T) {
	var (
		err              error
		service          gcliservices.SecretsService
		mapp             app.App
		def              []*config.Property
		appData          gcliservices.ApplicationData
		emptyElasticData common.ElasticData
	)
	t.Parallel()
	if service, mapp, err = buildMockupApp(""); err != nil {
		t.Error(err)
		return
	}
	if appData, err = am.NewApplicationData(map[string]string{}); err != nil {
		t.Error(err)
		return
	}
	if emptyElasticData, err = gclivarutil.NewElasticData(map[string]string{}); err != nil {
		t.Error(err)
		return
	}
	if def, err = service.ReadDefFromFS(mapp.Filespaces().Root(), emptyElasticData, appData); err != nil {
		t.Error(err)
		return
	}
	if len(def) != 0 {
		t.Errorf("result def should contains zero elements and it have %d", len(def))
		return
	}
}
