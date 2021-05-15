package properties

import (
	"testing"

	"github.com/goatcms/goatcli/cliapp/gcliservices"
	"github.com/goatcms/goatcore/app"
)

func TestPropertiesWriteDataFromDataFile(t *testing.T) {
	var (
		err     error
		service gcliservices.PropertiesService
		mapp    app.App
		data    map[string]string
	)
	t.Parallel()
	if service, mapp, err = buildMockupApp(""); err != nil {
		t.Error(err)
		return
	}
	if err = service.WriteDataToFS(mapp.Filespaces().Root(), map[string]string{
		"key1": "value1",
		"key2": "value2",
	}); err != nil {
		t.Error(err)
		return
	}
	if data, err = service.ReadDataFromFS(mapp.Filespaces().Root()); err != nil {
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
