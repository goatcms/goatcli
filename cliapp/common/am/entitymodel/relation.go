package entitymodel

import (
	"fmt"
	"strings"
)

// Relation struct represent single relation to other entity
type Relation struct {
	FullName Name
	Name     Name
	To       string
	Flags    FieldFlags
}

// RelationFlags contains flags represent a relation flags
type RelationFlags struct {
	System   bool
	Unique   bool
	Required bool
}

// NewRelation create new Relation instance
func NewRelation(name, to string) (instance *Relation, err error) {
	var (
		index int
	)
	if name == "" {
		return nil, fmt.Errorf("Relation: Name is required")
	}
	if to == "" {
		return nil, fmt.Errorf("Relation: To field is required")
	}
	instance = &Relation{}
	if instance.FullName, err = NewName(name); err != nil {
		return nil, err
	}
	index = strings.LastIndex(instance.FullName.Plain, ".")
	if index == -1 {
		instance.Name = instance.FullName
	} else {
		if instance.Name, err = NewName(instance.FullName.Plain[index:]); err != nil {
			return nil, err
		}
	}
	instance.To = strings.ToLower(to)
	return instance, nil
}

// NewRelationFromPlainmap create new Relation instance and load data from plainmap
func NewRelationFromPlainmap(prefix string, plainmap map[string]string) (instance *Relation, err error) {
	if instance, err = NewRelation(plainmap[prefix+".name"], plainmap[prefix+".to"]); err != nil {
		return nil, err
	}
	instance.Flags.System = plainmap[prefix+".system"] == "y"
	instance.Flags.Unique = plainmap[prefix+".unique"] == "y"
	instance.Flags.Required = plainmap[prefix+".required"] == "y"
	return instance, nil
}
