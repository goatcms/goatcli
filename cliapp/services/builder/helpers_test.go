package builder

import (
	"github.com/goatcms/goatcli/cliapp/common/am"
	"github.com/goatcms/goatcli/cliapp/services"
	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/app/scope"
	"github.com/goatcms/goatcore/filesystem"
)

type RenderFileDeps struct {
	BuilderService services.BuilderService `dependency:"BuilderService"`
}

func renderFile(fs filesystem.Filespace, deps RenderFileDeps, data, properties, secrets map[string]string) (err error) {
	ctxScope := scope.NewScope("test")
	appData := am.NewApplicationData(data)
	buildContext := deps.BuilderService.NewContext(ctxScope, appData, properties, secrets)
	if err = buildContext.Build(fs); err != nil {
		return err
	}
	if err = ctxScope.Wait(); err != nil {
		return err
	}
	return ctxScope.Trigger(app.CommitEvent, nil)
}
