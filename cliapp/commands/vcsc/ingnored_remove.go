package vcsc

import (
	"github.com/goatcms/goatcli/cliapp/services"
	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/filesystem"
	"github.com/goatcms/goatcore/varutil/goaterr"
)

// RunIgnoredRemove run vcs:ignored:remove command
func RunIgnoredRemove(a app.App, ctx app.IOContext) (err error) {
	var (
		deps struct {
			Path string `argument:"?$2" ,command:"?$2"`

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
	if deps.Path == "" {
		return goaterr.NewError("Path is required")
	}
	if vcsData, err = deps.VCSService.ReadDataFromFS(deps.CurrentFS); err != nil {
		return err
	}
	vcsData.VCSIgnoredFiles().RemovePath(deps.Path)
	if err = deps.VCSService.WriteDataToFS(deps.CurrentFS, vcsData); err != nil {
		return err
	}
	return ctx.IO().Out().Printf("Path removed\n")
}
