package gcliservices

import (
	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/app/modules/pipelinem/pipservices"
	"github.com/goatcms/goatcore/filesystem"
)

// ScriptsRunner run project scripts
type ScriptsRunner interface {
	Run(ctx ScriptsContext, fs filesystem.Filespace, scriptName string, properties, secrets map[string]string, data ApplicationData) (taskManager pipservices.TasksManager, err error)
}

// ScriptsContext contains script context
type ScriptsContext struct {
	Scope      app.Scope
	CWD        filesystem.Filespace
	Namespaces pipservices.Namespaces
}
