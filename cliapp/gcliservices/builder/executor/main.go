package executor

import (
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
	BuildProperties map[string]string
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
	PlainData  map[string]string
	Properties GlobalProperties
	FS         filesystem.Filespace
	AM         interface{}
}

// GlobalProperties contains task properties
type GlobalProperties struct {
	Project map[string]string
	Secrets map[string]string
}

// TaskProperties contains task properties
type TaskProperties struct {
	Build   map[string]string
	Project map[string]string
	Secrets map[string]string
}
