package am

import (
	"github.com/goatcms/goatcore/varutil/plainmap"
)

// RelationSet struct represent relations set
type RelationSet map[string]*Relation

// NewRelationSet create new RelationSet instance
func NewRelationSet(baseKey string, data map[string]string) (relations RelationSet, err error) {
	var (
		relation *Relation
		key      string
	)
	relations = RelationSet{}
	baseKey += "."
	for _, index := range plainmap.Keys(data, baseKey) {
		key = baseKey + index
		if relation, err = NewRelation(key, data); err != nil {
			return nil, err
		}
		relations[relation.FullName.CamelCaseUF] = relation
	}
	return relations, nil
}
