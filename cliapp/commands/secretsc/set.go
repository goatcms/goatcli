package secretsc

import (
	"github.com/goatcms/goatcli/cliapp/services"
	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/filesystem"
	"github.com/goatcms/goatcore/varutil/goaterr"
)

// RunSetSecretValue run command to set new secret value
func RunSetSecretValue(a app.App, ctx app.IOContext) (err error) {
	var (
		deps struct {
			Key            string                  `command:"?$1"`
			Value          string                  `command:"?$2"`
			CurrentFS      filesystem.Filespace    `filespace:"current"`
			SecretsService services.SecretsService `dependency:"SecretsService"`
		}
		secretsData map[string]string
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
	if deps.Value == "" {
		return goaterr.Errorf(ValueParameterIsRequired)
	}
	if secretsData, err = deps.SecretsService.ReadDataFromFS(deps.CurrentFS); err != nil {
		return err
	}
	secretsData[deps.Key] = deps.Value
	return deps.SecretsService.WriteDataToFS(deps.CurrentFS, secretsData)
}
