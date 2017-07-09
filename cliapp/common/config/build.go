package config

import (
	"fmt"
	"strings"

	"github.com/buger/jsonparser"
	"github.com/goatcms/goatcore/varutil/plainmap"
)

// Build is data definition structure
type Build struct {
	From       string
	To         string
	Layout     string
	Template   string
	Properties map[string]string
}

// NewBuilds parse json and return Build array instance
func NewBuilds(json []byte) (builds []*Build, de error) {
	builds = []*Build{}
	if _, err := jsonparser.ArrayEach(json, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		if err != nil || de != nil {
			return
		}
		if dataType != jsonparser.Object {
			de = fmt.Errorf("NewBuild array  must contains replace objects only")
			return
		}
		build, err2 := NewBuild(value)
		if err2 != nil {
			de = err2
			return
		}
		builds = append(builds, build)
	}); err != nil {
		return nil, err
	}
	if de != nil {
		return nil, de
	}
	return builds, nil
}

// NewBuild parse build data and return Build config object instance
func NewBuild(json []byte) (build *Build, err error) {
	build = &Build{}
	if err = jsonparser.ObjectEach(json, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
		switch strings.ToLower(string(key)) {
		case "from":
			if dataType != jsonparser.String {
				return fmt.Errorf("Builder.From expected string and take %s", value)
			}
			build.From = string(value)
		case "to":
			if dataType != jsonparser.String {
				return fmt.Errorf("Builder.To expected string and take %s", value)
			}
			build.To = string(value)
		case "layout":
			if dataType != jsonparser.String {
				return fmt.Errorf("Builder.Layout expected string and take %s", value)
			}
			build.Layout = string(value)
		case "template":
			if dataType != jsonparser.String {
				return fmt.Errorf("Builder.view expected string and take %s", value)
			}
			build.Template = string(value)
		case "properties":
			if dataType != jsonparser.Object {
				return fmt.Errorf("DataSet.Properties expected map and take %s", value)
			}
			if build.Properties, err = plainmap.JSONToPlainStringMap(value); err != nil {
				return err
			}
		case "comment":
			// ignore all comments
		default:
			return fmt.Errorf("config.NewReplace: Unknow key %s (value: %s)", key, value)
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return build, nil
}
