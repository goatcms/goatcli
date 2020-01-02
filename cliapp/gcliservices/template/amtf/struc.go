package amtf

import (
	"strings"

	"github.com/goatcms/goatcli/cliapp/common/am/entitymodel"
	"github.com/goatcms/goatcli/cliapp/common/naming"
)

// StructClassName return structure class name
func StructClassName(struc *entitymodel.Structure) string {
	var steps = []string{struc.Entity.Name.CamelCaseUF}
	steps = append(steps, struc.Path...)
	return strings.Join(steps, "")
}

// LinkStructureUF create variable link
func LinkStructureUF(varName string, struc *entitymodel.Structure) string {
	if len(struc.Path) > 0 {
		return varName + "." + strings.Join(struc.Path, ".")
	}
	return varName
}

// LinkStructureLF create variable link (use camelcase lower first names)
func LinkStructureLF(varName string, struc *entitymodel.Structure) string {
	var steps = []string{varName}
	for _, step := range struc.Path {
		steps = append(steps, naming.ToCamelCaseLF(step))
	}
	return strings.Join(steps, ".")
}
