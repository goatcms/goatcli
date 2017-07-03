package main

import (
	"log"
	"os"

	"github.com/goatcms/goatcli/cliapp"
	"github.com/goatcms/goatcore/app/bootstrap"
	"github.com/goatcms/goatcore/app/modules/terminal"
)

func main() {
	errLogs := log.New(os.Stderr, "", 0)
	app, err := cliapp.NewCLIApp(cliapp.AppName, cliapp.AppVersion)
	if err != nil {
		errLogs.Println(err)
		return
	}
	bootstrap := bootstrap.NewBootstrap(app)
	if err = bootstrap.Register(terminal.NewModule()); err != nil {
		errLogs.Println(err)
		return
	}
	if err = bootstrap.Register(cliapp.NewModule()); err != nil {
		errLogs.Println(err)
		return
	}
	if err := bootstrap.Init(); err != nil {
		errLogs.Println(err)
		return
	}
	if err := bootstrap.Run(); err != nil {
		errLogs.Println(err)
		return
	}
}
