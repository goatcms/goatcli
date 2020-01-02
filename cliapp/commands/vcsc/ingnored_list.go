package vcsc

import (
	"github.com/goatcms/goatcli/cliapp/services"
	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/filesystem"
	"github.com/goatcms/goatcore/varutil/goaterr"
)

// RunIgnoredList run vcs:ignored:list command
func RunIgnoredList(a app.App, ctx app.IOContext) (err error) {
	var (
		deps struct {
			CurrentFS  filesystem.Filespace `filespace:"current"`
			VCSService services.VCSService  `dependency:"VCSService"`
		}
		vcsData services.VCSData
	)
	if err = goaterr.ToErrors(goaterr.AppendError(nil,
		a.DependencyProvider().InjectTo(&deps),
		ctx.Scope().InjectTo(&deps))); err != nil {
		return err
	}
	if vcsData, err = deps.VCSService.ReadDataFromFS(deps.CurrentFS); err != nil {
		return err
	}
	for _, row := range vcsData.VCSIgnoredFiles().All() {
		ctx.IO().Out().Printf(" %s\n", row)
	}
	return nil
}
