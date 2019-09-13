package entitymodel

import (
	"strings"

	"github.com/goatcms/goatcli/cliapp/common/naming"
	"github.com/goatcms/goatcore/varutil/goaterr"
)

// Name struct contains entity name
type Name struct {
	Plain       string
	CamelCaseUF string
	CamelCaseLF string
	Lower       string
	Upper       string
	Underscore  string
}

// NewName create new entity name structure
func NewName(str string) (name Name, err error) {
	name.Plain = str
	if str == "" {
		return name, goaterr.Errorf("Name can not be empty")
	}
	name.CamelCaseUF = naming.ToCamelCaseUF(str)
	name.CamelCaseLF = naming.ToCamelCaseLF(str)
	name.Underscore = naming.ToUnderscore(str)
	name.Lower = strings.ToLower(name.CamelCaseUF)
	name.Upper = strings.ToUpper(name.CamelCaseUF)
	return name, nil
}
