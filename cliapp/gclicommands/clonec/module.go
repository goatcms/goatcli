package clonec

import (
	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/app/terminal"
)

func CloneCommand() app.TerminalCommand {
	return terminal.NewCommand(terminal.CommandParams{
		Name:          "clone",
		Help:          "clone project (and modules)",
		Callback:      Run,
		MainArguments: []string{"repo", "dest"},
		Arguments: terminal.NewArguments(
			terminal.NewArgument(terminal.ArgumentParams{
				Name:     "repo",
				Help:     "git repository url",
				Required: true,
			}),
			terminal.NewArgument(terminal.ArgumentParams{
				Name:     "dest",
				Help:     "local destination path",
				Required: true,
			}),
			terminal.NewArgument(terminal.ArgumentParams{
				Name:     "branch",
				Help:     "the branch name",
				Required: false,
			}),
			terminal.NewArgument(terminal.ArgumentParams{
				Name:     "rev",
				Help:     "repository revision",
				Required: false,
			}),
		),
	})
}

func Commands() []app.TerminalCommand {
	return []app.TerminalCommand{
		CloneCommand(),
	}
}
