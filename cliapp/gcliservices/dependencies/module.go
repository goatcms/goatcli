package dependencies

import "github.com/goatcms/goatcore/app"

// RegisterDependencies is init callback to register module dependencies
func RegisterDependencies(dp app.DependencyProvider) error {
	return dp.AddDefaultFactory("DependenciesService", Factory)
}
