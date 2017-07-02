package builder

import "github.com/goatcms/goatcore/dependency"

// RegisterDependencies is init callback to register module dependencies
func RegisterDependencies(dp dependency.Provider) error {
	if err := dp.AddDefaultFactory("BuilderService", ServiceFactory); err != nil {
		return err
	}
	return nil
}
