package secrets

import (
	"strings"

	"github.com/goatcms/goatcli/cliapp/gcliservices"
	"github.com/goatcms/goatcli/cliapp/gcliservices/template"
	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/app/gio"
	"github.com/goatcms/goatcore/app/goatapp"
)

const (
	testPropDefJSON      = `[{"key":"key1", "type":"alnum", "min":1, "max":22},{"key":"key2", "type":"alnum", "min":1, "max":22}]`
	testPropDataJSON     = `{"key1":"value1"}`
	testPropFullDataJSON = `{"key1":"value1","key2":"value2"}`
)

func buildMockupApp(input string) (p gcliservices.SecretsService, mapp app.App, err error) {
	if mapp, err = goatapp.NewMockupApp(goatapp.Params{
		IO: goatapp.IO{
			In: gio.NewAppInput(strings.NewReader(input)),
		},
	}); err != nil {
		return nil, nil, err
	}
	if err = RegisterDependencies(mapp.DependencyProvider()); err != nil {
		return nil, nil, err
	}
	if err = template.RegisterDependencies(mapp.DependencyProvider()); err != nil {
		return nil, nil, err
	}
	var deps struct {
		Secrets gcliservices.SecretsService `dependency:"SecretsService"`
	}
	if err = mapp.DependencyProvider().InjectTo(&deps); err != nil {
		return nil, nil, err
	}
	return deps.Secrets, mapp, nil
}
