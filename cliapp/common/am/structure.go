package am

import (
	"fmt"
	"strings"

	"github.com/goatcms/goatcli/cliapp/common/naming"
)

// Structure contains fields and relations structure
type Structure struct {
	Path       []string
	Fields     *FieldSet
	Relations  RelationSet
	Structures map[string]*Structure
}

// NewStructure create new Structure instance
func NewStructure() (instance *Structure) {
	return &Structure{
		Path:       []string{},
		Fields:     &FieldSet{},
		Relations:  RelationSet{},
		Structures: map[string]*Structure{},
	}
}

func (structure *Structure) ByPath(path string) (current *Structure, err error) {
	current = structure
	for _, key := range strings.Split(path, ".") {
		key = naming.ToCamelCaseUF(key)
		if newSet, ok := current.Structures[key]; ok {
			current = newSet
		} else {
			if current, err = current.NewStructure(key); ok {
				return nil, err
			}
		}
	}
	return current, nil
}

func (structure *Structure) AddField(field *Field) (err error) {
	if _, ok := structure.Relations[field.Name.CamelCaseUF]; ok {
		return fmt.Errorf("%s contains field and relation with the same name", field.FullName.CamelCaseUF)
	}
	if _, ok := structure.Structures[field.Name.CamelCaseUF]; ok {
		return fmt.Errorf("%s contains field and structure with the same name", field.FullName.CamelCaseUF)
	}
	structure.Fields.ByName[field.Name.CamelCaseUF] = field
	typeset := structure.Fields.ByType[field.Name.CamelCaseUF]
	structure.Fields.ByType[field.Name.CamelCaseUF] = append(typeset, field)
	return nil
}

func (structure *Structure) AddRelation(relation *Relation) (err error) {
	if _, ok := structure.Fields.ByName[relation.Name.CamelCaseUF]; ok {
		return fmt.Errorf("%s contains field and relation with the same name", relation.FullName.CamelCaseUF)
	}
	if _, ok := structure.Structures[relation.Name.CamelCaseUF]; ok {
		return fmt.Errorf("%s contains relation and structure with the same name", relation.FullName.CamelCaseUF)
	}
	structure.Relations[relation.Name.CamelCaseUF] = relation
	return nil
}

func (structure *Structure) NewStructure(name string) (node *Structure, err error) {
	name = naming.ToCamelCaseUF(name)
	if _, ok := structure.Fields.ByName[name]; ok {
		return nil, fmt.Errorf("%s contains field and structure with the same name", name)
	}
	if _, ok := structure.Relations[name]; ok {
		return nil, fmt.Errorf("%s contains structure and relation with the same name", name)
	}
	node = &Structure{
		Path:       append(structure.Path, name),
		Fields:     &FieldSet{},
		Relations:  RelationSet{},
		Structures: map[string]*Structure{},
	}
	structure.Structures[name] = node
	return node, nil
}
