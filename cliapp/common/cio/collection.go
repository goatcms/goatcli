package cio

import (
	"fmt"
	"strings"

	"github.com/goatcms/goatcli/cliapp/common/config"
	"github.com/goatcms/goatcore/app"
)

// ReadCollections read collections to data
func ReadCollections(baseKey string, in app.Input, out app.Output, def []*config.Collection, data map[string]string) (isChanged bool, err error) {
	var isCollChanged bool
	for _, coll := range def {
		if isCollChanged, err = ReadCollection(baseKey, in, out, coll, data); err != nil {
			return isChanged, err
		}
		isChanged = isChanged || isCollChanged
	}
	return isChanged, nil
}

// ReadCollection read collection to data
func ReadCollection(baseKey string, in app.Input, out app.Output, coll *config.Collection, data map[string]string) (isChanged bool, err error) {
	out.Printf("\n - (%s*) %s\n", baseKey, coll.Prompt)
	switch strings.ToLower(coll.Type) {
	case config.ArrayCollection:
		return ReadPropertiesArray(baseKey, in, out, coll.Properties, data)
	case config.MapCollection:
		return ReadPropertiesMap(baseKey, in, out, coll.Properties, data)
	default:
		return false, fmt.Errorf("cio.ReadCollection: incorrect type %s (expected 'map' or 'array')", coll.Type)
	}
}
