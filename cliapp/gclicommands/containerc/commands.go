package containerc

import (
	"github.com/goatcms/goatcli/cliapp/gclicommands/containerc/imagepip"
	"github.com/goatcms/goatcore/app/modules/pipelinem/pipcommands/pipc"
	"github.com/goatcms/goatcore/app/modules/terminalm/termservices"
	"github.com/goatcms/goatcore/app/modules/terminalm/termservices/terminals"
	"github.com/goatcms/goatcore/app/terminal"

	"github.com/goatcms/goatcore/app"
)

var (
	imagePipCommands = []app.TerminalCommand{
		pipc.RunCommand(),
		pipc.TryCommand(),
		imagepip.BuildCommand(),
		imagepip.PushCommand(),
	}
)

func newTerminal(a app.App) (term termservices.Terminal, err error) {
	return terminals.NewIOTerminal(a, terminal.NewCommands(imagePipCommands...)), nil
}
