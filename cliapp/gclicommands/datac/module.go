package datac

import (
	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/app/terminal"
)

func AddCommand() app.TerminalCommand {
	return terminal.NewCommand(terminal.CommandParams{
		Name:          "data:add",
		Help:          "add new data to project",
		Callback:      RunAdd,
		MainArguments: []string{"type", "data"},
		Arguments: terminal.NewArguments(
			terminal.NewArgument(terminal.ArgumentParams{
				Name:     "type",
				Help:     "the name of data type to add",
				Required: true,
				Type:     app.TerminalTextArgument,
			}),
			terminal.NewArgument(terminal.ArgumentParams{
				Name:     "data",
				Help:     "the name of data set to add",
				Required: true,
				Type:     app.TerminalTextArgument,
			}),
		),
	})
}

func Commands() []app.TerminalCommand {
	return []app.TerminalCommand{
		AddCommand(),
	}
}
