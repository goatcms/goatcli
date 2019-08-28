package am

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

// EntityModel contains single entity model data
type EntityModel struct {
	Name       Name
	Singular   Name
	Plural     Name
	LabelField *Field
	ACL        EntityACL
	Fields     *FieldSet
	Relations  RelationSet
	Structure  *Structure
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

// NewEntityModel create new EntityModel instance
func NewEntityModel(prefix string, data map[string]string) (instance *EntityModel, err error) {
	var (
		labelField string
		ok         bool
		value      string
	)
	instance = &EntityModel{}
	if value = data[prefix+".name"]; value == "" {
		return nil, fmt.Errorf("%s is required", prefix+".name")
	}
	if instance.Name, err = NewName(value); err != nil {
		return nil, err
	}
	if value = data[prefix+".singular"]; value == "" {
		instance.Singular = instance.Name
	} else {
		if instance.Singular, err = NewName(value); err != nil {
			return nil, err
		}
	}
	if value = data[prefix+".plural"]; value == "" {
		instance.Plural = instance.Name
	} else {
		if instance.Plural, err = NewName(value); err != nil {
			return nil, err
		}
	}
	instance.Structure = &Structure{
		Fields:     &FieldSet{},
		Relations:  RelationSet{},
		Structures: map[string]*Structure{},
	}
	if err = loadStructure(prefix, data, instance); err != nil {
		return nil, err
	}
	if labelField = data[prefix+".label"]; labelField == "" {
		for _, field := range instance.Fields.ByName {
			instance.LabelField = field
			break
		}
	} else {
		labelFieldUF := naming.ToCamelCaseUF(labelField)
		if instance.LabelField, ok = instance.Fields.ByName[labelFieldUF]; !ok {
			return nil, fmt.Errorf("%s label field is not exist for %s entity", labelFieldUF, instance.Name)
		}
	}
	loadACL(prefix, data, instance)
	return instance, nil
}

func loadACL(prefix string, data map[string]string, e *EntityModel) {
	e.ACL.Admin.EditRoles = arrayValue(data, prefix+".admin_edit_roles", DefaultACLRoles)
	e.ACL.Admin.ReadRoles = arrayValue(data, prefix+".admin_read_roles", DefaultACLRoles)
	e.ACL.Admin.InsertRoles = arrayValue(data, prefix+".admin_insert_roles", DefaultACLRoles)
	e.ACL.Admin.DeleteRoles = arrayValue(data, prefix+".admin_read_roles", DefaultACLRoles)
}

func loadStructure(prefix string, data map[string]string, e *EntityModel) (err error) {
	var (
		structure *Structure
		path      string
	)
	// load
	e.Structure = &Structure{
		Structures: map[string]*Structure{},
	}
	if e.Fields, err = NewFieldSet(prefix+".fields", data); err != nil {
		return err
	}
	if e.Relations, err = NewRelationSet(prefix+".relations", data); err != nil {
		return err
	}
	// build structure - fields
	for _, field := range e.Structure.Fields.ByName {
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
	for _, relation := range e.Structure.Relations {
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
