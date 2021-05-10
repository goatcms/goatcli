package gclicore

import (
	"github.com/goatcms/goatcore/dependency"
	"github.com/goatcms/goatcore/varutil/goaterr"
)

// RegisterDependencies is init callback to register module dependencies
func RegisterDependencies(dp dependency.Provider) error {
	return goaterr.ToError(goaterr.AppendError(nil,
		dp.AddDefaultFactory("GCLICoreArguments", ArgumentsFactory),
	))
}
