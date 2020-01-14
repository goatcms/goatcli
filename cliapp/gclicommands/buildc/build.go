package buildc

import (
	"github.com/goatcms/goatcli/cliapp/common/prevents"
	"github.com/goatcms/goatcli/cliapp/common/result"
	"github.com/goatcms/goatcli/cliapp/gclicommands/vcsc"
	"github.com/goatcms/goatcli/cliapp/gcliservices"
	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/app/gio"
	"github.com/goatcms/goatcore/filesystem"
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
		scope          app.Scope
		childCtx       app.IOContext
		out            app.Output
		fs             filesystem.Filespace
	)
	childCtx = gio.NewChildIOContext(ctx, gio.IOContextParams{})
	defer childCtx.Scope().Close()
	scope = childCtx.Scope()
	out = childCtx.IO().Out()
	fs = childCtx.IO().CWD()
	if err = vcsc.RunScan(a, ctx); err != nil {
		return err
	}
	if err = goaterr.ToErrors(goaterr.AppendError(nil,
		a.DependencyProvider().InjectTo(&deps),
		scope.InjectTo(&deps))); err != nil {
		return err
	}
	if err = prevents.RequireGoatProject(fs); err != nil {
		return err
	}
	if propertiesData, secretsData, appData, err = deps.GCLIInputs.Inputs(ctx); err != nil {
		return err
	}
	// Clone modules (if required)
	out.Printf("start clone modules... ")
	propertiesResult := result.NewPropertiesResult(propertiesData)
	if err = deps.ClonerService.CloneModules(fs, fs, propertiesResult); err != nil {
		return err
	}
	out.Printf("cloned\n")
	// Build
	out.Printf("start build... ")
	if err = deps.BuilderService.Build(ctx, fs, appData, propertiesData, secretsData); err != nil {
		return err
	}
	if err = scope.Wait(); err != nil {
		return goaterr.ToErrors(goaterr.AppendError(nil,
			err,
			scope.Trigger(app.RollbackEvent, nil)))
	}
	out.Printf("builded\n")
	out.Printf("start commit... ")
	if err = scope.Trigger(app.CommitEvent, nil); err != nil {
		return err
	}
	return out.Printf("commited\n")
}
