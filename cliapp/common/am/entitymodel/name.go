package entitymodel

import (
	"fmt"
	"strings"

	"github.com/goatcms/goatcli/cliapp/common/naming"
)

// Name struct contains entity name
type Name struct {
	Plain       string
	CamelCaseUF string
	CamelCaseLF string
	Underscore  string
	Lower       string
	Upper       string
}

// NewName create new entity name structure
func NewName(str string) (name Name, err error) {
	name.Plain = str
	if str == "" {
		return name, fmt.Errorf("Name can not be empty")
	}
	name.CamelCaseUF = naming.ToCamelCaseUF(str)
	name.CamelCaseLF = naming.ToCamelCaseLF(str)
	name.Underscore = naming.ToUnderscore(str)
	name.Lower = strings.ToLower(name.CamelCaseUF)
	name.Upper = strings.ToUpper(name.CamelCaseUF)
	return name, nil
}
