package cleanc

import (
	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/app/terminal"
)

func CleanCommand() app.TerminalCommand {
	return terminal.NewCommand(terminal.CommandParams{
		Name:     "clean",
		Help:     "clean builded files and dependencies",
		Callback: RunClean,
	})
}

func CleanBuildCommand() app.TerminalCommand {
	return terminal.NewCommand(terminal.CommandParams{
		Name:     "clean:build",
		Help:     "clean build files",
		Callback: RunCleanBuild,
	})
}

func CleanDependenciesCommand() app.TerminalCommand {
	return terminal.NewCommand(terminal.CommandParams{
		Name:     "clean:dependencies",
		Help:     "clean dependencies files",
		Callback: RunCleanDependencies,
	})
}

func Commands() []app.TerminalCommand {
	return []app.TerminalCommand{
		CleanBuildCommand(),
		CleanCommand(),
		CleanDependenciesCommand(),
	}
}
