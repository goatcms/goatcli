package secrets

import (
	"testing"

	"github.com/goatcms/goatcli/cliapp/common/config"
	"github.com/goatcms/goatcli/cliapp/gcliservices"
	"github.com/goatcms/goatcore/app"
)

func TesSecretsFillDataFile(t *testing.T) {
	var (
		mapp      app.App
		err       error
		service   gcliservices.SecretsService
		data      = map[string]string{}
		isChanged bool
	)
	t.Parallel()
	if service, mapp, err = buildMockupApp("my_insert_value1\nmy_insert_value2"); err != nil {
		t.Error(err)
		return
	}
	propertiesDef := []*config.Property{{
		Key:  "key1",
		Type: "line",
		Min:  1,
		Max:  22,
	}, {
		Key:  "key2",
		Type: "line",
		Min:  1,
		Max:  22,
	}}
	if isChanged, err = service.FillData(mapp.IOContext(), propertiesDef, data, map[string]string{}, true); err != nil {
		t.Error(err)
		return
	}
	if isChanged == false {
		t.Errorf("data was changed and return false isChanged")
		return
	}
	if data["key1"] != "my_insert_value1" {
		t.Errorf("expect key1 value equals to my_insert_value1 and it is '%s'", data["key1"])
		return
	}
	if data["key2"] != "my_insert_value2" {
		t.Errorf("expect key2 value equals to my_insert_value2 and it is '%s'", data["key2"])
		return
	}
}

func TestPropertieFillDataFileOmnit(t *testing.T) {
	var (
		mapp    app.App
		err     error
		service gcliservices.SecretsService
		data    = map[string]string{
			"key1": "value1",
			"key2": "value2",
		}
		isChanged bool
	)
	t.Parallel()
	if service, mapp, err = buildMockupApp(""); err != nil {
		t.Error(err)
		return
	}
	if isChanged, err = service.FillData(mapp.IOContext(), []*config.Property{{
		Key:  "key1",
		Type: "alnum",
		Min:  1,
		Max:  22,
	}, {
		Key:  "key2",
		Type: "alnum",
		Min:  1,
		Max:  22,
	}}, data, map[string]string{}, true); err != nil {
		t.Error(err)
		return
	}
	if isChanged == true {
		t.Errorf("data is not modyfied and isChanges is true")
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
