package main

import (
	"log"
	"os"

	"github.com/goatcms/goatcli/cliapp"
	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/app/bootstrap"
	"github.com/goatcms/goatcore/app/goatapp"
	"github.com/goatcms/goatcore/app/modules/commonm"
	"github.com/goatcms/goatcore/app/modules/pipelinem"
	"github.com/goatcms/goatcore/app/modules/terminalm"
	"github.com/goatcms/goatcore/varutil/goaterr"
)

func main() {
	var (
		gapp app.App
		err  error
		boot app.Bootstrap
	)
	errLogs := log.New(os.Stderr, "", 0)
	if gapp, err = goatapp.NewGoatApp(cliapp.AppName, cliapp.AppVersion, "./"); err != nil {
		errLogs.Println(err)
		return
	}
	boot = bootstrap.NewBootstrap(gapp)

	if err = goaterr.ToError(goaterr.AppendError(nil,
		boot.Register(terminalm.NewModule()),
		boot.Register(commonm.NewModule()),
		boot.Register(pipelinem.NewModule()),
		boot.Register(cliapp.NewModule()),
	)); err != nil {
		errLogs.Println(err)
		return
	}
	if err := boot.Init(); err != nil {
		errLogs.Println(err)
		return
	}
	if err := boot.Run(); err != nil {
		errLogs.Println(err)
		return
	}
}
