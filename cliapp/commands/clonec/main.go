package clonec

import (
	"github.com/goatcms/goatcli/cliapp/common/config"
	"github.com/goatcms/goatcli/cliapp/common/result"
	"github.com/goatcms/goatcli/cliapp/services"
	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/filesystem"
)

// Run run command in app.App context
func Run(a app.App) (err error) {
	var (
		deps struct {
			Command       string `argument:"$1"`
			RepositoryURL string `argument:"?$2"`
			RepositoryRev string `argument:"?rev"`
			DestPath      string `argument:"?$3"`

			RootFilespace filesystem.Filespace `filespace:"root"`

			RepositoriesService services.RepositoriesService `dependency:"RepositoriesService"`
			PropertiesService   services.PropertiesService   `dependency:"PropertiesService"`
			CloneService        services.ClonerService       `dependency:"ClonerService"`
			Input               app.Input                    `dependency:"InputService"`
			Output              app.Output                   `dependency:"OutputService"`
		}
		repofs         filesystem.Filespace
		propertiesDef  []*config.Property
		propertiesData map[string]string
		isChanged      bool
		destfs         filesystem.Filespace
	)
	if err = a.DependencyProvider().InjectTo(&deps); err != nil {
		return err
	}
	if deps.RepositoryURL == "" {
		deps.Output.Printf("Unknown url to clone")
		return nil
	}
	if deps.DestPath == "" {
		deps.Output.Printf("Unknown destination path")
		return nil
	}
	if repofs, err = deps.RepositoriesService.Filespace(deps.RepositoryURL, deps.RepositoryRev); err != nil {
		deps.Output.Printf("%s", err)
		return nil
	}
	if propertiesDef, err = deps.PropertiesService.ReadDefFromFS(repofs); err != nil {
		deps.Output.Printf("%s", err)
		return nil
	}
	if propertiesData, err = deps.PropertiesService.ReadDataFromFS(repofs); err != nil {
		deps.Output.Printf("%s", err)
		return nil
	}
	if isChanged, err = deps.PropertiesService.FillData(propertiesDef, propertiesData, map[string]string{}); err != nil {
		deps.Output.Printf("%s", err)
		return nil
	}
	if err = deps.RootFilespace.MkdirAll(deps.DestPath, 0766); err != nil {
		return err
	}
	if destfs, err = deps.RootFilespace.Filespace(deps.DestPath); err != nil {
		return err
	}
	if isChanged {
		if err = deps.PropertiesService.WriteDataToFS(destfs, propertiesData); err != nil {
			deps.Output.Printf("%s", err)
			return nil
		}
	}
	propertiesResult := result.NewPropertiesResult(propertiesData)
	if err = deps.CloneService.Clone(deps.RepositoryURL, deps.RepositoryRev, destfs, propertiesResult); err != nil {
		deps.Output.Printf("%s", err)
		return nil
	}
	deps.Output.Printf("cloned")
	return nil
}
