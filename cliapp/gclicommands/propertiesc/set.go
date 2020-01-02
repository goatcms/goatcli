package propertiesc

import (
	"github.com/goatcms/goatcli/cliapp/gcliservices"
	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/filesystem"
	"github.com/goatcms/goatcore/varutil/goaterr"
)

// RunSetPropertyValue run command to set property value
func RunSetPropertyValue(a app.App, ctx app.IOContext) (err error) {
	var (
		deps struct {
			Key               string                         `command:"?$1"`
			Value             string                         `command:"?$2"`
			CurrentFS         filesystem.Filespace           `filespace:"current"`
			PropertiesService gcliservices.PropertiesService `dependency:"PropertiesService"`
		}
		propertiesData map[string]string
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
	if deps.Value == "" {
		return goaterr.Errorf(ValueParameterIsRequired)
	}
	if propertiesData, err = deps.PropertiesService.ReadDataFromFS(deps.CurrentFS); err != nil {
		return err
	}
	propertiesData[deps.Key] = deps.Value
	return deps.PropertiesService.WriteDataToFS(deps.CurrentFS, propertiesData)
}
