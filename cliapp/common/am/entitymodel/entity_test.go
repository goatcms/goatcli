package entitymodel

import (
	"testing"

	"github.com/goatcms/goatcore/varutil/plainmap"
)

const (
	testNewEntityJSON = `
	{
	  "model": {
	    "user": {
	      "name": "user",
	      "plural": "users",
	      "admin_read_roles":   "admin",
	      "admin_insert_roles": "admin",
	      "admin_edit_roles":   "admin",
	      "admin_delete_roles": "admin",
	      "label":"person.email",
	      "fields": {
	        "0": {
	          "name": "person.firstname",
	          "system": "n",
	          "type": "string",
	          "unique": "n",
	          "required": "y"
	        },
	        "1": {
	          "name": "person.lastname",
	          "system": "n",
	          "type": "string",
	          "unique": "n",
	          "required": "y"
	        },
	        "2": {
	          "name": "person.email",
	          "system": "n",
	          "type": "email",
	          "unique": "y",
	          "required": "y"
	        },
	        "3": {
	          "name": "password",
	          "system": "y",
	          "type": "password",
	          "unique": "n",
	          "required": "n"
	        },
	        "4": {
	          "name": "roles",
	          "system": "y",
	          "type": "string",
	          "unique": "n",
	          "required": "n"
	        },
	        "5": {
	          "name": "username",
	          "system": "n",
	          "type": "string",
	          "unique": "y",
	          "required": "y"
	        }
	      },
	      "relations": {
	        "0": {
	          "name": "owner",
	          "to": "user"
	        }
				}
	    }
	  }
	}
`
)

func TestNewEntity(t *testing.T) {
	var (
		data   map[string]string
		entity *Entity
		err    error
	)
	t.Parallel()
	if data, err = plainmap.JSONToPlainStringMap([]byte(testNewEntityJSON)); err != nil {
		t.Error(err)
		return
	}
	if entity, err = NewEntityFromPlainmap("model.user", data); err != nil {
		t.Error(err)
		return
	}
	if entity.Name.CamelCaseUF != "User" {
		t.Errorf("entity.Name.CamelCaseUF should be equals to 'User' and takes %s", entity.Name.CamelCaseUF)
		return
	}
	if entity.Singular.CamelCaseUF != "User" {
		t.Errorf("entity.Singular.CamelCaseUF should be equals to 'User' and takes %s", entity.Singular.CamelCaseUF)
		return
	}
	if entity.Plural.CamelCaseUF != "Users" {
		t.Errorf("entity.Plural.CamelCaseUF should be equals to 'Users' and takes %s", entity.Plural.CamelCaseUF)
		return
	}
	if entity.AllFields.ByName["Password"].Type != "password" {
		t.Errorf("Password must be a password type")
		return
	}
	if entity.Structure.Fields.ByName["Password"].Type != "password" {
		t.Errorf("Password (in root structure) must be a password type")
		return
	}
	if entity.AllFields.ByName["PersonFirstname"].Type != "string" {
		t.Errorf("PersonFirstname must be a string type")
		return
	}
	if entity.Structure.Structures["Person"].Fields.ByName["Firstname"].Type != "string" {
		t.Errorf("Firstname in Person structure must be a string type")
		return
	}
	if len(entity.ACL.Admin.DeleteRoles) != 1 || entity.ACL.Admin.DeleteRoles[0] != "admin" {
		t.Errorf("expected single admin role for entity.ACL.Admin.DeleteRoles and takes %s", entity.ACL.Admin.DeleteRoles)
		return
	}
	if len(entity.ACL.Admin.EditRoles) != 1 || entity.ACL.Admin.EditRoles[0] != "admin" {
		t.Errorf("expected single admin role for entity.ACL.Admin.EditRoles and takes %s", entity.ACL.Admin.EditRoles)
		return
	}
	if len(entity.ACL.Admin.InsertRoles) != 1 || entity.ACL.Admin.InsertRoles[0] != "admin" {
		t.Errorf("expected single admin role for entity.ACL.Admin.InsertRoles and takes %s", entity.ACL.Admin.InsertRoles)
		return
	}
	if len(entity.ACL.Admin.ReadRoles) != 1 || entity.ACL.Admin.ReadRoles[0] != "admin" {
		t.Errorf("expected single admin role for entity.ACL.Admin.ReadRoles and takes %s", entity.ACL.Admin.ReadRoles)
		return
	}
	if entity.LabelField == nil || entity.LabelField.Name.CamelCaseUF != "Email" {
		t.Errorf("expected Email like label field")
		return
	}
}
