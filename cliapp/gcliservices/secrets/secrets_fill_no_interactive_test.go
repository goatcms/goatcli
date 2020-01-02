package secrets

import (
	"testing"

	"github.com/goatcms/goatcli/cliapp/common/config"
	"github.com/goatcms/goatcli/cliapp/gcliservices"
	"github.com/goatcms/goatcore/app"
)

func TesSecretsFillNoInteractive(t *testing.T) {
	var (
		mapp    app.App
		err     error
		service gcliservices.SecretsService
		data    = map[string]string{}
	)
	t.Parallel()
	if service, mapp, err = buildMockupApp("my_insert_value1\nmy_insert_value2"); err != nil {
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
	if _, err = service.FillData(mapp.IOContext(), propertiesDef, data, map[string]string{}, false); err == nil {
		t.Errorf("FillData should return error for no-interactive mode and no data")
		return
	}
}
