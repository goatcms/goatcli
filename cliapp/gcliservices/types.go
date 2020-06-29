package gcliservices

import "github.com/goatcms/goatcli/cliapp/common"

// ApplicationData provide data api
type ApplicationData struct {
	ElasticData common.ElasticData
	AM          interface{}
}
