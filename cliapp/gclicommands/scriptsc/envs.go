package scriptsc

import (
	"github.com/goatcms/goatcli/cliapp/gcliservices"
	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/app/modules/commonm/commservices"
)

// RunScriptsEnvs run scripts:envs command
func RunScriptsEnvs(a app.App, ctx app.IOContext) (err error) {
	var (
		deps struct {
			Name             string                        `command:"?$1"`
			ScriptsRunner    gcliservices.ScriptsRunner    `dependency:"ScriptsRunner"`
			GCLIInputs       gcliservices.GCLIInputs       `dependency:"GCLIInputs"`
			GCLIEnvironment  gcliservices.GCLIEnvironment  `dependency:"GCLIEnvironment"`
			EnvironmentsUnit commservices.EnvironmentsUnit `dependency:"CommonEnvironmentsUnit"`
		}
		propertiesData map[string]string
		secretsData    map[string]string
		envs           commservices.Environments
		ctxScope       = ctx.Scope()
	)
	if err = a.DependencyProvider().InjectTo(&deps); err != nil {
		return err
	}
	if err = ctxScope.InjectTo(&deps); err != nil {
		return err
	}
	// load variables and envs
	if propertiesData, secretsData, _, err = deps.GCLIInputs.Inputs(ctx); err != nil {
		return err
	}
	if err = deps.GCLIEnvironment.LoadEnvs(ctxScope, propertiesData, secretsData); err != nil {
		return err
	}
	// list envs
	if envs, err = deps.EnvironmentsUnit.Envs(ctxScope); err != nil {
		return err
	}
	for name := range envs.All() {
		ctx.IO().Out().Printf("%s\n", name)
	}
	return nil
}
