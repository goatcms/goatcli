package scriptsc

import (
	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/app/terminal"
)

func EnvsCommand() app.TerminalCommand {
	return terminal.NewCommand(terminal.CommandParams{
		Name:     "scripts:envs",
		Help:     "list scripts environments",
		Callback: RunScript,
	})
}
func RunCommand() app.TerminalCommand {
	return terminal.NewCommand(terminal.CommandParams{
		Name:          "scripts:run",
		Help:          "run script by name",
		Callback:      RunScript,
		MainArguments: []string{"name"},
		Arguments: terminal.NewArguments(
			terminal.NewArgument(terminal.ArgumentParams{
				Name:     "name",
				Help:     "the script name",
				Required: true,
				Type:     app.TerminalTextArgument,
			}),
		),
	})
}

func Commands() []app.TerminalCommand {
	return []app.TerminalCommand{
		EnvsCommand(),
		RunCommand(),
	}
}
