package vcsc

import (
	"github.com/goatcms/goatcli/cliapp/gcliservices"
	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/filesystem"
	"github.com/goatcms/goatcore/varutil/goaterr"
)

// RunPersistedAdd run vcs:persisted:add command
func RunPersistedAdd(a app.App, ctx app.IOContext) (err error) {
	var (
		deps struct {
			Path       string                  `argument:"?$2" ,command:"?$2"`
			CurrentFS  filesystem.Filespace    `filespace:"current"`
			VCSService gcliservices.VCSService `dependency:"VCSService"`
		}
		vcsData gcliservices.VCSData
	)
	if err = goaterr.ToError(goaterr.AppendError(nil,
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
	vcsData.VCSPersistedFiles().AddPath(deps.Path)
	if err = deps.VCSService.WriteDataToFS(deps.CurrentFS, vcsData); err != nil {
		return err
	}
	return ctx.IO().Out().Printf("Path added\n")
}
