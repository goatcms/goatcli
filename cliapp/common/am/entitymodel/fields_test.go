package entitymodel

import (
	"testing"
)

func TestNewFields(t *testing.T) {
	var (
		data   map[string]string
		fields *Fields
		err    error
		ok     bool
	)
	t.Parallel()
	data = map[string]string{
		"prefix.somestring.name": "title",
		"prefix.somestring.type": "text",
		"prefix.someint.name":    "counter",
		"prefix.someint.type":    "int",
		"prefix.someStruct.name": "flags.private",
		"prefix.someStruct.type": "bool",
	}
	if fields, err = NewFieldsFromPlainmap("prefix", data); err != nil {
		t.Error(err)
		return
	}
	if _, ok = fields.ByName["Title"]; !ok {
		t.Errorf("expected Title field")
		return
	}
	if _, ok = fields.ByName["Counter"]; !ok {
		t.Errorf("expected Counter field")
		return
	}
	if _, ok = fields.ByName["FlagsPrivate"]; !ok {
		t.Errorf("expected FlagsPrivate field")
		return
	}
	if len(fields.ByType["int"]) != 1 {
		t.Errorf("expected one int type field")
		return
	}
	if len(fields.ByType["text"]) != 1 {
		t.Errorf("expected one text type field")
		return
	}
	if len(fields.ByType["unknowntype"]) != 0 {
		t.Errorf("unknowntype should be empty")
		return
	}
}

func TestNameIsRequired(t *testing.T) {
	var (
		data map[string]string
		err  error
	)
	t.Parallel()
	data = map[string]string{
		"prefix.somestring.entity": "",
		"prefix.somestring.type":   "text",
	}
	if _, err = NewFieldsFromPlainmap("prefix", data); err == nil {
		t.Errorf("NewFields should return error. Name is required")
		return
	}
}

func TestTypeIsRequired(t *testing.T) {
	var (
		data map[string]string
		err  error
	)
	t.Parallel()
	data = map[string]string{
		"prefix.somestring.entity": "name",
		"prefix.somestring.type":   "",
	}
	if _, err = NewFieldsFromPlainmap("prefix", data); err == nil {
		t.Errorf("NewFields should return error. Type is required")
		return
	}
}
