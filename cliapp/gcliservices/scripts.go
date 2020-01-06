package gcliservices

import (
	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/filesystem"
)

// ScriptsRunner run project scripts
type ScriptsRunner interface {
	Run(ctx app.IOContext, fs filesystem.Filespace, scriptName string, properties, secrets map[string]string, data ApplicationData) (err error)
}
