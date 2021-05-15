package containerc

import (
	"github.com/goatcms/goatcli/cliapp/gclicommands/containerc/imagepip"
	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/app/terminal"
)

func Commands() []app.TerminalCommand {
	return []app.TerminalCommand{
		terminal.NewCommand(terminal.CommandParams{
			Name:     "container:image",
			Help:     "Build container image by pipeline commands (see pip argument). Run build command to build image from Dockerfile pip and push to push result image to remote repository/repositories.",
			Callback: RunContainerImagePip,
			Arguments: terminal.NewArguments(
				terminal.NewArgument(terminal.ArgumentParams{
					Name:     "pip",
					Help:     "The argument describe build and push flow.",
					Required: true,
					Type:     app.TerminalPIPArgument,
					Commands: terminal.NewCommands(imagepip.Commands()...),
				}),
			),
		}),
	}
}
