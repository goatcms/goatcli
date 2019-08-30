package entitymodel

import (
	"github.com/goatcms/goatcore/varutil/plainmap"
)

// Entities contains entities
type Entities map[string]*Entity

// NewEntities create new Entities instance
func NewEntities() (instance Entities) {
	return map[string]*Entity{}
}

// NewEntitiesFromPlainmap create new Entities instance and load data from plainmap
func NewEntitiesFromPlainmap(prefix string, data map[string]string) (instance Entities, err error) {
	var entity *Entity
	instance = map[string]*Entity{}
	prefix += "."
	for _, key := range plainmap.Keys(data, prefix) {
		if entity, err = NewEntityFromPlainmap(prefix+key, data); err != nil {
			return nil, err
		}
		instance[entity.Name.CamelCaseUF] = entity
	}
	return instance, nil
}
