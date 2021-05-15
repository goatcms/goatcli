package secretsc

import (
	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/app/terminal"
)

func GetCommand() app.TerminalCommand {
	return terminal.NewCommand(terminal.CommandParams{
		Callback:      RunGetSecretValue,
		Help:          "display secret by name",
		Name:          "secrets:get",
		MainArguments: []string{"name"},
		Arguments: terminal.NewArguments(
			terminal.NewArgument(terminal.ArgumentParams{
				Name:     "name",
				Help:     "the secret name",
				Required: true,
				Type:     app.TerminalTextArgument,
			}),
		),
	})
}

func SetCommand() app.TerminalCommand {
	return terminal.NewCommand(terminal.CommandParams{
		Callback:      RunSetSecretValue,
		Help:          "adds or updates an secret with specified key and value",
		Name:          "secrets:set",
		MainArguments: []string{"name", "value"},
		Arguments: terminal.NewArguments(
			terminal.NewArgument(terminal.ArgumentParams{
				Name:     "name",
				Help:     "the secret name",
				Required: true,
				Type:     app.TerminalTextArgument,
			}),
			terminal.NewArgument(terminal.ArgumentParams{
				Name:     "value",
				Help:     "the secret value",
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
