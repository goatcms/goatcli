package buildc

import (
	"github.com/goatcms/goatcli/cliapp/common/prevents"
	"github.com/goatcms/goatcli/cliapp/common/result"
	"github.com/goatcms/goatcli/cliapp/gclicommands/vcsc"
	"github.com/goatcms/goatcli/cliapp/gcliservices"
	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/varutil/goaterr"
)

// RunBuild run build command
func RunBuild(a app.App, ctx app.IOContext) (err error) {
	var (
		deps struct {
			Interactive    string                      `argument:"?interactive" ,command:"?interactive"`
			BuilderService gcliservices.BuilderService `dependency:"BuilderService"`
			ClonerService  gcliservices.ClonerService  `dependency:"ClonerService"`
			GCLIInputs     gcliservices.GCLIInputs     `dependency:"GCLIInputs"`
		}
		propertiesData map[string]string
		secretsData    map[string]string
		appData        gcliservices.ApplicationData
		fs             = ctx.IO().CWD()
	)
	if err = vcsc.RunScan(a, ctx); err != nil {
		return nil
	}
	if err = goaterr.ToErrors(goaterr.AppendError(nil,
		a.DependencyProvider().InjectTo(&deps),
		ctx.Scope().InjectTo(&deps))); err != nil {
		return err
	}
	if err = prevents.RequireGoatProject(fs); err != nil {
		return err
	}
	if propertiesData, secretsData, appData, err = deps.GCLIInputs.Inputs(ctx); err != nil {
		return err
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
