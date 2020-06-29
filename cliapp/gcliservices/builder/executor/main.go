package executor

import (
	"github.com/goatcms/goatcli/cliapp/common"
	"github.com/goatcms/goatcli/cliapp/gcliservices"
	"github.com/goatcms/goatcore/filesystem"
)

// Task is single task for generator
type Task struct {
	// Template to run
	Template TemplateHandler
	// external data
	DotData interface{}
	// External properties
	BuildProperties common.ElasticData
	// destination filesystem
	FSPath string
}

// TemplateHandler describe template and block to run
type TemplateHandler struct {
	Layout string
	Path   string
	Name   string
}

// SharedData contains shared data for all executor tasks
type SharedData struct {
	VCSData    gcliservices.VCSData
	Data       common.ElasticData
	Properties GlobalProperties
	FS         filesystem.Filespace
	AM         interface{}
}

// GlobalProperties contains task properties
type GlobalProperties struct {
	Project common.ElasticData
	Secrets common.ElasticData
}

// TaskProperties contains task properties
type TaskProperties struct {
	Build   common.ElasticData
	Project common.ElasticData
	Secrets common.ElasticData
}
