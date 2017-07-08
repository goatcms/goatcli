package config

import (
	"fmt"
	"strings"

	"github.com/buger/jsonparser"
)

const (
	// ArrayCollection is array type name
	ArrayCollection = "array"
	// MapCollection is map type name
	MapCollection = "map"
)

// Collection is a properties collection definition
type Collection struct {
	Key        string
	Type       string
	Prompt     string
	Properties []*Property
}

// NewCollections parse json and return DataSet array instance
func NewCollections(json []byte) (colls []*Collection, de error) {
	colls = []*Collection{}
	if _, err := jsonparser.ArrayEach(json, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		if err != nil || de != nil {
			return
		}
		if dataType != jsonparser.Object {
			de = fmt.Errorf("NewCollections: array must contains objects")
			return
		}
		coll, err2 := NewCollection(value)
		if err2 != nil {
			de = err2
			return
		}
		colls = append(colls, coll)
	}); err != nil {
		return nil, err
	}
	if de != nil {
		return nil, de
	}
	return colls, nil
}

// NewCollection parse collection json and return Collection object instance
func NewCollection(json []byte) (coll *Collection, err error) {
	coll = &Collection{}
	if err = jsonparser.ObjectEach(json, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
		switch strings.ToLower(string(key)) {
		case "type":
			if dataType != jsonparser.String {
				return fmt.Errorf("Collection.Type expected string and take %s", value)
			}
			coll.Type = string(value)
			if coll.Type != MapCollection && coll.Type != ArrayCollection {
				return fmt.Errorf("Collection.Type must be equals to 'map' or 'array' and it is equals to '%s'", coll.Type)
			}
		case "key":
			if dataType != jsonparser.String {
				return fmt.Errorf("Collection.Key expected string and take %s", value)
			}
			coll.Key = string(value)
		case "properties":
			if dataType != jsonparser.Array {
				return fmt.Errorf("Collection.Properties expected array and take %s", value)
			}
			if coll.Properties, err = NewProperties(value); err != nil {
				return err
			}
		case "prompt":
			if dataType != jsonparser.String {
				return fmt.Errorf("property.prompt expected string and take %s", value)
			}
			coll.Prompt = string(value)
		case "comment":
			// ignore all comments
		default:
			return fmt.Errorf("Collection.NewCollection: Unknow key %s (value: %s)", key, value)
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return coll, nil
}
