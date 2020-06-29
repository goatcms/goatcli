package scripts

import (
	"github.com/goatcms/goatcli/cliapp/common"
)

// Context is script context data
type Context struct {
	Data       common.ElasticData
	AM         interface{}
	Properties TaskProperties
}

// TaskProperties contains task properties
type TaskProperties struct {
	Project common.ElasticData
	Secrets common.ElasticData
}
