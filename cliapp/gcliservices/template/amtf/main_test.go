package amtf

import (
	"github.com/goatcms/goatcli/cliapp/common/am"
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
	      "label":"person.firstname",
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
	          "to": "user",
	          "name": "family.parent",
	          "required": "y",
	          "system": "y",
	          "unique": "n"
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

func testApplicationModel() (applicationModel *am.ApplicationModel, err error) {
	var data map[string]string
	if data, err = plainmap.JSONToPlainStringMap([]byte(testNewApplicationModelJSON)); err != nil {
		return nil, err
	}
	return am.NewApplicationModel(data), nil
}
