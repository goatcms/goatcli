package entitymodel

import (
	"testing"
)

func TestStructureAddField(t *testing.T) {
	var (
		instance = NewStructure(nil)
		err      error
	)
	t.Parallel()
	name, _ := NewName("name")
	field := &Field{
		Name: name,
	}
	if err = instance.AddField(field); err != nil {
		t.Error(err)
		return
	}
	if field.Structure != instance {
		t.Errorf("Field structure should be updated")
	}
	if _, ok := instance.Fields.ByName["Name"]; !ok {
		t.Errorf("Name field should be defined")
	}
}

func TestStructureAddRelation(t *testing.T) {
	var (
		instance = NewStructure(nil)
		err      error
	)
	t.Parallel()
	name, _ := NewName("name")
	relation := &Relation{
		Name: name,
	}
	if err = instance.AddRelation(relation); err != nil {
		t.Error(err)
		return
	}
	if relation.Structure != instance {
		t.Errorf("Relation structure should be updated")
	}
	if _, ok := instance.Relations["Name"]; !ok {
		t.Errorf("Name relation should be defined")
	}
}

func TestStructureAddFieldAndRelationWithTheSameName(t *testing.T) {
	var (
		instance = NewStructure(nil)
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
