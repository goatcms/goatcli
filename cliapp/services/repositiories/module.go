package repositories

import "github.com/goatcms/goatcore/dependency"

// RegisterDependencies is init callback to register module dependencies
func RegisterDependencies(dp dependency.Provider) error {
	return dp.AddDefaultFactory("RepositoriesService", Factory)
}
