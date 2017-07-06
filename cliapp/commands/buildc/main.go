package buildc

import (
	"github.com/goatcms/goatcli/cliapp/common/config"
	"github.com/goatcms/goatcli/cliapp/common/prevents"
	"github.com/goatcms/goatcli/cliapp/services"
	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/filesystem"
)

// Run run command in app.App context
func Run(a app.App) (err error) {
	var (
		deps struct {
			CurrentFS filesystem.Filespace `filespace:"current"`

			PropertiesService services.PropertiesService `dependency:"PropertiesService"`
			BuilderService    services.BuilderService    `dependency:"BuilderService"`
			DataService       services.DataService       `dependency:"DataService"`
			Input             app.Input                  `dependency:"InputService"`
			Output            app.Output                 `dependency:"OutputService"`
		}
		propertiesDef  []*config.Property
		propertiesData map[string]string
		isChanged      bool
		builderDef     []*config.Build
		data           map[string]string
	)
	if err = a.DependencyProvider().InjectTo(&deps); err != nil {
		return err
	}
	if err = prevents.RequireGoatProject(deps.CurrentFS); err != nil {
		return err
	}
	// load properties
	if propertiesData, err = deps.PropertiesService.ReadDataFromFS(deps.CurrentFS); err != nil {
		return err
	}
	if isChanged, err = deps.PropertiesService.FillData(propertiesDef, propertiesData, map[string]string{}); err != nil {
		return err
	}
	if isChanged {
		if err = deps.PropertiesService.WriteDataToFS(deps.CurrentFS, propertiesData); err != nil {
			return err
		}
	}
	// load data
	if data, err = deps.DataService.ReadDataFromFS(deps.CurrentFS); err != nil {
		return err
	}
	// build
	if builderDef, err = deps.BuilderService.ReadDefFromFS(deps.CurrentFS); err != nil {
		return err
	}
	if err = deps.BuilderService.Build(deps.CurrentFS, builderDef, data, propertiesData); err != nil {
		return err
	}
	deps.Output.Printf("builded")
	return nil
}
