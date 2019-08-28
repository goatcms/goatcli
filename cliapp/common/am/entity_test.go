package am

import (
	"testing"

	"github.com/goatcms/goatcore/varutil/plainmap"
)

const (
	testNewEntityModelJSON = `
	{
	  "model": {
	    "user": {
	      "name": "user",
	      "plural": "users",
	      "admin_read_roles":   "admin",
	      "admin_insert_roles": "admin",
	      "admin_edit_roles":   "admin",
	      "admin_delete_roles": "admin",
	      "label":"email",
	      "fields": {
	        "0": {
	          "name": "firstname",
	          "system": "n",
	          "type": "string",
	          "unique": "n",
	          "required": "y"
	        },
	        "1": {
	          "name": "lastname",
	          "system": "n",
	          "type": "string",
	          "unique": "n",
	          "required": "y"
	        },
	        "2": {
	          "name": "email",
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
	      }
	    }
	  }
	}
`
)

func TestNewEntityModel(t *testing.T) {
	var (
		data        map[string]string
		entityModel *EntityModel
		err         error
	)
	t.Parallel()
	if data, err = plainmap.JSONToPlainStringMap([]byte(testNewEntityModelJSON)); err != nil {
		t.Error(err)
		return
	}

	if entityModel, err = NewEntityModel("model.user", data); err != nil {
		t.Error(err)
		return
	}
	if entityModel.Name.CamelCaseUF != "User" {
		t.Errorf("entityModel.Name.CamelCaseUF should be equals to 'User' and takes %s", entityModel.Name.CamelCaseUF)
		return
	}
	if entityModel.Singular.CamelCaseUF != "User" {
		t.Errorf("entityModel.Singular.CamelCaseUF should be equals to 'User' and takes %s", entityModel.Singular.CamelCaseUF)
		return
	}
	if entityModel.Plural.CamelCaseUF != "Users" {
		t.Errorf("entityModel.Plural.CamelCaseUF should be equals to 'Users' and takes %s", entityModel.Plural.CamelCaseUF)
		return
	}
	if len(entityModel.ACL.Admin.DeleteRoles) != 1 || entityModel.ACL.Admin.DeleteRoles[0] != "admin" {
		t.Errorf("expected single admin role for entityModel.ACL.Admin.DeleteRoles and takes %s", entityModel.ACL.Admin.DeleteRoles)
		return
	}
	if len(entityModel.ACL.Admin.EditRoles) != 1 || entityModel.ACL.Admin.EditRoles[0] != "admin" {
		t.Errorf("expected single admin role for entityModel.ACL.Admin.EditRoles and takes %s", entityModel.ACL.Admin.EditRoles)
		return
	}
	if len(entityModel.ACL.Admin.InsertRoles) != 1 || entityModel.ACL.Admin.InsertRoles[0] != "admin" {
		t.Errorf("expected single admin role for entityModel.ACL.Admin.InsertRoles and takes %s", entityModel.ACL.Admin.InsertRoles)
		return
	}
	if len(entityModel.ACL.Admin.ReadRoles) != 1 || entityModel.ACL.Admin.ReadRoles[0] != "admin" {
		t.Errorf("expected single admin role for entityModel.ACL.Admin.ReadRoles and takes %s", entityModel.ACL.Admin.ReadRoles)
		return
	}
	if entityModel.LabelField == nil || entityModel.LabelField.Name.CamelCaseUF != "Email" {
		t.Errorf("expected Email like label field")
		return
	}
}
