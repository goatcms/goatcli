package initc

import (
	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/app/terminal"
)

func InitCommands() app.TerminalCommand {
	return terminal.NewCommand(terminal.CommandParams{
		Name:     "init",
		Help:     "initialize new goat project",
		Callback: RunInit,
	})
}

func Commands() []app.TerminalCommand {
	return []app.TerminalCommand{
		InitCommands(),
	}
}
