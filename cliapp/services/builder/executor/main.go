package executor

import (
	"github.com/goatcms/goatcli/cliapp/services"
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
	VCSData    services.VCSData
	PlainData  map[string]string
	Properties GlobalProperties
	FS         filesystem.Filespace
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
