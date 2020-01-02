package cleanc

import (
	"github.com/goatcms/goatcli/cliapp/common/config"
	"github.com/goatcms/goatcli/cliapp/services"
	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/filesystem"
	"github.com/goatcms/goatcore/varutil/goaterr"
)

// RunCleanDependencies run clean:dependencies command
func RunCleanDependencies(a app.App, ctx app.IOContext) (err error) {
	var (
		deps struct {
			CurrentFS    filesystem.Filespace         `filespace:"current"`
			Dependencies services.DependenciesService `dependency:"DependenciesService"`
		}
		configDeps []*config.Dependency
	)
	if err = goaterr.ToErrors(goaterr.AppendError(nil,
		a.DependencyProvider().InjectTo(&deps),
		ctx.Scope().InjectTo(&deps))); err != nil {
		return err
	}
	// load data
	if configDeps, err = deps.Dependencies.ReadDefFromFS(deps.CurrentFS); err != nil {
		return err
	}
	for _, row := range configDeps {
		if err = deps.CurrentFS.RemoveAll(row.Dest); err != nil {
			return err
		}
	}
	return ctx.IO().Out().Printf("dependencies cleaned\n")
}
