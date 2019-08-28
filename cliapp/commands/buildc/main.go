package buildc

import (
	"github.com/goatcms/goatcli/cliapp/common/config"
	"github.com/goatcms/goatcli/cliapp/common/prevents"
	"github.com/goatcms/goatcli/cliapp/common/result"
	"github.com/goatcms/goatcli/cliapp/services"
	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/filesystem"
	"github.com/goatcms/goatcore/filesystem/fscache"
	"github.com/goatcms/goatcore/varutil/goaterr"
)

// Run run command in app.App context
func Run(a app.App, ctxScope app.Scope) (err error) {
	var (
		deps struct {
			Interactive string `argument:"?interactive",command:"?interactive"`

			CurrentFS filesystem.Filespace `filespace:"current"`

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
		builderDef     []*config.Build
		data           map[string]string
		interactive    bool
		fs             fscache.Cache
	)
	if err = a.DependencyProvider().InjectTo(&deps); err != nil {
		return err
	}
	if err = ctxScope.InjectTo(&deps); err != nil {
		return err
	}
	interactive = !(deps.Interactive == "false")
	if fs, err = fscache.NewMemCache(fs); err != nil {
		return err
	}
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
	// load secrets
	if secretsDef, err = deps.SecretsService.ReadDefFromFS(fs); err != nil {
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
	// load data
	if data, err = deps.DataService.ReadDataFromFS(fs); err != nil {
		return err
	}
	// Clone modules (if required)
	propertiesResult := result.NewPropertiesResult(propertiesData)
	if err = deps.ClonerService.CloneModules(fs, fs, propertiesResult); err != nil {
		return err
	}
	// build
	if builderDef, err = deps.BuilderService.ReadDefFromFS(fs); err != nil {
		return err
	}
	if err = deps.BuilderService.Build(ctxScope, fs, builderDef, data, propertiesData, secretsData); err != nil {
		return err
	}
	if err = ctxScope.Wait(); err != nil {
		return goaterr.ToErrors(goaterr.AppendError(nil,
			err,
			ctxScope.Trigger(app.RollbackEvent, nil)))
	}
	deps.Output.Printf("builded\n")
	if err = fs.Commit(); err != nil {
		return err
	}
	deps.Output.Printf("commit to filesystem...\n")
	if err = ctxScope.Trigger(app.CommitEvent, nil); err != nil {
		return err
	}
	deps.Output.Printf("start after commit tasks...\n")
	if err = ctxScope.Trigger(app.AfterCommitEvent, nil); err != nil {
		return err
	}
	deps.Output.Printf("after commit tasks done ")
	return nil
}
