package main

import (
	"fmt"
	"os"

	"github.com/goatcms/goatcli/cliapp"
	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/app/bootstrap"
	"github.com/goatcms/goatcore/app/goatapp"
	"github.com/goatcms/goatcore/app/modules/commonm"
	"github.com/goatcms/goatcore/app/modules/ocm"
	"github.com/goatcms/goatcore/app/modules/pipelinem"
	"github.com/goatcms/goatcore/app/modules/systemm"
	"github.com/goatcms/goatcore/app/modules/terminalm"
	"github.com/goatcms/goatcore/varutil/goaterr"
)

func main() {
	var (
		gapp app.App
		boot *bootstrap.Bootstrap
		err  error
		code int
	)
	if gapp, err = goatapp.NewGoatApp(goatapp.Params{
		Name:    "GoatCLI",
		Version: goatapp.NewVersion(0, 0, 3, "-dev"),
	}); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	boot = bootstrap.NewBootstrap(gapp)
	if err = goaterr.ToError(goaterr.AppendError(nil,
		boot.Register(terminalm.NewModule()),
		boot.Register(commonm.NewModule()),
		boot.Register(ocm.NewModule()),
		boot.Register(pipelinem.NewModule()),
		boot.Register(cliapp.NewModule()),
		boot.Register(systemm.NewModule()),
	)); err != nil {
		fmt.Println(err)
		os.Exit(2)
		return
	}
	if err = boot.Init(); err != nil {
		fmt.Println(err)
		os.Exit(3)
		return
	}
	if err = boot.Run(); err != nil {
		fmt.Printf("\n")
		if code, err = boot.ShowError(err); err != nil {
			fmt.Println(err)
		}
		os.Exit(code)
	}
}
