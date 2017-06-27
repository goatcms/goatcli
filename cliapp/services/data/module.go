package compiler

import "github.com/goatcms/goatcore/dependency"

// RegisterDependencies is init callback to register module dependencies
func RegisterDependencies(dp dependency.Provider) error {
	if err := dp.AddDefaultFactory("DataService", BuilderFactory); err != nil {
		return err
	}
	return nil
}
