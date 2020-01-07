package cliapp

import (
	"os/exec"

	"github.com/goatcms/goatcore/app"
)

// GitHealthChecker check if system contains git
func GitHealthChecker(a app.App, ctxScope app.Scope) (msg string, err error) {
	if err = exec.Command("git", "version").Run(); err != nil {
		return "Workers goatcli require pre-installed git tools (install: https://git-scm.com/ )", err
	}
	return "Git", nil
}
