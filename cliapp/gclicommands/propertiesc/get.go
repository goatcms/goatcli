package propertiesc

import (
	"github.com/goatcms/goatcli/cliapp/gcliservices"
	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/filesystem"
	"github.com/goatcms/goatcore/varutil/goaterr"
)

// RunGetPropertyValue run command to set property value
func RunGetPropertyValue(a app.App, ctx app.IOContext) (err error) {
	var (
		deps struct {
			Key               string                         `command:"?$1"`
			CurrentFS         filesystem.Filespace           `filespace:"current"`
			PropertiesService gcliservices.PropertiesService `dependency:"PropertiesService"`
		}
		propertiesData map[string]string
		ok             bool
		value          string
		ctxScope       = ctx.Scope()
	)
	if err = a.DependencyProvider().InjectTo(&deps); err != nil {
		return err
	}
	if err = ctxScope.InjectTo(&deps); err != nil {
		return err
	}
	if deps.Key == "" {
		return goaterr.Errorf(FirstKeyParameterIsRequired)
	}
	if propertiesData, err = deps.PropertiesService.ReadDataFromFS(deps.CurrentFS); err != nil {
		return err
	}
	if value, ok = propertiesData[deps.Key]; !ok {
		return goaterr.Errorf("Unknow Value for %s key", deps.Key)
	}
	return ctx.IO().Out().Printf(value)
}
