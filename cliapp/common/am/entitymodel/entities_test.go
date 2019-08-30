package entitymodel

import (
	"testing"

	"github.com/goatcms/goatcore/varutil/plainmap"
)

const (
	testNewEntitiesJSON = `
	{
	  "model": {
	    "user": {
	      "name": "user",
	      "plural": "users",
	      "label":"person.firstname",
	      "fields": {
	        "0": {
	          "name": "person.firstname",
	          "system": "n",
	          "type": "string",
	          "unique": "n",
	          "required": "y"
	        },
	      },
	      "relations": {
	        "0": {
	          "name": "owner",
	          "to": "user"
	        }
				}
	    },
	    "car": {
	      "name": "car",
	      "plural": "cars",
	      "label":"model",
	      "fields": {
	        "0": {
	          "name": "model",
	          "system": "n",
	          "type": "string",
	          "unique": "n",
	          "required": "y"
	        },
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

func TestNewEntities(t *testing.T) {
	var (
		data     map[string]string
		entities Entities
		err      error
	)
	t.Parallel()
	if data, err = plainmap.JSONToPlainStringMap([]byte(testNewEntitiesJSON)); err != nil {
		t.Error(err)
		return
	}
	if entities, err = NewEntitiesFromPlainmap("model", data); err != nil {
		t.Error(err)
		return
	}
	if len(entities) != 2 {
		t.Errorf("expected two entities and take %v", len(entities))
		return
	}
}
