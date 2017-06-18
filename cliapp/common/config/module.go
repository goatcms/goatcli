package config

import (
	"fmt"
	"strings"

	"github.com/buger/jsonparser"
	"github.com/goatcms/goatcli/cliapp/common"
)

// Module is configuration container for one module
type Module struct {
	SourceURL string
	SourceRev string
	SourceDir string
	TestURL   string
	TestRev   string
	TestDir   string
	Replaces  []*Replace
}

// NewModules parse json and return Modules array instance
func NewModules(json []byte, si common.StringInjector) ([]*Module, error) {
	var de error = nil
	modules := []*Module{}
	if _, err := jsonparser.ArrayEach(json, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		if err != nil || de != nil {
			return
		}
		if dataType != jsonparser.Object {
			de = fmt.Errorf("NewModules array must contains replace objects only")
			return
		}
		m, err2 := NewModule(value, si)
		if err2 != nil {
			de = err2
			return
		}
		modules = append(modules, m)
	}); err != nil {
		return nil, err
	}
	if de != nil {
		return nil, de
	}
	return modules, nil
}

// NewModule parse module data and return module object instance
func NewModule(json []byte, si common.StringInjector) (*Module, error) {
	var err error
	c := &Module{}
	if err = jsonparser.ObjectEach(json, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
		switch strings.ToLower(string(key)) {
		case "srcclone":
			if dataType != jsonparser.String {
				return fmt.Errorf("expected string and take %s", value)
			}
			c.SourceURL = string(value)
		case "srcrev":
			if dataType != jsonparser.String {
				return fmt.Errorf("expected string and take %s", value)
			}
			c.SourceRev = string(value)
		case "srcdir":
			if dataType != jsonparser.String {
				return fmt.Errorf("expected string and take %s", value)
			}
			c.SourceDir = string(value)

		case "testclone":
			if dataType != jsonparser.String {
				return fmt.Errorf("expected string and take %s", value)
			}
			c.TestURL = string(value)
		case "testrev":
			if dataType != jsonparser.String {
				return fmt.Errorf("expected string and take %s", value)
			}
			c.TestRev = string(value)
		case "testdir":
			if dataType != jsonparser.String {
				return fmt.Errorf("expected string and take %s", value)
			}
			c.TestDir = string(value)
		case "replaces":
			if dataType != jsonparser.Array {
				return fmt.Errorf("expected replaces array and take %s", value)
			}
			c.Replaces, err = NewReplaces(value, si)
			if err != nil {
				return err
			}
		case "comment":
			// ignore all comments
		default:
			return fmt.Errorf("config.NewModules: Unknow key %s (value: %s)", key, value)
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return c, nil
}
