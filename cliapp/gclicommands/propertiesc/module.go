package propertiesc

import (
	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/app/terminal"
)

func GetCommand() app.TerminalCommand {
	return terminal.NewCommand(terminal.CommandParams{
		Name:          "properties:get",
		Help:          "display property by name",
		Callback:      RunSetPropertyValue,
		MainArguments: []string{"name"},
		Arguments: terminal.NewArguments(
			terminal.NewArgument(terminal.ArgumentParams{
				Name:     "name",
				Help:     "the property name",
				Required: true,
				Type:     app.TerminalTextArgument,
			}),
		),
	})
}

func SetCommand() app.TerminalCommand {
	return terminal.NewCommand(terminal.CommandParams{
		Name:          "properties:set",
		Help:          "adds or updates a property with a specified key and value",
		Callback:      RunSetPropertyValue,
		MainArguments: []string{"name", "value"},
		Arguments: terminal.NewArguments(
			terminal.NewArgument(terminal.ArgumentParams{
				Name:     "name",
				Help:     "the property name",
				Required: true,
				Type:     app.TerminalTextArgument,
			}),
			terminal.NewArgument(terminal.ArgumentParams{
				Name:     "value",
				Help:     "the property value",
				Required: true,
				Type:     app.TerminalTextArgument,
			}),
		),
	})
}

func Commands() []app.TerminalCommand {
	return []app.TerminalCommand{
		GetCommand(),
		SetCommand(),
	}
}
