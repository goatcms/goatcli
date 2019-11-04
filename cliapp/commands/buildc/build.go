package buildc

import (
	"strings"

	"github.com/goatcms/goatcli/cliapp/common/am"
	"github.com/goatcms/goatcli/cliapp/common/config"
	"github.com/goatcms/goatcli/cliapp/common/prevents"
	"github.com/goatcms/goatcli/cliapp/common/result"
	"github.com/goatcms/goatcli/cliapp/services"
	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/filesystem"
	"github.com/goatcms/goatcore/varutil/goaterr"
)

// RunBuild run build command
func RunBuild(a app.App, ctxScope app.Scope) (err error) {
	var (
		deps struct {
			Interactive string `argument:"?interactive" ,command:"?interactive"`

			CurrentFS filesystem.Filespace `filespace:"current"`

			VCSService        services.VCSService        `dependency:"VCSService"`
			PropertiesService services.PropertiesService `dependency:"PropertiesService"`
			SecretsService    services.SecretsService    `dependency:"SecretsService"`
			BuilderService    services.BuilderService    `dependency:"BuilderService"`
			ClonerService     services.ClonerService     `dependency:"ClonerService"`
			DataService       services.DataService       `dependency:"DataService"`
			Input             app.Input                  `dependency:"InputService"`
			Output            app.Output                 `dependency:"OutputService"`
		}
		propertiesDef  []*config.Property
		propertiesData map[string]string
		secretsDef     []*config.Property
		secretsData    map[string]string
		isChanged      bool
		data           map[string]string
		interactive    bool
		fs             filesystem.Filespace
		appData        services.ApplicationData
	)
	if err = goaterr.ToErrors(goaterr.AppendError(nil,
		a.DependencyProvider().InjectTo(&deps),
		ctxScope.InjectTo(&deps))); err != nil {
		return err
	}
	interactive = strings.ToLower(deps.Interactive) != "false"
	fs = deps.CurrentFS
	if err = prevents.RequireGoatProject(fs); err != nil {
		return err
	}
	// load properties
	if propertiesDef, err = deps.PropertiesService.ReadDefFromFS(fs); err != nil {
		return err
	}
	if propertiesData, err = deps.PropertiesService.ReadDataFromFS(fs); err != nil {
		return err
	}
	if isChanged, err = deps.PropertiesService.FillData(propertiesDef, propertiesData, map[string]string{}, interactive); err != nil {
		return err
	}
	if isChanged {
		if err = deps.PropertiesService.WriteDataToFS(fs, propertiesData); err != nil {
			return err
		}
	}
	// load data
	if data, err = deps.DataService.ReadDataFromFS(fs); err != nil {
		return err
	}
	appData = am.NewApplicationData(data)
	// load secrets
	if secretsDef, err = deps.SecretsService.ReadDefFromFS(fs, propertiesData, appData); err != nil {
		return err
	}
	if secretsData, err = deps.SecretsService.ReadDataFromFS(fs); err != nil {
		return err
	}
	if isChanged, err = deps.SecretsService.FillData(secretsDef, secretsData, map[string]string{}, interactive); err != nil {
		return err
	}
	if isChanged {
		if err = deps.SecretsService.WriteDataToFS(fs, secretsData); err != nil {
			return err
		}
	}
	// Clone modules (if required)
	deps.Output.Printf("start clone modules... ")
	propertiesResult := result.NewPropertiesResult(propertiesData)
	if err = deps.ClonerService.CloneModules(fs, fs, propertiesResult); err != nil {
		return err
	}
	deps.Output.Printf("cloned\n")
	// Build
	deps.Output.Printf("start build... ")
	buildContext := deps.BuilderService.NewContext(ctxScope, appData, propertiesData, secretsData)
	if err = buildContext.Build(fs); err != nil {
		return err
	}
	if err = ctxScope.Wait(); err != nil {
		return goaterr.ToErrors(goaterr.AppendError(nil,
			err,
			ctxScope.Trigger(app.RollbackEvent, nil)))
	}
	deps.Output.Printf("builded\n")
	deps.Output.Printf("start commit... ")
	if err = ctxScope.Trigger(app.CommitEvent, nil); err != nil {
		return err
	}
	deps.Output.Printf("commited\n")
	return nil
}
