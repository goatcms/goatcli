package config

import (
	"fmt"
	"strings"

	"github.com/buger/jsonparser"
)

// Replace is configuration container for one replace description
type Replace struct {
	From   string
	To     string
	Suffix []string
}

// NewReplaces parse json and return Replace array instance
func NewReplaces(json []byte) ([]*Replace, error) {
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
		replace, err2 := NewReplace(value)
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
func NewReplace(json []byte) (*Replace, error) {
	var err error
	r := &Replace{}
	if err = jsonparser.ObjectEach(json, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
		switch strings.ToLower(string(key)) {
		case "from":
			if dataType != jsonparser.String {
				return fmt.Errorf("expected string and take %s", value)
			}
			r.From = string(value)
		case "to":
			if dataType != jsonparser.String {
				return fmt.Errorf("expected string and take %s", value)
			}
			r.To = string(value)
		case "suffix":
			if dataType == jsonparser.String {
				r.Suffix = []string{string(value)}
			} else if dataType == jsonparser.Array {
				r.Suffix, err = parseStringArray(value)
				if err != nil {
					return err
				}
			} else {
				return fmt.Errorf("Incorrect Replace.Suffix type (allow strings array or single string value)")
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

func parseStringArray(json []byte) ([]string, error) {
	var de error
	sr := []string{}
	if _, err := jsonparser.ArrayEach(json, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		if err != nil || de != nil {
			return
		}
		if dataType != jsonparser.String {
			de = fmt.Errorf("parseStringArray support only string array")
			return
		}
		sr = append(sr, string(value))
	}); err != nil {
		return nil, err
	}
	if de != nil {
		return nil, de
	}
	return sr, nil
}
