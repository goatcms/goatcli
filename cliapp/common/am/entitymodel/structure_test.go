package entitymodel

import (
	"testing"
)

func TestStructureAddField(t *testing.T) {
	var (
		instance *Structure
		err      error
		field    *Field
	)
	t.Parallel()
	if instance, err = NewStructure("ROOT", nil); err != nil {
		t.Error(err)
		return
	}
	if field, err = NewField("name"); err != nil {
		t.Error(err)
		return
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
		instance *Structure
		err      error
		relation *Relation
	)
	t.Parallel()
	if instance, err = NewStructure("ROOT", nil); err != nil {
		t.Error(err)
		return
	}
	if relation, err = NewRelation("name", "to"); err != nil {
		t.Error(err)
		return
	}
	if err = instance.AddRelation(relation); err != nil {
		t.Error(err)
		return
	}
	if relation.Structure != instance {
		t.Errorf("Relation structure should be updated")
	}
	if _, ok := instance.Relations.ByName["Name"]; !ok {
		t.Errorf("Name relation should be defined")
	}
}

func TestStructureAddFieldAndRelationWithTheSameName(t *testing.T) {
	var (
		instance *Structure
		err      error
		relation *Relation
		field    *Field
	)
	t.Parallel()
	if instance, err = NewStructure("ROOT", nil); err != nil {
		t.Error(err)
		return
	}
	if relation, err = NewRelation("name", "to"); err != nil {
		t.Error(err)
		return
	}
	if field, err = NewField("name"); err != nil {
		t.Error(err)
		return
	}
	if err = instance.AddRelation(relation); err != nil {
		t.Error(err)
		return
	}
	if err = instance.AddField(field); err == nil {
		t.Errorf("Expected error when define relation and field with the same name")
		return
	}
}
