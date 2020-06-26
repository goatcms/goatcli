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
			Name               string                          `command:"?$1"`
			ScriptsRunner      gcliservices.ScriptsRunner      `dependency:"ScriptsRunner"`
			GCLIProjectManager gcliservices.GCLIProjectManager `dependency:"GCLIProjectManager"`
			GCLIEnvironment    gcliservices.GCLIEnvironment    `dependency:"GCLIEnvironment"`
			EnvironmentsUnit   commservices.EnvironmentsUnit   `dependency:"CommonEnvironmentsUnit"`
		}
		project  *gcliservices.Project
		envs     commservices.Environments
		ctxScope = ctx.Scope()
	)
	if err = a.DependencyProvider().InjectTo(&deps); err != nil {
		return err
	}
	if err = ctxScope.InjectTo(&deps); err != nil {
		return err
	}
	// load variables and envs
	if project, err = deps.GCLIProjectManager.Project(ctx); err != nil {
		return err
	}
	if err = deps.GCLIEnvironment.LoadEnvs(ctxScope, project.Properties, project.Secrets); err != nil {
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
