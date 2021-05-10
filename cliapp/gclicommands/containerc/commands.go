package containerc

import (
	"github.com/goatcms/goatcli/cliapp/gclicommands/containerc/imagepip"
	"github.com/goatcms/goatcore/app/modules"
	"github.com/goatcms/goatcore/app/modules/pipelinem/pipcommands"
	"github.com/goatcms/goatcore/app/modules/pipelinem/pipcommands/pipc"
	"github.com/goatcms/goatcore/app/modules/terminalm"
	"github.com/goatcms/goatcore/app/scope"

	"github.com/goatcms/goatcore/app"
)

var (
	commandScope = newCommandScope()
)

func newCommandScope() (scp app.Scope) {
	scp = scope.NewScope(scope.Params{})
	app.RegisterScopeCommand(scp, "pip:run", pipc.Run, pipcommands.PipRun)
	app.RegisterScopeCommand(scp, "pip:try", pipc.Try, pipcommands.PipTry)
	app.RegisterScopeCommand(scp, "build", imagepip.RunBuild, "")
	app.RegisterScopeCommand(scp, "push", imagepip.RunPush, "")
	return scp
}

func newTerminal(a app.App) (term modules.Terminal, err error) {
	return terminalm.NewIOTerminal(a, commandScope)
}
