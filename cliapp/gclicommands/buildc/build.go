package buildc

import (
	"strings"

	"github.com/goatcms/goatcli/cliapp/common/am"
	"github.com/goatcms/goatcli/cliapp/common/config"
	"github.com/goatcms/goatcli/cliapp/common/prevents"
	"github.com/goatcms/goatcli/cliapp/common/result"
	"github.com/goatcms/goatcli/cliapp/gclicommands/vcsc"
	"github.com/goatcms/goatcli/cliapp/gcliservices"
	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/filesystem"
	"github.com/goatcms/goatcore/varutil/goaterr"
)

// RunBuild run build command
func RunBuild(a app.App, ctx app.IOContext) (err error) {
	var (
		deps struct {
			Interactive       string                         `argument:"?interactive" ,command:"?interactive"`
			CurrentFS         filesystem.Filespace           `filespace:"current"`
			VCSService        gcliservices.VCSService        `dependency:"VCSService"`
			PropertiesService gcliservices.PropertiesService `dependency:"PropertiesService"`
			SecretsService    gcliservices.SecretsService    `dependency:"SecretsService"`
			BuilderService    gcliservices.BuilderService    `dependency:"BuilderService"`
			ClonerService     gcliservices.ClonerService     `dependency:"ClonerService"`
			DataService       gcliservices.DataService       `dependency:"DataService"`
		}
		propertiesDef  []*config.Property
		propertiesData map[string]string
		secretsDef     []*config.Property
		secretsData    map[string]string
		isChanged      bool
		data           map[string]string
		interactive    bool
		fs             filesystem.Filespace
		appData        gcliservices.ApplicationData
	)
	if err = vcsc.RunScan(a, ctx); err != nil {
		return nil
	}
	if err = goaterr.ToErrors(goaterr.AppendError(nil,
		a.DependencyProvider().InjectTo(&deps),
		ctx.Scope().InjectTo(&deps))); err != nil {
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
	if isChanged, err = deps.PropertiesService.FillData(ctx, propertiesDef, propertiesData, map[string]string{}, interactive); err != nil {
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
	if isChanged, err = deps.SecretsService.FillData(ctx, secretsDef, secretsData, map[string]string{}, interactive); err != nil {
		return err
	}
	if isChanged {
		if err = deps.SecretsService.WriteDataToFS(fs, secretsData); err != nil {
			return err
		}
	}
	// Clone modules (if required)
	ctx.IO().Out().Printf("start clone modules... ")
	propertiesResult := result.NewPropertiesResult(propertiesData)
	if err = deps.ClonerService.CloneModules(fs, fs, propertiesResult); err != nil {
		return err
	}
	ctx.IO().Out().Printf("cloned\n")
	// Build
	ctx.IO().Out().Printf("start build... ")
	if err = deps.BuilderService.Build(ctx, fs, appData, propertiesData, secretsData); err != nil {
		return err
	}
	if err = ctx.Scope().Wait(); err != nil {
		return goaterr.ToErrors(goaterr.AppendError(nil,
			err,
			ctx.Scope().Trigger(app.RollbackEvent, nil)))
	}
	ctx.IO().Out().Printf("builded\n")
	ctx.IO().Out().Printf("start commit... ")
	if err = ctx.Scope().Trigger(app.CommitEvent, nil); err != nil {
		return err
	}
	return ctx.IO().Out().Printf("commited\n")
}
