package properties

import (
	"strings"

	"github.com/goatcms/goatcli/cliapp/gclimock"
	"github.com/goatcms/goatcli/cliapp/gcliservices"
	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/app/gio"
	"github.com/goatcms/goatcore/app/mockupapp"
)

const (
	testPropDefJSON      = `[{"key":"key1", "type":"alnum", "min":1, "max":22},{"key":"key2", "type":"alnum", "min":1, "max":22}]`
	testPropDataJSON     = `{"key1":"value1"}`
	testPropFullDataJSON = `{"key1":"value1","key2":"value2"}`
)

func buildMockupApp(input string) (p gcliservices.PropertiesService, mapp app.App, err error) {
	if mapp, err = gclimock.NewApp(mockupapp.MockupOptions{
		Input: gio.NewInput(strings.NewReader(input)),
	}); err != nil {
		return nil, nil, err
	}
	if err = RegisterDependencies(mapp.DependencyProvider()); err != nil {
		return nil, nil, err
	}
	var deps struct {
		Properties gcliservices.PropertiesService `dependency:"PropertiesService"`
	}
	if err = mapp.DependencyProvider().InjectTo(&deps); err != nil {
		return nil, nil, err
	}
	return deps.Properties, mapp, nil
}
