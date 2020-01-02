package secretsc

import (
	"github.com/goatcms/goatcli/cliapp/gcliservices"
	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/filesystem"
	"github.com/goatcms/goatcore/varutil/goaterr"
)

// RunGetSecretValue run command return secret value
func RunGetSecretValue(a app.App, ctx app.IOContext) (err error) {
	var (
		deps struct {
			Key            string                      `command:"?$1"`
			CurrentFS      filesystem.Filespace        `filespace:"current"`
			SecretsService gcliservices.SecretsService `dependency:"SecretsService"`
		}
		secretsData map[string]string
		ok          bool
		value       string
	)
	if err = a.DependencyProvider().InjectTo(&deps); err != nil {
		return err
	}
	if err = ctx.Scope().InjectTo(&deps); err != nil {
		return err
	}
	if deps.Key == "" {
		return goaterr.Errorf(FirstKeyParameterIsRequired)
	}
	if secretsData, err = deps.SecretsService.ReadDataFromFS(deps.CurrentFS); err != nil {
		return err
	}
	if value, ok = secretsData[deps.Key]; !ok {
		return goaterr.Errorf("Unknow Value for %s key", deps.Key)
	}
	return ctx.IO().Out().Printf(value)
}
