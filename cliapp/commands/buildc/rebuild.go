package buildc

import (
	"github.com/goatcms/goatcli/cliapp/commands/cleanc"
	"github.com/goatcms/goatcore/app"
)

// RunRebuild run rebuild command
func RunRebuild(a app.App, ctxScope app.Scope) (err error) {
	if err = cleanc.RunCleanBuild(a, ctxScope); err != nil {
		return nil
	}
	return RunBuild(a, ctxScope)
}
