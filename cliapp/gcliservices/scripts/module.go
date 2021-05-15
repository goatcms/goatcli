package scripts

import "github.com/goatcms/goatcore/app"

// RegisterDependencies is init callback to register module dependencies
func RegisterDependencies(dp app.DependencyProvider) (err error) {
	return dp.AddDefaultFactory("ScriptsRunner", RunnerFactory)
}
