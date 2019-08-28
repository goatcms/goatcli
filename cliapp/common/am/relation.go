package am

import (
	"fmt"
	"strings"
)

/*
"relations": {
	"0": {
		"model": "user",
		"name": "user",
		"required": "y",
		"system": "y",
		"unique": "n"
	}
}
*/

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
func NewRelation(prefix string, plainmap map[string]string) (instance *Relation, err error) {
	var (
		index   int
		nameStr = plainmap[prefix+".name"]
	)
	instance = &Relation{}
	if instance.FullName, err = NewName(nameStr); err != nil {
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
	instance.To = strings.ToLower(plainmap[prefix+".to"])
	if instance.To == "" {
		return nil, fmt.Errorf("Name is required")
	}
	instance.Flags.System = plainmap[prefix+".system"] == "y"
	instance.Flags.Unique = plainmap[prefix+".unique"] == "y"
	instance.Flags.Required = plainmap[prefix+".required"] == "y"
	return instance, nil
}
