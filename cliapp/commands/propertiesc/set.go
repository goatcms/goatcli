package propertiesc

import (
	"fmt"

	"github.com/goatcms/goatcli/cliapp/services"
	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/filesystem"
)

// RunSetPropertyValue run command to set property value
func RunSetPropertyValue(a app.App, ctxScope app.Scope) (err error) {
	var (
		deps struct {
			Key               string                     `command:"?$1"`
			Value             string                     `command:"?$2"`
			CurrentFS         filesystem.Filespace       `filespace:"current"`
			PropertiesService services.PropertiesService `dependency:"PropertiesService"`
			Input             app.Input                  `dependency:"InputService"`
			Output            app.Output                 `dependency:"OutputService"`
		}
		propertiesData map[string]string
	)
	if err = a.DependencyProvider().InjectTo(&deps); err != nil {
		return err
	}
	if err = ctxScope.InjectTo(&deps); err != nil {
		return err
	}
	if deps.Key == "" {
		deps.Output.Printf(FirstKeyParameterIsRequired)
		return fmt.Errorf(FirstKeyParameterIsRequired)
	}
	if deps.Value == "" {
		deps.Output.Printf(ValueParameterIsRequired)
		return fmt.Errorf(ValueParameterIsRequired)
	}
	if propertiesData, err = deps.PropertiesService.ReadDataFromFS(deps.CurrentFS); err != nil {
		return err
	}
	propertiesData[deps.Key] = deps.Value
	if err = deps.PropertiesService.WriteDataToFS(deps.CurrentFS, propertiesData); err != nil {
		return err
	}
	return nil
}
