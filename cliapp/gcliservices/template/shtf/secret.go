package shtf

import (
	"strings"

	"github.com/goatcms/goatcli/cliapp/common/naming"
)

// SecretEnv return shell secret environment name
func SecretEnv(key string) string {
	return "$SECRET_" + strings.ToUpper(naming.ToUnderscore(key))
}
