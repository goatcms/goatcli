package entitymodel

import (
	"fmt"
	"strings"
)

// Field struct represent single entity field
type Field struct {
	// Path     []string
	FullName  Name
	Name      Name
	Type      string
	Flags     FieldFlags
	Structure *Structure
}

// FieldFlags contains flags represent a field
type FieldFlags struct {
	System   bool
	Unique   bool
	Required bool
}

// NewField create new Field instance
func NewField(name string) (instance *Field, err error) {
	var index int
	instance = &Field{}
	if instance.FullName, err = NewName(name); err != nil {
		return nil, err
	}
	index = strings.LastIndex(name, ".")
	if index == -1 {
		instance.Name = instance.FullName
	} else {
		if instance.Name, err = NewName(instance.FullName.Plain[index:]); err != nil {
			return nil, err
		}
	}
	return instance, nil
}

// NewFieldFromPlainmap create new Field instance and load data from plainmap
func NewFieldFromPlainmap(prefix string, plainmap map[string]string) (instance *Field, err error) {
	if instance, err = NewField(plainmap[prefix+".name"]); err != nil {
		return nil, err
	}
	instance.Type = strings.ToLower(plainmap[prefix+".type"])
	if instance.Type == "" {
		return nil, fmt.Errorf("Name is required")
	}
	instance.Flags.System = plainmap[prefix+".system"] == "y"
	instance.Flags.Unique = plainmap[prefix+".unique"] == "y"
	instance.Flags.Required = plainmap[prefix+".required"] == "y"
	return instance, nil
}
