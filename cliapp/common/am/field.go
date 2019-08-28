package am

import (
	"fmt"
	"strings"
)

// Field struct represent single entity field
type Field struct {
	// Path     []string
	FullName Name
	Name     Name
	Type     string
	Flags    FieldFlags
}

// FieldFlags contains flags represent a field
type FieldFlags struct {
	System   bool
	Unique   bool
	Required bool
}

// NewField create new Field instance
func NewField(prefix string, plainmap map[string]string) (instance *Field, err error) {
	var (
		index   int
		nameStr = plainmap[prefix+".name"]
	)
	instance = &Field{}
	if instance.FullName, err = NewName(nameStr); err != nil {
		return nil, err
	}
	index = strings.LastIndex(nameStr, ".")
	if index == -1 {
		instance.Name = instance.FullName
		// instance.Path = make([]string, 0)
	} else {
		if instance.Name, err = NewName(instance.FullName.Plain[index:]); err != nil {
			return nil, err
		}
		// instance.Path = strings.Split(nameStr[:index], ".")
		// for i := 0; i < len(instance.Path); i++ {
		// 	instance.Path[i] = naming.ToCamelCaseUF(instance.Path[i])
		// }
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
