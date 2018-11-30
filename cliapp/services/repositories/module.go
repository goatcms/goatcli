package repositories

import "github.com/goatcms/goatcore/dependency"

// RegisterDependencies is init callback to register module dependencies
func RegisterDependencies(dp dependency.Provider) (err error) {
	if err = dp.AddDefaultFactory("RepositoriesService", Factory); err != nil {
		return err
	}
	return dp.AddDefaultFactory("RepositoriesConnector", ConnectorFactory)
}
