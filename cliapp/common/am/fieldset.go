package am

import (
	"fmt"
	"strings"

	"github.com/goatcms/goatcli/cliapp/common/naming"
	"github.com/goatcms/goatcore/varutil/plainmap"
)

// FieldSet struct represent fields set
type FieldSet struct {
	ByName  map[string]*Field
	ByType  map[string][]*Field
	Structs map[string]*FieldSet
}

// NewFieldSet create new FieldSet instance
func NewFieldSet(baseKey string, data map[string]string) (root *FieldSet, err error) {
	var (
		structSet *FieldSet
		field     *Field
		key       string
		pos       int
		ok        bool
	)
	root = &FieldSet{
		ByName:  map[string]*Field{},
		ByType:  map[string][]*Field{},
		Structs: map[string]*FieldSet{},
	}
	baseKey += "."
	for _, index := range plainmap.Keys(data, baseKey) {
		key = baseKey + index
		if field, err = NewField(key, data); err != nil {
			return nil, err
		}
		pos = strings.LastIndex(field.FullName.Plain, ".")
		if pos != -1 {
			if structSet, err = root.structByPath(field.FullName.Plain[:pos]); err != nil {
				return nil, err
			}
			if _, ok = structSet.Structs[field.Name.CamelCaseUF]; ok {
				return nil, fmt.Errorf("%s field can not be node and left", field.FullName.Plain)
			}
			structSet.ByName[field.Name.CamelCaseUF] = field
			structSet.ByType[field.Type] = append(structSet.ByType[field.Type], field)
		}
		root.ByName[field.FullName.CamelCaseUF] = field
		root.ByType[field.Type] = append(root.ByType[field.Type], field)
	}
	return root, nil
}

func (fs *FieldSet) structByPath(path string) (current *FieldSet, err error) {
	current = fs
	for _, key := range strings.Split(path, ".") {
		key = naming.ToCamelCaseUF(key)
		if newSet, ok := current.Structs[key]; ok {
			current = newSet
		} else {
			if _, ok := current.ByName[key]; ok {
				return nil, fmt.Errorf("%s (%s) field can not be node and left", path, key)
			}
			newSet := &FieldSet{
				ByName:  map[string]*Field{},
				ByType:  map[string][]*Field{},
				Structs: map[string]*FieldSet{},
			}
			current.Structs[key] = newSet
			current = newSet
		}
	}
	return current, nil
}
