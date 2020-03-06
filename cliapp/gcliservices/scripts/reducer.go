package scripts

import (
	"github.com/goatcms/goatcli/cliapp/gcliservices"
	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/dependency"
)

// Reducer run application scripts
type Reducer struct {
	deps struct {
		CWD           string                     `argument:"?cwd"`
		ScriptsRunner gcliservices.ScriptsRunner `dependency:"ScriptsRunner"`
	}
}

// ReducerFactory create new Reducer instance
func ReducerFactory(dp dependency.Provider) (result interface{}, err error) {
	r := &Reducer{}
	if err = dp.InjectTo(&r.deps); err != nil {
		return nil, err
	}
	return gcliservices.ScriptsReducer(r), nil
}

// Run reduce script and run
func (runner *Reducer) Run(ctx app.IOContext, reducerName string, params gcliservices.ScriptsRunnerParams) (err error) {
	//TODO:
	return nil
}
