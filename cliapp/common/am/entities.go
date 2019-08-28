package am

import (
	"github.com/goatcms/goatcore/varutil/plainmap"
)

// EntitiesModel contains entities
type EntitiesModel map[string]*EntityModel

// NewEntitiesModel create new EntityModel instance
func NewEntitiesModel(prefix string, data map[string]string) (instance EntitiesModel, err error) {
	var entityModel *EntityModel
	instance = map[string]*EntityModel{}
	prefix += "."
	for _, key := range plainmap.Keys(data, prefix) {
		if entityModel, err = NewEntityModel(prefix+key, data); err != nil {
			return nil, err
		}
		instance[entityModel.Name.CamelCaseUF] = entityModel
	}
	return instance, nil
}
