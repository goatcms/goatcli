package entitymodel

import (
	"sort"
	"strings"

	"github.com/goatcms/goatcli/cliapp/common/naming"
	"github.com/goatcms/goatcore/varutil"
	"github.com/goatcms/goatcore/varutil/goaterr"
	"github.com/goatcms/goatcore/varutil/plainmap"
)

var (
	// DefaultACLRoles is default role for ACL
	DefaultACLRoles = []string{"admin"}
	pluralSuffixes  = []string{"ch", "sh", "s", "x", "z"}
)

// Entity contains single entity model data
type Entity struct {
	Name          Name
	Singular      Name
	Plural        Name
	ACL           EntityACL
	LabelField    *Field
	AllFields     EntityFields
	AllRelations  EntityRelations
	RootStructure *Structure
}

// EntityFields contains all entity fields
type EntityFields struct {
	ByFullName map[string]*Field
	ByType     map[string][]*Field
	Ordered    []*Field
}

// EntityRelations contains all entity relations
type EntityRelations struct {
	ByFullName map[string]*Relation
	Ordered    []*Relation
}

// EntityACL contains criteria to access entity
type EntityACL struct {
	Admin AdminEntityACL
}

// AdminEntityACL contains criteria to access admin
type AdminEntityACL struct {
	ReadRoles   []string
	InsertRoles []string
	EditRoles   []string
	DeleteRoles []string
}

// NewEntity create new Entity instance
func NewEntity(name, singular, plural string) (instance *Entity, err error) {
	instance = &Entity{
		AllFields: EntityFields{
			ByFullName: map[string]*Field{},
			ByType:     map[string][]*Field{},
			Ordered:    []*Field{},
		},
		AllRelations: EntityRelations{
			ByFullName: map[string]*Relation{},
			Ordered:    []*Relation{},
		},
	}
	if instance.RootStructure, err = NewRootStructure(instance); err != nil {
		return nil, err
	}
	if name == "" {
		return nil, goaterr.Errorf("Entity name is required")
	}
	if instance.Name, err = NewName(name); err != nil {
		return nil, err
	}
	if singular == "" {
		instance.Singular = instance.Name
	} else {
		if instance.Singular, err = NewName(singular); err != nil {
			return nil, err
		}
	}
	if plural == "" {
		if varutil.HasOneSuffix(name, pluralSuffixes) {
			if instance.Plural, err = NewName(name + "es"); err != nil {
				return nil, err
			}
		} else {
			if instance.Plural, err = NewName(name + "s"); err != nil {
				return nil, err
			}
		}
	} else {
		if instance.Plural, err = NewName(plural); err != nil {
			return nil, err
		}
	}
	return instance, nil
}

// NewEntityFromPlainmap create new Entity instance and load data from plainmap
func NewEntityFromPlainmap(prefix string, data map[string]string) (instance *Entity, err error) {
	if instance, err = NewEntity(
		data[prefix+".name"],
		data[prefix+".singular"],
		data[prefix+".plural"]); err != nil {
		return nil, err
	}
	if err = instance.loadFields(prefix, data); err != nil {
		return nil, err
	}
	if err = instance.loadRelations(prefix, data); err != nil {
		return nil, err
	}
	if err = instance.loadStructure(prefix, data); err != nil {
		return nil, err
	}
	loadACL(prefix, data, instance)
	return instance, nil
}

func loadACL(prefix string, data map[string]string, e *Entity) {
	e.ACL.Admin.EditRoles = arrayValue(data, prefix+".admin_edit_roles", DefaultACLRoles)
	e.ACL.Admin.ReadRoles = arrayValue(data, prefix+".admin_read_roles", DefaultACLRoles)
	e.ACL.Admin.InsertRoles = arrayValue(data, prefix+".admin_insert_roles", DefaultACLRoles)
	e.ACL.Admin.DeleteRoles = arrayValue(data, prefix+".admin_read_roles", DefaultACLRoles)
}

