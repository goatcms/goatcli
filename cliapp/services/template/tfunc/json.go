package tfunc

import (
	"fmt"

	"github.com/goatcms/goatcore/varutil/plainmap"
)

// ToJSON convert plainmap to json
func ToJSON(data map[string]string) (json string, err error) {
	if json, err = plainmap.PlainStringMapToFormattedJSON(data); err != nil {
		LogCriticError(fmt.Sprintf("ToJSON Error: %v", err))
		return "", err
	}
	return json, nil
}
