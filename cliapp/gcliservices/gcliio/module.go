package gcliio

import "github.com/goatcms/goatcore/dependency"

// RegisterDependencies is init callback to register module dependencies
func RegisterDependencies(dp dependency.Provider) (err error) {
	if err = dp.AddDefaultFactory("GCLIInputs", InputsFactory); err != nil {
		return err
	}
	return dp.AddDefaultFactory("GCLIEnvironment", EnvironmentFactory)
}
