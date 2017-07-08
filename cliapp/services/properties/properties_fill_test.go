package properties

import (
	"testing"

	"github.com/goatcms/goatcli/cliapp/common/config"
	"github.com/goatcms/goatcli/cliapp/services"
)

func TestPropertieFillDataFile(t *testing.T) {
	var (
		err       error
		service   services.PropertiesService
		data      = map[string]string{}
		isChanged bool
	)
	t.Parallel()
	if service, _, err = buildMockupApp("my_insert_value1\nmy_insert_value2"); err != nil {
		t.Error(err)
		return
	}
	propertiesDef := []*config.Property{
		&config.Property{
			Key:  "key1",
			Type: "line",
			Min:  1,
			Max:  22,
		},
		&config.Property{
			Key:  "key2",
			Type: "line",
			Min:  1,
			Max:  22,
		}}
	if isChanged, err = service.FillData(propertiesDef, data, map[string]string{}); err != nil {
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
		err     error
		service services.PropertiesService
		data    = map[string]string{
			"key1": "value1",
			"key2": "value2",
		}
		isChanged bool
	)
	t.Parallel()
	if service, _, err = buildMockupApp(""); err != nil {
		t.Error(err)
		return
	}
	if isChanged, err = service.FillData([]*config.Property{
		&config.Property{
			Key:  "key1",
			Type: "alnum",
			Min:  1,
			Max:  22,
		},
		&config.Property{
			Key:  "key2",
			Type: "alnum",
			Min:  1,
			Max:  22,
		}}, data, map[string]string{}); err != nil {
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
