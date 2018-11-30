package config

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/buger/jsonparser"
	"github.com/goatcms/goatcore/varutil/totype"
)

// Property is a record describe one property
type Property struct {
	Key     string
	Type    string
	Prompt  string
	Min     int
	Max     int
	Pattern *regexp.Regexp
}

// NewProperties parse json and return Property array instance
func NewProperties(json []byte) ([]*Property, error) {
	var de error
	properties := []*Property{}
	if _, err := jsonparser.ArrayEach(json, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		if err != nil || de != nil {
			return
		}
		if dataType != jsonparser.Object {
			de = fmt.Errorf("NewProperties array must contains property objects only")
			return
		}
		property, err2 := NewProperty(value)
		if err2 != nil {
			de = err2
			return
		}
		properties = append(properties, property)
	}); err != nil {
		return nil, err
	}
	if de != nil {
		return nil, de
	}
	return properties, nil
}

// NewProperty parse property data and return Property object instance
func NewProperty(json []byte) (*Property, error) {
	var err error
	p := &Property{}
	if err = jsonparser.ObjectEach(json, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
		switch strings.ToLower(string(key)) {
		case "key":
			if dataType != jsonparser.String {
				return fmt.Errorf("expected string and take %s", value)
			}
			p.Key = string(value)
		case "type":
			if dataType != jsonparser.String {
				return fmt.Errorf("expected string and take %s", value)
			}
			p.Type = string(value)
		case "min":
			if dataType != jsonparser.Number {
				return fmt.Errorf("expected int and take %s", value)
			}
			if p.Min, err = totype.StringToInt(string(value)); err != nil {
				return err
			}
		case "max":
			if dataType != jsonparser.Number {
				return fmt.Errorf("expected int and take %s", value)
			}
			if p.Max, err = totype.StringToInt(string(value)); err != nil {
				return err
			}
		case "pattern":
			if dataType != jsonparser.String {
				return fmt.Errorf("property.pattern expected string and take %s", value)
			}
			if p.Pattern, err = regexp.Compile(string(value)); err != nil {
				return err
			}
		case "prompt":
			if dataType != jsonparser.String {
				return fmt.Errorf("property.prompt expected string and take %s", value)
			}
			p.Prompt = string(value)
		case commentProperty:
			// ignore all comments
		default:
			return fmt.Errorf("config.NewProperty: Unknow key %s (value: %s)", key, value)
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return p, nil
}
