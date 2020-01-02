package vcsc

import (
	"time"

	"github.com/goatcms/goatcli/cliapp/gcliservices"
	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/filesystem"
	"github.com/goatcms/goatcore/varutil/goaterr"
)

// RunGeneratedList run vcs:ignored:list command
func RunGeneratedList(a app.App, ctx app.IOContext) (err error) {
	var (
		deps struct {
			CurrentFS  filesystem.Filespace    `filespace:"current"`
			VCSService gcliservices.VCSService `dependency:"VCSService"`
		}
		vcsData gcliservices.VCSData
	)
	if err = goaterr.ToErrors(goaterr.AppendError(nil,
		a.DependencyProvider().InjectTo(&deps),
		ctx.Scope().InjectTo(&deps))); err != nil {
		return err
	}
	if vcsData, err = deps.VCSService.ReadDataFromFS(deps.CurrentFS); err != nil {
		return err
	}
	for _, row := range vcsData.VCSGeneratedFiles().All() {
		ctx.IO().Out().Printf(" %s: %s\n", row.ModTime.Format(time.RFC3339), row.Path)
	}
	return nil
}
