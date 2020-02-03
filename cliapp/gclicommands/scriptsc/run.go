package scriptsc

import (
	"strings"

	"github.com/goatcms/goatcli/cliapp/common/naming"
	"github.com/goatcms/goatcli/cliapp/gcliservices"
	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/app/modules/commonm/commservices"
	"github.com/goatcms/goatcore/varutil/goaterr"
)

// RunScript run script by name
func RunScript(a app.App, ctx app.IOContext) (err error) {
	var (
		deps struct {
			Name             string                        `command:"?$1"`
			ScriptsRunner    gcliservices.ScriptsRunner    `dependency:"ScriptsRunner"`
			GCLIInputs       gcliservices.GCLIInputs       `dependency:"GCLIInputs"`
			EnvironmentsUnit commservices.EnvironmentsUnit `dependency:"CommonEnvironmentsUnit"`
		}
		propertiesData map[string]string
		secretsData    map[string]string
		appData        gcliservices.ApplicationData
		envs           commservices.Environments
		ctxScope       = ctx.Scope()
	)
	if err = a.DependencyProvider().InjectTo(&deps); err != nil {
		return err
	}
	if err = ctxScope.InjectTo(&deps); err != nil {
		return err
	}
	if deps.Name == "" {
		return goaterr.Errorf(FirstKeyParameterIsRequired)
	}
	// load variables
	if propertiesData, secretsData, appData, err = deps.GCLIInputs.Inputs(ctx); err != nil {
		return err
	}
	if envs, err = deps.EnvironmentsUnit.Envs(ctxScope); err != nil {
		return err
	}
	for key, value := range secretsData {
		key = "SECRET_" + strings.ToUpper(naming.ToUnderscore(key))
		if err = envs.Set(key, value); err != nil {
			return err
		}
	}
	// run script
	return deps.ScriptsRunner.Run(ctx, ctx.IO().CWD(), deps.Name, propertiesData, secretsData, appData)
}
