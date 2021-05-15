package vcsc

import (
	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/app/terminal"
)

func CleanCommand() app.TerminalCommand {
	return terminal.NewCommand(terminal.CommandParams{
		Callback: RunClean,
		Help:     "clean vcs persisted files",
		Name:     "vcs:clean",
	})
}

func GeneratedListCommand() app.TerminalCommand {
	return terminal.NewCommand(terminal.CommandParams{
		Callback: RunGeneratedList,
		Help:     "show generated files listing",
		Name:     "vcs:generated:list",
	})
}

func PersistedAddCommand() app.TerminalCommand {
	return terminal.NewCommand(terminal.CommandParams{
		Callback: RunPersistedAdd,
		Help:     "add new vcs persisted file",
		Name:     "vcs:persisted:add",
		Arguments: terminal.NewArguments(
			terminal.NewArgument(terminal.ArgumentParams{
				Name:     "path",
				Help:     "the path of file to add",
				Required: true,
				Type:     app.TerminalBoolArgument,
			}),
		),
	})
}

func PersistedRemoveCommand() app.TerminalCommand {
	return terminal.NewCommand(terminal.CommandParams{
		Callback: RunPersistedRemove,
		Help:     "remove a vcs persisted file",
		Name:     "vcs:persisted:remove",
		Arguments: terminal.NewArguments(
			terminal.NewArgument(terminal.ArgumentParams{
				Name:     "path",
				Help:     "the path of file to remove",
				Required: true,
				Type:     app.TerminalBoolArgument,
			}),
		),
	})
}

func PersistedListCommand() app.TerminalCommand {
	return terminal.NewCommand(terminal.CommandParams{
		Callback: RunPersistedList,
		Help:     "show persisted files listing",
		Name:     "vcs:persisted:list",
	})
}

func ScanCommand() app.TerminalCommand {
	return terminal.NewCommand(terminal.CommandParams{
		Callback: RunScan,
		Help:     "scan files for changes and add it to vcs persisted files",
		Name:     "vcs:scan",
		Arguments: terminal.NewArguments(
			terminal.NewArgument(terminal.ArgumentParams{
				Name:     "interactive",
				Help:     "the interactive turn on / off interactive mode",
				Required: true,
				Type:     app.TerminalBoolArgument,
			}),
		),
	})
}

func Commands() []app.TerminalCommand {
	return []app.TerminalCommand{
		CleanCommand(),
		GeneratedListCommand(),
		PersistedAddCommand(),
		PersistedListCommand(),
		PersistedRemoveCommand(),
		ScanCommand(),
	}
}
