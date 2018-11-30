package config

import (
	"fmt"
	"strings"

	"github.com/buger/jsonparser"
)

// Dependency is configuration container for one dependency
type Dependency struct {
	Repo   string `json:"repo,omitempty"`
	Branch string `json:"branch,omitempty"`
	Rev    string `json:"rev,omitempty"`
	Dest   string `json:"dest,omitempty"`
}

// NewDependencies parse json and return dependencies array instance
func NewDependencies(json []byte) ([]*Dependency, error) {
	var (
		de   error
		err  error
		deps = []*Dependency{}
	)
	if _, err = jsonparser.ArrayEach(json, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		if err != nil || de != nil {
			return
		}
		if dataType != jsonparser.Object {
			de = fmt.Errorf("NewDependencies array must contains dependency objects only")
			return
		}
		d, err2 := NewDependency(value)
		if err2 != nil {
			de = err2
			return
		}
		deps = append(deps, d)
	}); err != nil {
		return nil, err
	}
	if de != nil {
		return nil, de
	}
	return deps, nil
}

// NewDependency parse module data and return dependency object instance
func NewDependency(json []byte) (*Dependency, error) {
	var err error
	d := &Dependency{}
	if err = jsonparser.ObjectEach(json, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
		switch strings.ToLower(string(key)) {
		case "repo":
			if dataType != jsonparser.String {
				return fmt.Errorf("expected string and take %s", value)
			}
			d.Repo = string(value)
		case "branch":
			if dataType != jsonparser.String {
				return fmt.Errorf("expected string and take %s", value)
			}
			d.Branch = string(value)
		case "rev":
			if dataType != jsonparser.String {
				return fmt.Errorf("expected string and take %s", value)
			}
			d.Rev = string(value)
		case "dest":
			if dataType != jsonparser.String {
				return fmt.Errorf("expected string and take %s", value)
			}
			d.Dest = string(value)
		case commentProperty:
			// ignore all comments
		default:
			return fmt.Errorf("config.NewModules: Unknow key %s (value: %s)", key, value)
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return d, nil
}
