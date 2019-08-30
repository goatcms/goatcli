package entitymodel

import (
	"github.com/goatcms/goatcore/varutil/plainmap"
)

// Fields struct represent fields set
type Fields struct {
	ByName map[string]*Field
	ByType map[string][]*Field
}

// NewFields create new Fields instance
func NewFields() *Fields {
	return &Fields{
		ByName: map[string]*Field{},
		ByType: map[string][]*Field{},
	}
}

// NewFieldsFromPlainmap create new Fields instance and load data from plainmap
func NewFieldsFromPlainmap(baseKey string, data map[string]string) (fields *Fields, err error) {
	var (
		field *Field
		key   string
	)
	baseKey += "."
	fields = NewFields()
	for _, index := range plainmap.Keys(data, baseKey) {
		key = baseKey + index
		if field, err = NewFieldFromPlainmap(key, data); err != nil {
			return nil, err
		}
		fields.ByName[field.FullName.CamelCaseUF] = field
		fields.ByType[field.Type] = append(fields.ByType[field.Type], field)
	}
	return fields, nil
}