func (e *Entity) loadFields(prefix string, data map[string]string) (err error) {
	var (
		field          *Field
		key            string
		ok             bool
		labelFieldName string
	)
	fieldsPrefix := prefix + ".fields."
	for _, index := range plainmap.Keys(data, fieldsPrefix) {
		key = fieldsPrefix + index
		if field, err = NewFieldFromPlainmap(key, data); err != nil {
			return err
		}
		if err = e.addField(field); err != nil {
			return err
		}
	}
	// add ID field by default
	if _, ok = e.AllFields.ByFullName["ID"]; !ok {
		var field *Field
		if field, err = NewField("ID"); err != nil {
			return err
		}
		field.Type = "id"
		field.Flags = FieldFlags{
			System:   true,
			Unique:   true,
			Required: false,
		}
		if err = e.addField(field); err != nil {
			return err
		}
	}
	// set label field
	if labelFieldName = data[prefix+".label"]; labelFieldName == "" {
		if e.LabelField, ok = e.AllFields.ByFullName["ID"]; !ok {
			for _, field := range e.AllFields.ByFullName {
				e.LabelField = field
				break
			}
		}
	} else {
		labelFieldUF := naming.ToCamelCaseUF(labelFieldName)
		if e.LabelField, ok = e.AllFields.ByFullName[labelFieldUF]; !ok {
			return goaterr.Errorf("%s label field is not exist for %s entity", labelFieldName, e.Name)
		}
	}
	ordered := e.AllFields.Ordered
	sort.Slice(ordered, func(i, j int) bool {
		if ordered[i].FullName.CamelCaseUF == "ID" {
			return true
		}
		if ordered[j].FullName.CamelCaseUF == "ID" {
			return false
		}
		return ordered[i].FullName.CamelCaseUF < ordered[j].FullName.CamelCaseUF
	})
	return nil
}

func (e *Entity) loadRelations(prefix string, data map[string]string) (err error) {
	var (
		relation *Relation
		key      string
	)
	relationPrefix := prefix + ".relations."
	for _, index := range plainmap.Keys(data, relationPrefix) {
		key = relationPrefix + index
		if relation, err = NewRelationFromPlainmap(key, data); err != nil {
			return err
		}
		if e.addRelation(relation); err != nil {
			return err
		}
	}
	ordered := e.AllRelations.Ordered
	sort.Slice(ordered, func(i, j int) bool {
		return ordered[i].FullName.CamelCaseUF < ordered[j].FullName.CamelCaseUF
	})
	return nil
}

func (e *Entity) loadStructure(prefix string, data map[string]string) (err error) {
	var (
		structure *Structure
		path      string
	)
	// build structure - fields
	for _, field := range e.AllFields.Ordered {
		path, _ = splitName(field.FullName.Plain)
		if structure, err = e.StructureByPath(path); err != nil {
			return err
		}
		if err = structure.AddField(field); err != nil {
			return err
		}
	}
	// build structure - relations
	for _, relation := range e.AllRelations.Ordered {
		path, _ = splitName(relation.FullName.Plain)
		if structure, err = e.StructureByPath(path); err != nil {
			return err
		}
		if err = structure.AddRelation(relation); err != nil {
			return err
		}
	}
	return nil
}

// addField add new field to entity
func (e *Entity) addField(field *Field) (err error) {
	if err = e.preventDuplicate(field.FullName.CamelCaseUF); err != nil {
		return err
	}
	e.AllFields.ByFullName[field.FullName.CamelCaseUF] = field
	e.AllFields.ByType[field.Type] = append(e.AllFields.ByType[field.Type], field)
	e.AllFields.Ordered = append(e.AllFields.Ordered, field)
	return nil
}

// addField add new field to entity
func (e *Entity) addRelation(relation *Relation) (err error) {
	if err = e.preventDuplicate(relation.FullName.CamelCaseUF); err != nil {
		return err
	}
	e.AllRelations.ByFullName[relation.FullName.CamelCaseUF] = relation
	e.AllRelations.Ordered = append(e.AllRelations.Ordered, relation)
	return nil
}

// addField add new field to entity
func (e *Entity) preventDuplicate(nameUF string) (err error) {
	if _, ok := e.AllFields.ByFullName[nameUF]; ok {
		return goaterr.Errorf(duplicateStructureElementError, nameUF)
	}
	if _, ok := e.AllRelations.ByFullName[nameUF]; ok {
		return goaterr.Errorf(duplicateStructureElementError, nameUF)
	}
	return nil
}

// StructureByPath return child structrure by path
func (e *Entity) StructureByPath(path string) (current *Structure, err error) {
	current = e.RootStructure
	if path == "" {
		return current, nil
	}
	for _, key := range strings.Split(path, ".") {
		if current, err = current.instanceOfStructure(key); err != nil {
			return nil, err
		}
	}
	return current, nil
}
