package depsc

import (
	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/app/terminal"
)

func AddCommand() app.TerminalCommand {
	return terminal.NewCommand(terminal.CommandParams{
		Name:          "deps:add",
		Help:          "add new static dependency (like golang vendor or js node module)",
		Callback:      RunAddDep,
		MainArguments: []string{"path", "repo", "branch", "revision"},
		Arguments: terminal.NewArguments(
			terminal.NewArgument(terminal.ArgumentParams{
				Name:     "path",
				Help:     "the destination dependency path",
				Required: true,
				Type:     app.TerminalTextArgument,
			}),
			terminal.NewArgument(terminal.ArgumentParams{
				Name:     "repo",
				Help:     "the repository URL",
				Required: true,
				Type:     app.TerminalURLArgument,
			}),
			terminal.NewArgument(terminal.ArgumentParams{
				Name:     "branch",
				Help:     "the branch name",
				Required: false,
				Type:     app.TerminalTextArgument,
			}),
			terminal.NewArgument(terminal.ArgumentParams{
				Name:     "revision",
				Help:     "the revision repository name",
				Required: false,
				Type:     app.TerminalTextArgument,
			}),
		),
	})
}

func AddGOCommand() app.TerminalCommand {
	return terminal.NewCommand(terminal.CommandParams{
		Name:          "deps:add:go",
		Help:          "Add new golang dependency",
		Callback:      RunAddDep,
		MainArguments: []string{"repo", "branch", "revision"},
		Arguments: terminal.NewArguments(
			terminal.NewArgument(terminal.ArgumentParams{
				Name:     "repo",
				Help:     "the go repository like 'github.com/goatcms/goatcore'",
				Required: true,
				Type:     app.TerminalURLArgument,
			}),
			terminal.NewArgument(terminal.ArgumentParams{
				Name:     "branch",
				Help:     "the branch name",
				Required: false,
				Type:     app.TerminalTextArgument,
			}),
			terminal.NewArgument(terminal.ArgumentParams{
				Name:     "revision",
				Help:     "the revision repository name",
				Required: false,
				Type:     app.TerminalTextArgument,
			}),
		),
	})
}
func AddGOImportCommand() app.TerminalCommand {
	return terminal.NewCommand(terminal.CommandParams{
		Name:     "deps:add:go:import",
		Help:     "scan project and add imports",
		Callback: RunAddGOImportsDep,
	})
}

func Commands() []app.TerminalCommand {
	return []app.TerminalCommand{
		AddCommand(),
		AddGOCommand(),
		AddGOImportCommand(),
	}
}
