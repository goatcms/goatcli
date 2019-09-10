package entitymodel

import (
	"fmt"
	"strings"

	"github.com/goatcms/goatcli/cliapp/common/naming"
)

const (
	duplicateStructureElementError = "Duplicated structure elemenet %s"
)

// Structure contains fields and relations structure
type Structure struct {
	Entity     *Entity
	Fields     *Fields
	Relations  Relations
	Path       []string
	Structures map[string]*Structure
}

// NewStructure create new Structure instance
func NewStructure(entity *Entity) (instance *Structure) {
	return &Structure{
		Entity:     entity,
		Fields:     NewFields(),
		Relations:  NewRelations(),
		Path:       []string{},
		Structures: map[string]*Structure{},
	}
}

// ByPath return child structrure by path
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

// AddField add new field to structure
func (structure *Structure) AddField(field *Field) (err error) {
	if err = structure.preventDuplicateNames(field.Name.CamelCaseUF); err != nil {
		return err
	}
	field.Structure = structure
	structure.Fields.ByName[field.Name.CamelCaseUF] = field
	typeset := structure.Fields.ByType[field.Name.CamelCaseUF]
	structure.Fields.ByType[field.Name.CamelCaseUF] = append(typeset, field)
	return nil
}

// AddRelation add new relation to structure
func (structure *Structure) AddRelation(relation *Relation) (err error) {
	if err = structure.preventDuplicateNames(relation.Name.CamelCaseUF); err != nil {
		return err
	}
	relation.Structure = structure
	structure.Relations[relation.Name.CamelCaseUF] = relation
	return nil
}

// NewStructure create new child structure
func (structure *Structure) NewStructure(name string) (node *Structure, err error) {
	name = naming.ToCamelCaseUF(name)
	if err = structure.preventDuplicateNames(name); err != nil {
		return nil, err
	}
	node = NewStructure(structure.Entity)
	node.Path = append(structure.Path, name)
	structure.Structures[name] = node
	return node, nil
}

func (structure *Structure) preventDuplicateNames(name string) (err error) {
	if _, ok := structure.Fields.ByName[name]; ok {
		return fmt.Errorf(duplicateStructureElementError, name)
	}
	if _, ok := structure.Relations[name]; ok {
		return fmt.Errorf(duplicateStructureElementError, name)
	}
	if _, ok := structure.Structures[name]; ok {
		return fmt.Errorf(duplicateStructureElementError, name)
	}
	return nil
}
