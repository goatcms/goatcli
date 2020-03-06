package gcliservices

import (
	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/filesystem"
)

// ScriptsRunnerParams is data object for ScriptsRunner
type ScriptsRunnerParams struct {
	FS         filesystem.Filespace
	ScriptName string
	Once       bool
	Secrets    map[string]string
	Properties map[string]string
	Data       ApplicationData
}

// ScriptsRunner run project scripts
type ScriptsRunner interface {
	RunByName(ctx app.IOContext, params ScriptsRunnerParams) (err error)
	// Runned return true if script by name was runed in the context
	IsRunnedInScope(scp app.Scope, name string) (is bool, err error)
	// Runned return true if script by name was runed
	IsRunned(name string) (is bool, err error)
}

// ScriptsQueueReducer reduce, run and control execution of project scripts
type ScriptsQueueReducer interface {
	Run(ctx app.IOContext, queueType string, params ScriptsRunnerParams) (err error)
}
