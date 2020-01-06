package gcliio

import "github.com/goatcms/goatcore/dependency"

// RegisterDependencies is init callback to register module dependencies
func RegisterDependencies(dp dependency.Provider) (err error) {
	return dp.AddDefaultFactory("GCLIInputs", InputsFactory)
}
