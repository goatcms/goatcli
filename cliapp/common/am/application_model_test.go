package am

import (
	"testing"

	"github.com/goatcms/goatcore/varutil/plainmap"
)

const (
	testNewApplicationModelJSON = `
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
	    },
	    "session": {
	      "name": "session",
	      "fields": {
	        "0": {
	          "name": "secret",
	          "required": "y",
	          "system": "n",
	          "type": "string",
	          "unique": "y"
	        },
	        "1": {
	          "name": "expires",
	          "required": "y",
	          "system": "n",
	          "type": "int",
	          "unique": "n"
	        }
	      },
	      "relations": {
	        "0": {
	          "model": "user",
	          "name": "user",
	          "required": "y",
	          "system": "y",
	          "unique": "n"
	        }
	      }
	    }
	  }
	}
`
)

func TestNewApplicationModel(t *testing.T) {
	var (
		data             map[string]string
		applicationModel *ApplicationModel
		entities         EntitiesModel
		err              error
	)
	t.Parallel()
	if data, err = plainmap.JSONToPlainStringMap([]byte(testNewApplicationModelJSON)); err != nil {
		t.Error(err)
		return
	}
	applicationModel = NewApplicationModel(data)
	if entities, err = applicationModel.Entities(); err != nil {
		t.Error(err)
		return
	}
	if len(entities) != 2 {
		t.Errorf("expected two entities in entitiesModel and take %v", entities)
		return
	}
}
