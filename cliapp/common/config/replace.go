package config

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/buger/jsonparser"
	"github.com/goatcms/goatcli/cliapp/common"
)

// Replace is configuration container for one replace description
type Replace struct {
	From    *regexp.Regexp
	To      string
	Pattern *regexp.Regexp
}

// NewReplaces parse json and return Replace array instance
func NewReplaces(json []byte, si common.StringInjector) ([]*Replace, error) {
	var de error = nil
	replaces := []*Replace{}
	if _, err := jsonparser.ArrayEach(json, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		if err != nil || de != nil {
			return
		}
		if dataType != jsonparser.Object {
			de = fmt.Errorf("NewReplaces array  must contains replace objects only")
			return
		}
		replace, err2 := NewReplace(value, si)
		if err2 != nil {
			de = err2
			return
		}
		replaces = append(replaces, replace)
	}); err != nil {
		return nil, err
	}
	if de != nil {
		return nil, de
	}
	return replaces, nil
}

// NewReplace parse replace data and return Replace object instance
func NewReplace(json []byte, si common.StringInjector) (*Replace, error) {
	var err error
	r := &Replace{}
	if err = jsonparser.ObjectEach(json, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
		switch strings.ToLower(string(key)) {
		case "from":
			if dataType != jsonparser.String {
				return fmt.Errorf("Replace.From expected string and take %s", value)
			}
			if r.From, err = regexp.Compile(string(value)); err != nil {
				return err
			}
		case "to":
			if dataType != jsonparser.String {
				return fmt.Errorf(" Replace.To expected string and take %s", value)
			}
			if r.To, err = si.InjectToString(string(value)); err != nil {
				return err
			}
		case "pattern":
			if dataType != jsonparser.String {
				return fmt.Errorf("Replace.Pattern must be a string")
			}
			if r.Pattern, err = regexp.Compile(string(value)); err != nil {
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
	return r, nil
}
