package vcsc

import (
	"github.com/goatcms/goatcli/cliapp/services"
	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/filesystem"
	"github.com/goatcms/goatcore/varutil/goaterr"
)

// RunIgnoredList run vcs:ignored:list command
func RunIgnoredList(a app.App, ctxScope app.Scope) (err error) {
	var (
		deps struct {
			CurrentFS filesystem.Filespace `filespace:"current"`

			VCSService services.VCSService `dependency:"VCSService"`
			Input      app.Input           `dependency:"InputService"`
			Output     app.Output          `dependency:"OutputService"`
		}
		vcsData services.VCSData
	)
	if err = goaterr.ToErrors(goaterr.AppendError(nil,
		a.DependencyProvider().InjectTo(&deps),
		ctxScope.InjectTo(&deps))); err != nil {
		return err
	}
	if vcsData, err = deps.VCSService.ReadDataFromFS(deps.CurrentFS); err != nil {
		return err
	}
	for _, row := range vcsData.VCSIgnoredFiles().All() {
		deps.Output.Printf(" %s\n", row)
	}
	return nil
}
