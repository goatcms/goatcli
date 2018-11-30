package repositories

import (
	"github.com/goatcms/goatcore/dependency"
	"github.com/goatcms/goatcore/repositories"
	"github.com/goatcms/goatcore/repositories/git"
)

// ConnectorFactory create new repository connector instance
func ConnectorFactory(dp dependency.Provider) (interface{}, error) {
	connector := repositories.NewMultiConnector([]repositories.ConnectorAdapter{
		git.NewConnector(),
	})
	return repositories.Connector(connector), nil
}
