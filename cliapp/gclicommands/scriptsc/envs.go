package scriptsc

import (
	"sort"

	"github.com/goatcms/goatcli/cliapp/gcliservices"
	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/app/modules/commonm/commservices"
)

// RunScriptsEnvs run scripts:envs command
func RunScriptsEnvs(a app.App, ctx app.IOContext) (err error) {
	var (
		deps struct {
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
	allEnvs := envs.All()
	keys := make([]string, 0, len(allEnvs))
	for key := range allEnvs {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		ctx.IO().Out().Printf("%s\n", key)
	}
	return nil
}
