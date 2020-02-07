package gcliservices

import (
	"github.com/goatcms/goatcore/app"
)

// GCLIInputs return goat cli application inputs
type GCLIInputs interface {
	Inputs(ctx app.IOContext) (propertiesData, secretsData map[string]string, appData ApplicationData, err error)
}

// GCLIEnvironment manage environments
type GCLIEnvironment interface {
	LoadEnvs(ctxScope app.Scope, propertiesData, secretsData map[string]string) (err error)
}
