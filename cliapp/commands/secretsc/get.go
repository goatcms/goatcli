package secretsc

import (
	"fmt"

	"github.com/goatcms/goatcli/cliapp/services"
	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/filesystem"
)

// RunGetSecretValue run command return secret value
func RunGetSecretValue(a app.App) (err error) {
	var (
		deps struct {
			Key            string                  `argument:"?$2"`
			CurrentFS      filesystem.Filespace    `filespace:"current"`
			SecretsService services.SecretsService `dependency:"SecretsService"`
			Input          app.Input               `dependency:"InputService"`
			Output         app.Output              `dependency:"OutputService"`
		}
		secretsData map[string]string
		ok          bool
		value       string
	)
	if err = a.DependencyProvider().InjectTo(&deps); err != nil {
		return err
	}
	if deps.Key == "" {
		deps.Output.Printf(FirstKeyParameterIsRequired)
		return fmt.Errorf(FirstKeyParameterIsRequired)
	}
	if secretsData, err = deps.SecretsService.ReadDataFromFS(deps.CurrentFS); err != nil {
		return err
	}
	if value, ok = secretsData[deps.Key]; !ok {
		return fmt.Errorf("Unknow Value for %s key", deps.Key)
	}
	return deps.Output.Printf(value)
}
