package properties

import (
	"testing"

	"github.com/goatcms/goatcli/cliapp/common/config"
	"github.com/goatcms/goatcli/cliapp/services"
	"github.com/goatcms/goatcore/app"
)

func TestPropertiesReadDefFromDataFile(t *testing.T) {
	var (
		err     error
		service services.PropertiesService
		mapp    app.App
		def     []*config.Property
	)
	t.Parallel()
	if service, mapp, err = buildMockupApp(""); err != nil {
		t.Error(err)
		return
	}
	if err = mapp.RootFilespace().WriteFile(PropertiesDefPath, []byte(testPropDefJSON), 0766); err != nil {
		t.Error(err)
		return
	}
	if def, err = service.ReadDefFromFS(mapp.RootFilespace()); err != nil {
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
		err     error
		service services.PropertiesService
		mapp    app.App
		def     []*config.Property
	)
	t.Parallel()
	if service, mapp, err = buildMockupApp(""); err != nil {
		t.Error(err)
		return
	}
	if def, err = service.ReadDefFromFS(mapp.RootFilespace()); err != nil {
		t.Error(err)
		return
	}
	if len(def) != 0 {
		t.Errorf("result def should contains zero elements and it have %d", len(def))
		return
	}
}
