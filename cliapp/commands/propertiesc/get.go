package propertiesc

import (
	"fmt"

	"github.com/goatcms/goatcli/cliapp/services"
	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/filesystem"
)

// RunGetPropertyValue run command to set property value
func RunGetPropertyValue(a app.App) (err error) {
	var (
		deps struct {
			Key               string                     `argument:"?$2"`
			CurrentFS         filesystem.Filespace       `filespace:"current"`
			PropertiesService services.PropertiesService `dependency:"PropertiesService"`
			Input             app.Input                  `dependency:"InputService"`
			Output            app.Output                 `dependency:"OutputService"`
		}
		propertiesData map[string]string
		ok             bool
		value          string
	)
	if err = a.DependencyProvider().InjectTo(&deps); err != nil {
		return err
	}
	if deps.Key == "" {
		deps.Output.Printf(FirstKeyParameterIsRequired)
		return fmt.Errorf(FirstKeyParameterIsRequired)
	}
	if propertiesData, err = deps.PropertiesService.ReadDataFromFS(deps.CurrentFS); err != nil {
		return err
	}
	if value, ok = propertiesData[deps.Key]; !ok {
		return fmt.Errorf("Unknow Value for %s key", deps.Key)
	}
	return deps.Output.Printf(value)
}
