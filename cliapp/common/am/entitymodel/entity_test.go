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
	        },
	        "6": {
	          "name": "access",
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
	if _, ok := entity.AllFields.ByFullName["ID"]; !ok {
		t.Errorf("Should add a ID field if it is undefined")
		return
	}
	if entity.AllFields.ByFullName["ID"].Type != "id" {
		t.Errorf("default ID type is 'id'")
		return
	}
	if entity.AllFields.ByFullName["Password"].Type != "password" {
		t.Errorf("Password must be a password type")
		return
	}
	if entity.RootStructure.Fields.ByName["Password"].Type != "password" {
		t.Errorf("Password (in root structure) must be a password type")
		return
	}
	if entity.AllFields.ByFullName["PersonFirstname"].Type != "string" {
		t.Errorf("PersonFirstname must be a string type")
		return
	}
	if len(entity.RootStructure.Structures.Ordered) != 1 {
		t.Errorf("Expected one sub sructure")
		return
	}
	if _, ok := entity.RootStructure.Structures.ByName["Person"]; !ok {
		t.Errorf("Person structure should be defined")
		return
	}
	personStructure := entity.RootStructure.Structures.ByName["Person"]
	if personStructure.FullName.CamelCaseUF != "Person" {
		t.Errorf("Expected FullName equals to 'Person' and get %s", personStructure.FullName.CamelCaseUF)
		return
	}
	if _, ok := personStructure.Fields.ByName["Firstname"]; !ok {
		t.Errorf("Person structure should have Firstname field")
		return
	}
	if personStructure.Fields.ByName["Firstname"].Type != "string" {
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
	if entity.LabelField == nil {
		t.Errorf("entity must defined Label field")
		return
	}
	if entity.LabelField.Name.CamelCaseUF != "Email" {
		t.Errorf("expected Email like label field and take %s", entity.LabelField.Name.CamelCaseUF)
		return
	}
}

func TestNewEntityFieldsOrder(t *testing.T) {
	var (
		data   map[string]string
		entity *Entity
		err    error
		take   string
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
	take = entity.AllFields.Ordered[0].FullName.CamelCaseUF
	if take != "ID" {
		t.Errorf("Expected id at first field and take %v", take)
		return
	}
	take = entity.AllFields.Ordered[1].FullName.CamelCaseUF
	if take != "Access" {
		t.Errorf("Expected Access at second field and take %v", take)
		return
	}
	take = entity.AllFields.Ordered[2].FullName.CamelCaseUF
	if take != "Password" {
		t.Errorf("Expected password at third field and take %v", take)
		return
	}
	take = entity.AllFields.Ordered[3].FullName.CamelCaseUF
	if take != "PersonEmail" {
		t.Errorf("Expected PersonEmail at fourth field and take %v", take)
		return
	}
	take = entity.AllFields.Ordered[4].FullName.CamelCaseUF
	if take != "PersonFirstname" {
		t.Errorf("Expected PersonFirstname at fifth field and take %v", take)
		return
	}
	take = entity.AllFields.Ordered[5].FullName.CamelCaseUF
	if take != "PersonLastname" {
		t.Errorf("Expected PersonLastname at sixth field and take %v", take)
		return
	}
	take = entity.AllFields.Ordered[6].FullName.CamelCaseUF
	if take != "Roles" {
		t.Errorf("Expected Roles at seventh field and take %v", take)
		return
	}
	take = entity.AllFields.Ordered[7].FullName.CamelCaseUF
	if take != "Username" {
		t.Errorf("Expected Username at eighth field and take %v", take)
		return
	}
}

func TestPersonStructOrder(t *testing.T) {
	var (
		data   map[string]string
		entity *Entity
		err    error
		take   string
		person *Structure
		ok     bool
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
	if person, ok = entity.RootStructure.Structures.ByName["Person"]; !ok {
		t.Errorf("Expected person structure")
		return
	}
	take = person.Fields.Ordered[0].Name.CamelCaseUF
	if take != "Email" {
		t.Errorf("Expected Email at first field and take %v", take)
		return
	}
	take = person.Fields.Ordered[1].Name.CamelCaseUF
	if take != "Firstname" {
		t.Errorf("Expected Firstname at second field and take %v", take)
		return
	}
	take = person.Fields.Ordered[2].Name.CamelCaseUF
	if take != "Lastname" {
		t.Errorf("Expected Lastname at third field and take %v", take)
		return
	}
}

func TestRootStructOrder(t *testing.T) {
	var (
		data   map[string]string
		entity *Entity
		err    error
		take   string
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
	ordered := entity.RootStructure.Fields.Ordered
	take = ordered[0].FullName.CamelCaseUF
	if take != "ID" {
		t.Errorf("Expected id at first field and take %v", take)
		return
	}
	take = ordered[1].FullName.CamelCaseUF
	if take != "Access" {
		t.Errorf("Expected Access at second field and take %v", take)
		return
	}
	take = ordered[2].FullName.CamelCaseUF
	if take != "Password" {
		t.Errorf("Expected password at third field and take %v", take)
		return
	}
	take = ordered[3].FullName.CamelCaseUF
	if take != "Roles" {
		t.Errorf("Expected Roles at fourth and take %v", take)
		return
	}
	take = ordered[4].FullName.CamelCaseUF
	if take != "Username" {
		t.Errorf("Expected Username at fifth and take %v", take)
		return
	}
}

func TestNewEntityStructure(t *testing.T) {
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
	if entity.RootStructure.Entity != entity {
		t.Errorf("main structure should have entity handler")
		return
	}
	if entity.RootStructure.Structures.ByName["Person"].Entity != entity {
		t.Errorf("evry child structure should have entity handler")
		return
	}
}
