package buildc

import (
	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/app/terminal"
)

func BuildCommand() app.TerminalCommand {
	return terminal.NewCommand(terminal.CommandParams{
		Name:     "build",
		Help:     "build goat project in current directory",
		Callback: RunBuild,
	})
}

func RebuildCommand() app.TerminalCommand {
	return terminal.NewCommand(terminal.CommandParams{
		Name:     "rebuild",
		Help:     "clean and build project",
		Callback: RunRebuild,
	})
}

func Commands() []app.TerminalCommand {
	return []app.TerminalCommand{
		BuildCommand(),
		RebuildCommand(),
	}
}
