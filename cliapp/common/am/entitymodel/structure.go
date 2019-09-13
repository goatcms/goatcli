package entitymodel

import (
	"strings"

	"github.com/goatcms/goatcli/cliapp/common/naming"
	"github.com/goatcms/goatcore/varutil/goaterr"
)

// Structure contains fields and relations structure
type Structure struct {
	Path       []string
	FullName   Name
	Name       Name
	Entity     *Entity
	Fields     StructureFields
	Relations  StructureRelations
	Structures StructureChilds
}

// StructureFields struct represent fields set
type StructureFields struct {
	ByName  map[string]*Field
	ByType  map[string][]*Field
	Ordered []*Field
}

func newStructureFields() StructureFields {
	return StructureFields{
		ByName:  map[string]*Field{},
		ByType:  map[string][]*Field{},
		Ordered: []*Field{},
	}
}

// StructureRelations struct represent relations set
type StructureRelations struct {
	ByName  map[string]*Relation
	Ordered []*Relation
}

func newStructureRelations() StructureRelations {
	return StructureRelations{
		ByName:  map[string]*Relation{},
		Ordered: []*Relation{},
	}
}

// StructureChilds struct represent relations set
type StructureChilds struct {
	ByName  map[string]*Structure
	Ordered []*Structure
}

func newStructureChilds() StructureChilds {
	return StructureChilds{
		ByName:  map[string]*Structure{},
		Ordered: []*Structure{},
	}
}

// NewRootStructure create new Structure instance
func NewRootStructure(entity *Entity) (instance *Structure, err error) {
	return &Structure{
		Path:       nil,
		FullName:   Name{},
		Name:       Name{},
		Entity:     entity,
		Fields:     newStructureFields(),
		Relations:  newStructureRelations(),
		Structures: newStructureChilds(),
	}, nil
}

// NewStructure create new Structure instance
func NewStructure(plainName string, entity *Entity) (instance *Structure, err error) {
	var (
		name string
	)
	plainName = strings.TrimSpace(plainName)
	if plainName == "" {
		return nil, goaterr.Errorf("Structure name is required")
	}
	instance = &Structure{
		Entity:     entity,
		Fields:     newStructureFields(),
		Relations:  newStructureRelations(),
		Structures: newStructureChilds(),
	}
	instance.Path = strings.Split(plainName, ".")
	name = instance.Path[len(instance.Path)-1]
	instance.Path = instance.Path[:len(instance.Path)-1]
	for i := 0; i < len(instance.Path); i++ {
		instance.Path[i] = naming.ToCamelCaseUF(instance.Path[i])
	}
	if instance.Name, err = NewName(name); err != nil {
		return nil, err
	}
	if instance.FullName, err = NewName(plainName); err != nil {
		return nil, err
	}
	return instance, nil
}

// AddField add new field to structure
func (structure *Structure) AddField(field *Field) (err error) {
	if err = structure.preventDuplicateNames(field.Name.CamelCaseUF); err != nil {
		return err
	}
	field.Structure = structure
	structure.Fields.ByName[field.Name.CamelCaseUF] = field
	structure.Fields.ByType[field.Type] = append(structure.Fields.ByType[field.Type], field)
	structure.Fields.Ordered = append(structure.Fields.Ordered, field)
	return nil
}

// AddRelation add new relation to structure
func (structure *Structure) AddRelation(relation *Relation) (err error) {
	if err = structure.preventDuplicateNames(relation.Name.CamelCaseUF); err != nil {
		return err
	}
	relation.Structure = structure
	structure.Relations.ByName[relation.Name.CamelCaseUF] = relation
	structure.Relations.Ordered = append(structure.Relations.Ordered, relation)
	return nil
}

// NewStructure create new child structure
func (structure *Structure) instanceOfStructure(name string) (node *Structure, err error) {
	var (
		fullName   Name
		structName Name
		ok         bool
	)
	name = naming.ToCamelCaseUF(name)
	if node, ok = structure.Structures.ByName[name]; ok {
		return node, nil
	}
	if err = structure.preventDuplicateNames(name); err != nil {
		return nil, err
	}
	if fullName, err = NewName(structure.FullName.CamelCaseUF + name); err != nil {
		return nil, err
	}
	if structName, err = NewName(name); err != nil {
		return nil, err
	}
	node = &Structure{
		Path:       append(structure.Path, name),
		FullName:   fullName,
		Name:       structName,
		Entity:     structure.Entity,
		Fields:     newStructureFields(),
		Relations:  newStructureRelations(),
		Structures: newStructureChilds(),
	}
	structure.Structures.ByName[name] = node
	structure.Structures.Ordered = append(structure.Structures.Ordered, node)
	return node, nil
}

func (structure *Structure) preventDuplicateNames(name string) (err error) {
	if _, ok := structure.Fields.ByName[name]; ok {
		return goaterr.Errorf(duplicateStructureElementError, name)
	}
	if _, ok := structure.Relations.ByName[name]; ok {
		return goaterr.Errorf(duplicateStructureElementError, name)
	}
	if _, ok := structure.Structures.ByName[name]; ok {
		return goaterr.Errorf(duplicateStructureElementError, name)
	}
	return nil
}
