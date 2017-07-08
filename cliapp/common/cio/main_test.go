package cio

import "github.com/goatcms/goatcli/cliapp/common/config"

const (
	testPropDefJSON = `[{"key":"key1", "type":"alnum", "min":1, "max":22},{"key":"key2", "type":"alnum", "min":1, "max":22}]`
)

var (
	testProperties = []*config.Property{
		&config.Property{
			Key:  "key1",
			Type: "alnum",
			Min:  1,
			Max:  22,
		},
		&config.Property{
			Key:  "key2",
			Type: "alnum",
			Min:  1,
			Max:  22,
		},
	}
)
