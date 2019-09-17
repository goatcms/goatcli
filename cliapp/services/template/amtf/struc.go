package amtf

import (
	"strings"

	"github.com/goatcms/goatcli/cliapp/common/am/entitymodel"
)

// StructClassName return structure class name
func StructClassName(struc *entitymodel.Structure) string {
	var steps = []string{struc.Entity.Name.CamelCaseUF}
	steps = append(steps, struc.Path...)
	return strings.Join(steps, "")
}
