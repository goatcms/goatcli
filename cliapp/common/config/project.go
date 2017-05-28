package config

import (
	"fmt"

	"github.com/buger/jsonparser"
)

// Project is configuration container for one project
type Project struct {
	ProjectPath string
	Modules     []*Module
}

// NewProject parse project data and return object
func NewProject(json []byte) (*Project, error) {
	var err error
	c := &Project{}
	if err := jsonparser.ObjectEach(json, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
		keys := string(key)
		if keys == "projectPath" {
			if dataType != jsonparser.String {
				return fmt.Errorf("expected string and take %s", value)
			}
			c.ProjectPath = string(value)
		}
		if keys == "modules" {
			if dataType != jsonparser.Array {
				return fmt.Errorf("expected array and take %s", value)
			}
			c.Modules, err = NewModules(value)
			if err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return c, nil
}
