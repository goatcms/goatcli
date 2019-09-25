package amtf

import (
	"strings"

	"github.com/goatcms/goatcli/cliapp/common/am/entitymodel"
	"github.com/goatcms/goatcli/cliapp/common/naming"
)

// LinkRelationUF create variable link
func LinkRelationUF(varName string, relation *entitymodel.Relation) string {
	if len(relation.Structure.Path) > 0 {
		return varName + "." + strings.Join(relation.Structure.Path, ".") + "." + relation.Name.CamelCaseUF
	}
	return varName + "." + relation.Name.CamelCaseUF
}

// LinkRelationLF create variable link (use camelcase lower first names)
func LinkRelationLF(varName string, relation *entitymodel.Relation) string {
	var steps = []string{varName}
	for _, step := range relation.Structure.Path {
		steps = append(steps, naming.ToCamelCaseLF(step))
	}
	steps = append(steps, relation.Name.CamelCaseLF)
	return strings.Join(steps, ".")
}
