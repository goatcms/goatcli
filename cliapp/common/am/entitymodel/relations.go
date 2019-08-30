package entitymodel

import (
	"github.com/goatcms/goatcore/varutil/plainmap"
)

// Relations struct represent relations set
type Relations map[string]*Relation

// NewRelations create new Relations instance
func NewRelations() Relations {
	return Relations{}
}

// NewRelationsFromPlainmap create new Relations instance and load data from plainmap
func NewRelationsFromPlainmap(baseKey string, data map[string]string) (relations Relations, err error) {
	var (
		relation *Relation
		key      string
	)
	relations = NewRelations()
	baseKey += "."
	for _, index := range plainmap.Keys(data, baseKey) {
		key = baseKey + index
		if relation, err = NewRelationFromPlainmap(key, data); err != nil {
			return nil, err
		}
		relations[relation.FullName.CamelCaseUF] = relation
	}
	return relations, nil
}
