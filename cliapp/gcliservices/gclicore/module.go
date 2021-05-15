package gclicore

import (
	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/varutil/goaterr"
)

// RegisterDependencies is init callback to register module dependencies
func RegisterDependencies(dp app.DependencyProvider) error {
	return goaterr.ToError(goaterr.AppendError(nil,
		dp.AddDefaultFactory("GCLICoreArguments", ArgumentsFactory),
	))
}
