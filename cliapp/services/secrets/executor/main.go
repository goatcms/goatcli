package executor

import (
	"github.com/goatcms/goatcli/cliapp/services"
)

// Task is single task for generator
type Task struct {
	TemplateName string
}

// SharedData contains shared data for all executor tasks
type SharedData struct {
	AppData    services.ApplicationData
	Properties GlobalProperties
	DotData    interface{}
}

// GlobalProperties contains task properties
type GlobalProperties struct {
	Project map[string]string
}
