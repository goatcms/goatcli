package secretsc

import (
	"fmt"

	"github.com/goatcms/goatcli/cliapp/services"
	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/filesystem"
)

// RunSetSecretValue run command to set new secret value
func RunSetSecretValue(a app.App) (err error) {
	var (
		deps struct {
			Key            string                  `argument:"?$2"`
			Value          string                  `argument:"?$3"`
			CurrentFS      filesystem.Filespace    `filespace:"current"`
			SecretsService services.SecretsService `dependency:"SecretsService"`
			Input          app.Input               `dependency:"InputService"`
			Output         app.Output              `dependency:"OutputService"`
		}
		secretsData map[string]string
	)
	if err = a.DependencyProvider().InjectTo(&deps); err != nil {
		return err
	}
	if deps.Key == "" {
		deps.Output.Printf(FirstKeyParameterIsRequired)
		return fmt.Errorf(FirstKeyParameterIsRequired)
	}
	if deps.Value == "" {
		deps.Output.Printf(ValueParameterIsRequired)
		return fmt.Errorf(ValueParameterIsRequired)
	}
	if secretsData, err = deps.SecretsService.ReadDataFromFS(deps.CurrentFS); err != nil {
		return err
	}
	secretsData[deps.Key] = deps.Value
	return deps.SecretsService.WriteDataToFS(deps.CurrentFS, secretsData)
}
