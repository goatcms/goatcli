package entitymodel

import (
	"testing"
)

func TestStructureAddField(t *testing.T) {
	var (
		instance = NewStructure()
		err      error
	)
	t.Parallel()
	name, _ := NewName("name")
	if err = instance.AddField(&Field{
		Name: name,
	}); err != nil {
		t.Error(err)
		return
	}
	if _, ok := instance.Fields.ByName["Name"]; !ok {
		t.Errorf("Name field should be defined")
		return
	}
}

func TestStructureAddFieldAndRelationWithTheSameName(t *testing.T) {
	var (
		instance = NewStructure()
		err      error
	)
	t.Parallel()
	name, _ := NewName("name")
	if err = instance.AddRelation(&Relation{
		Name: name,
	}); err != nil {
		t.Error(err)
		return
	}
	if err = instance.AddField(&Field{
		Name: name,
	}); err == nil {
		t.Errorf("Expected error when define relation and field with the same name")
		return
	}
}
