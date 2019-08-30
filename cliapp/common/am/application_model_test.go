package am

import (
	"testing"

	"github.com/goatcms/goatcli/cliapp/common/am/entitymodel"
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
	      "label":"firstname",
	      "fields": {
	        "0": {
	          "name": "firstname",
	          "system": "n",
	          "type": "string",
	          "unique": "n",
	          "required": "y"
	        },
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
	        }
	      },
	      "relations": {
	        "0": {
	          "to": "user",
	          "name": "owner",
	          "required": "y",
	          "system": "y",
	          "unique": "n"
	        }
	      }
	    }
	  },
	  "dto": {
	    "user": {
	      "name": "user",
	      "plural": "users",
	      "fields": {
	        "0": {
	          "name": "firstname",
	          "system": "n",
	          "type": "string",
	          "unique": "n",
	          "required": "y"
	        },
	      }
	    }
	  },
	  "options": {
	    "user": {
	      "name": "user",
	      "plural": "users",
	      "fields": {
	        "0": {
	          "name": "firstname",
	          "system": "n",
	          "type": "string",
	          "unique": "n",
	          "required": "y"
	        },
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
		entities         entitymodel.Entities
		dtos             entitymodel.Entities
		options          entitymodel.Entities
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
	if options, err = applicationModel.Options(); err != nil {
		t.Error(err)
		return
	}
	if len(options) != 1 {
		t.Errorf("expected one options in options and take %v", options)
		return
	}
	if dtos, err = applicationModel.DTO(); err != nil {
		t.Error(err)
		return
	}
	if len(dtos) != 1 {
		t.Errorf("expected one dtos in dtos and take %v", dtos)
		return
	}
}
