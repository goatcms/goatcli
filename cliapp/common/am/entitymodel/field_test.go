package entitymodel

import (
	"testing"
)

func TestNewField(t *testing.T) {
	var (
		data  map[string]string
		field *Field
		err   error
	)
	t.Parallel()
	data = map[string]string{
		"prefix.name":     "struct.title",
		"prefix.type":     "Text",
		"prefix.system":   "y",
		"prefix.unique":   "y",
		"prefix.required": "y",
	}
	if field, err = NewFieldFromPlainmap("prefix", data); err != nil {
		t.Error(err)
		return
	}
	if field.Name.CamelCaseUF != "Title" {
		t.Errorf("field.Name.CamelCaseUF should be equals to 'Title' and take %s", field.Name.CamelCaseUF)
		return
	}
	if field.FullName.CamelCaseUF != "StructTitle" {
		t.Errorf("field.FullName.CamelCaseUF should be equals to 'StructTitle' and take %s", field.FullName.CamelCaseUF)
		return
	}
	if field.Type != "text" {
		t.Errorf("field.Type should be equals to 'text' and take %s", field.Name.CamelCaseUF)
		return
	}
	if field.Flags.Required != true {
		t.Errorf("field.Flags.Required should be equals to true")
		return
	}
	if field.Flags.System != true {
		t.Errorf("field.Flags.System should be equals to true")
		return
	}
	if field.Flags.Unique != true {
		t.Errorf("field.Flags.Unique should be equals to true")
		return
	}
}
