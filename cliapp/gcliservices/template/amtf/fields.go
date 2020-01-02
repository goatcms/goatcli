package amtf

import (
	"strings"

	"github.com/goatcms/goatcli/cliapp/common/am/entitymodel"
	"github.com/goatcms/goatcli/cliapp/common/naming"
)

// LinkFieldUF create variable link
func LinkFieldUF(varName string, field *entitymodel.Field) string {
	if len(field.Structure.Path) > 0 {
		return varName + "." + strings.Join(field.Structure.Path, ".") + "." + field.Name.CamelCaseUF
	}
	return varName + "." + field.Name.CamelCaseUF
}

// LinkFieldLF create variable link (use camelcase lower first names)
func LinkFieldLF(varName string, field *entitymodel.Field) string {
	var steps = []string{varName}
	for _, step := range field.Structure.Path {
		steps = append(steps, naming.ToCamelCaseLF(step))
	}
	steps = append(steps, field.Name.CamelCaseLF)
	return strings.Join(steps, ".")
}
