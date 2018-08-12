package cio

import (
	"github.com/goatcms/goatcli/cliapp/common/config"
	"github.com/goatcms/goatcore/app"
)

// ReadDataSet read properties from app.Input
func ReadDataSet(baseKey string, in app.Input, out app.Output, def *config.DataSet, data map[string]string) (isChanged bool, err error) {
	var isChangedTmp bool
	if isChangedTmp, err = ReadProperties(baseKey, in, out, def.Properties, data, map[string]string{}, true); err != nil {
		return false, err
	}
	isChanged = isChangedTmp || isChanged
	if isChangedTmp, err = ReadCollections(baseKey, in, out, def.Collections, data); err != nil {
		return false, err
	}
	return isChangedTmp || isChanged, nil
}
