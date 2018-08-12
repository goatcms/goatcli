package properties

import (
	"testing"

	"github.com/goatcms/goatcli/cliapp/common/config"
	"github.com/goatcms/goatcli/cliapp/services"
)

func TestPropertieFillNoInteractive(t *testing.T) {
	var (
		err     error
		service services.PropertiesService
		data    = map[string]string{}
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
	if _, err = service.FillData(propertiesDef, data, map[string]string{}, false); err == nil {
		t.Errorf("FillData should return error for no-interactive mode and no data")
		return
	}
}
