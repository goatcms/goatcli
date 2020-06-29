package executor

import (
	"github.com/goatcms/goatcli/cliapp/common"
	"github.com/goatcms/goatcli/cliapp/gcliservices"
)

// Task is single task for generator
type Task struct {
	TemplateName string
}

// SharedData contains shared data for all executor tasks
type SharedData struct {
	AppData    gcliservices.ApplicationData
	Properties GlobalProperties
	DotData    interface{}
}

// GlobalProperties contains task properties
type GlobalProperties struct {
	Project common.ElasticData
}
