package gcliservices

import (
	"github.com/goatcms/goatcli/cliapp/common"
	"github.com/goatcms/goatcli/cliapp/common/config"
	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/filesystem"
)

// BuilderService build project structure
type BuilderService interface {
	ReadDefFromFS(fs filesystem.Filespace) ([]*config.Build, error)
	Build(ctx app.IOContext, fs filesystem.Filespace, appData ApplicationData, properties, secrets common.ElasticData) (err error)
}
