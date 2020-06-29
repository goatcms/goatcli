package gclivarutil

import (
	"github.com/goatcms/goatcli/cliapp/common"
)

// NewElasticData create elastic data from palin string map
func NewElasticData(plain map[string]string) (elastic common.ElasticData, err error) {
	elastic.Plain = plain
	elastic.Tree, err = ToTree(plain)
	return
}
