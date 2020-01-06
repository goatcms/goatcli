package scriptsc

import (
	"github.com/goatcms/goatcli/cliapp/gcliservices"
	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/varutil/goaterr"
)

// RunScript run script by name
func RunScript(a app.App, ctx app.IOContext) (err error) {
	var (
		deps struct {
			Name          string                     `command:"?$1"`
			ScriptsRunner gcliservices.ScriptsRunner `dependency:"ScriptsRunner"`
			GCLIInputs    gcliservices.GCLIInputs    `dependency:"GCLIInputs"`
		}
		propertiesData map[string]string
		secretsData    map[string]string
		appData        gcliservices.ApplicationData
		ctxScope       = ctx.Scope()
	)
	if err = a.DependencyProvider().InjectTo(&deps); err != nil {
		return err
	}
	if err = ctxScope.InjectTo(&deps); err != nil {
		return err
	}
	if propertiesData, secretsData, appData, err = deps.GCLIInputs.Inputs(ctx); err != nil {
		return err
	}
	if deps.Name == "" {
		return goaterr.Errorf(FirstKeyParameterIsRequired)
	}
	return deps.ScriptsRunner.Run(ctx, ctx.IO().CWD(), deps.Name, propertiesData, secretsData, appData)
}
