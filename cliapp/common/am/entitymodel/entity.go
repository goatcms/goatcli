package entitymodel

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/goatcms/goatcli/cliapp/common/naming"
)

var (
	// DefaultACLRoles is default role for ACL
	DefaultACLRoles = []string{"admin"}
)

// Entity contains single entity model data
type Entity struct {
	Name         Name
	Singular     Name
	Plural       Name
	LabelField   *Field
	ACL          EntityACL
	AllFields    *Fields
	AllRelations Relations
	Structure    *Structure
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
		AllFields:    NewFields(),
		AllRelations: NewRelations(),
		Structure:    NewStructure(),
	}
	if name == "" {
		return nil, fmt.Errorf("Entity name is required")
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
		instance.Plural = instance.Name
	} else {
		if instance.Plural, err = NewName(plural); err != nil {
			return nil, err
		}
	}
	return instance, nil
}

// NewEntityFromPlainmap create new Entity instance and load data from plainmap
func NewEntityFromPlainmap(prefix string, data map[string]string) (instance *Entity, err error) {
	var (
		labelField string
		ok         bool
	)
	if instance, err = NewEntity(
		data[prefix+".name"],
		data[prefix+".singular"],
		data[prefix+".plural"]); err != nil {
		return nil, err
	}
	if err = loadStructure(prefix, data, instance); err != nil {
		return nil, err
	}
	if labelField = data[prefix+".label"]; labelField == "" {
		for _, field := range instance.AllFields.ByName {
			instance.LabelField = field
			break
		}
	} else {
		labelFieldUF := naming.ToCamelCaseUF(labelField)
		if instance.LabelField, ok = instance.AllFields.ByName[labelFieldUF]; !ok {
			return nil, fmt.Errorf("%s label field is not exist for %s entity", labelFieldUF, instance.Name)
		}
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

func loadStructure(prefix string, data map[string]string, e *Entity) (err error) {
	var (
		structure *Structure
		path      string
	)
	// load
	e.Structure = NewStructure()
	if e.AllFields, err = NewFieldsFromPlainmap(prefix+".fields", data); err != nil {
		return err
	}
	if e.AllRelations, err = NewRelationsFromPlainmap(prefix+".relations", data); err != nil {
		return err
	}
	// build structure - fields
	for _, field := range e.AllFields.ByName {
		structure = e.Structure
		if path = pathFromPlainName(field.FullName.Plain); path != "" {
			if structure, err = e.Structure.ByPath(path); err != nil {
				return err
			}
		}
		if err = structure.AddField(field); err != nil {
			return err
		}
	}
	// build structure - relations
	for _, relation := range e.AllRelations {
		structure = e.Structure
		if path = pathFromPlainName(relation.FullName.Plain); path != "" {
			if structure, err = e.Structure.ByPath(path); err != nil {
				return err
			}
		}
		if err = structure.AddRelation(relation); err != nil {
			return err
		}
	}
	return nil
}

func arrayValue(data map[string]string, key string, defaultValue []string) (v []string) {
	var s string
	if s = data[key]; s == "" {
		return defaultValue
	}
	rx := regexp.MustCompile(`\s`)
	s = rx.ReplaceAllString(s, "")
	return strings.Split(s, ",")
}

func pathFromPlainName(name string) string {
	if pos := strings.LastIndex(name, "."); pos != -1 {
		return name[:pos]
	}
	return name
}
