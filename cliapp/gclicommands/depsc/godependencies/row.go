package godependencies

import (
	"github.com/goatcms/goatcli/cliapp/common/config"
)

// SetRow is single record in Set
type SetRow struct {
	Dependency *config.Dependency
	Imported   bool
}

// SetImported update imported value
func (row *SetRow) SetImported(value bool) {
	row.Imported = value
}
