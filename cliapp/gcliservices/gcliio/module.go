package gcliio

import "github.com/goatcms/goatcore/app"

// RegisterDependencies is init callback to register module dependencies
func RegisterDependencies(dp app.DependencyProvider) (err error) {
	if err = dp.AddDefaultFactory("GCLIInputs", InputsFactory); err != nil {
		return err
	}
	return dp.AddDefaultFactory("GCLIEnvironment", EnvironmentFactory)
}
