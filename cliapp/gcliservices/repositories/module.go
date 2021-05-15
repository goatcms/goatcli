package repositories

import "github.com/goatcms/goatcore/app"

// RegisterDependencies is init callback to register module dependencies
func RegisterDependencies(dp app.DependencyProvider) (err error) {
	if err = dp.AddDefaultFactory("RepositoriesService", Factory); err != nil {
		return err
	}
	return dp.AddDefaultFactory("RepositoriesConnector", ConnectorFactory)
}
