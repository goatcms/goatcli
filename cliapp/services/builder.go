package services

import (
	"github.com/goatcms/goatcli/cliapp/common/config"
	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/filesystem"
)

// BuildContext is data for builer
type BuildContext interface {
	Build(fs filesystem.Filespace) error
}

// BuilderService build project structure
type BuilderService interface {
	ReadDefFromFS(fs filesystem.Filespace) ([]*config.Build, error)
	NewContext(scope app.Scope, data, properties, secrets map[string]string) BuildContext
}
