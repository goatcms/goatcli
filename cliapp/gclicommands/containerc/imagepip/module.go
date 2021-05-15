package imagepip

import (
	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/app/modules/pipelinem/pipcommands/pipc"
	"github.com/goatcms/goatcore/app/terminal"
)

// BuildCommand return build command definition
func BuildCommand() app.TerminalCommand {
	return terminal.NewCommand(terminal.CommandParams{
		Callback: RunBuild,
		Help:     "build container image",
		Name:     "build",
		Arguments: terminal.NewArguments([]app.TerminalArgument{
			terminal.NewArgument(terminal.ArgumentParams{
				Help:     "steps to build. Use syntax from Dockerfile (https://docs.docker.com/engine/reference/builder/)",
				Name:     "steps",
				Required: true,
				Type:     app.TerminalTextArgument,
			}),
		}...),
	})
}

// PushCommand return push command definition
func PushCommand() app.TerminalCommand {
	return terminal.NewCommand(terminal.CommandParams{
		Callback: RunPush,
		Help:     "push image to remote repository. ",
		Name:     "push",
		Arguments: terminal.NewArguments([]app.TerminalArgument{
			terminal.NewArgument(terminal.ArgumentParams{
				Help:     "your login for remote ropository (insert seret key name)",
				Name:     "login",
				Required: true,
				Type:     app.TerminalTextArgument,
			}),
			terminal.NewArgument(terminal.ArgumentParams{
				Help:     "your password for remote ropository (insert seret key name)",
				Name:     "password",
				Required: true,
				Type:     app.TerminalTextArgument,
			}),
			terminal.NewArgument(terminal.ArgumentParams{
				Help:     "full destination repository like docker.io/library/hello-world or quay.io/libpod/alpine",
				Name:     "dest",
				Required: true,
				Type:     app.TerminalTextArgument,
			}),
			terminal.NewArgument(terminal.ArgumentParams{
				Help:     "tls-verify turn on / off tls verification",
				Name:     "tls-verify",
				Required: false,
				Type:     app.TerminalBoolArgument,
			}),
		}...),
	})
}

// Commands return image pipeline commands list
func Commands() []app.TerminalCommand {
	return []app.TerminalCommand{
		BuildCommand(),
		PushCommand(),
		pipc.RunCommand(),
		pipc.TryCommand(),
	}
}
