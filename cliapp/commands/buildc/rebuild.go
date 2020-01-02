package buildc

import (
	"github.com/goatcms/goatcli/cliapp/commands/cleanc"
	"github.com/goatcms/goatcore/app"
)

// RunRebuild run rebuild command
func RunRebuild(a app.App, ctx app.IOContext) (err error) {
	if err = cleanc.RunCleanBuild(a, ctx); err != nil {
		return nil
	}
	return RunBuild(a, ctx)
}
