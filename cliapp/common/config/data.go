package config

import (
	"fmt"
	"strings"

	"github.com/buger/jsonparser"
)

// DataSet is data definition structure
type DataSet struct {
	TypeName    string
	Properties  []*Property
	Collections []*Collection
}

// NewDataSets parse json and return DataSet array instance
func NewDataSets(json []byte) (dataSets []*DataSet, de error) {
	dataSets = []*DataSet{}
	if _, err := jsonparser.ArrayEach(json, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		if err != nil || de != nil {
			return
		}
		if dataType != jsonparser.Object {
			de = fmt.Errorf("NewReplaces array  must contains replace objects only")
			return
		}
		dataSet, err2 := NewDataSet(value)
		if err2 != nil {
			de = err2
			return
		}
		dataSets = append(dataSets, dataSet)
	}); err != nil {
		return nil, err
	}
	if de != nil {
		return nil, de
	}
	return dataSets, nil
}

// NewDataSet parse replace data and return DataSet object instance
func NewDataSet(json []byte) (dataSet *DataSet, err error) {
	dataSet = &DataSet{}
	if err = jsonparser.ObjectEach(json, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
		switch strings.ToLower(string(key)) {
		case "type":
			if dataType != jsonparser.String {
				return fmt.Errorf("DataSet.TypeName expected string and take %s", value)
			}
			dataSet.TypeName = string(value)
		case "properties":
			if dataType != jsonparser.Array {
				return fmt.Errorf("DataSet.Properties expected array and take %s", value)
			}
			if dataSet.Properties, err = NewProperties(value); err != nil {
				return err
			}
		case "collections":
			if dataType != jsonparser.Array {
				return fmt.Errorf("DataSet.Collections expected array and take %s", value)
			}
			if dataSet.Collections, err = NewCollections(value); err != nil {
				return err
			}
		case commentProperty:
			// ignore all comments
		default:
			return fmt.Errorf("config.NewReplace: Unknow key %s (value: %s)", key, value)
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return dataSet, nil
}
