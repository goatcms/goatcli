package am

import (
	"testing"
)

func TestNewFieldSet(t *testing.T) {
	var (
		data  map[string]string
		set   *FieldSet
		flags *FieldSet
		err   error
		ok    bool
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
	if set, err = NewFieldSet("prefix", data); err != nil {
		t.Error(err)
		return
	}
	if _, ok = set.ByName["Title"]; !ok {
		t.Errorf("expected Title field")
		return
	}
	if _, ok = set.ByName["Counter"]; !ok {
		t.Errorf("expected Counter field")
		return
	}
	if _, ok = set.ByName["FlagsPrivate"]; !ok {
		t.Errorf("expected FlagsPrivate field")
		return
	}
	if flags, ok = set.Structs["Flags"]; !ok {
		t.Errorf("expected Flags struct")
		return
	}
	if _, ok = flags.ByName["Private"]; !ok {
		t.Errorf("expected Private field in Flags struct")
		return
	}
	if len(set.ByType["int"]) != 1 {
		t.Errorf("expected one int type field")
		return
	}
	if len(set.ByType["text"]) != 1 {
		t.Errorf("expected one text type field")
		return
	}
	if len(set.ByType["unknowntype"]) != 0 {
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
	if _, err = NewFieldSet("prefix", data); err == nil {
		t.Errorf("NewFieldSet should return error. Name is required")
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
	if _, err = NewFieldSet("prefix", data); err == nil {
		t.Errorf("NewFieldSet should return error. Type is required")
		return
	}
}
