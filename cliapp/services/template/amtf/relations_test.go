package amtf

import (
	"testing"

	"github.com/goatcms/goatcli/cliapp/common/am"
	"github.com/goatcms/goatcli/cliapp/common/am/entitymodel"
)

func TestLinkRelationUF(t *testing.T) {
	var (
		model    *am.ApplicationModel
		relation *entitymodel.Relation
		entities entitymodel.Entities
		err      error
		ok       bool
	)
	t.Parallel()
	if model, err = testApplicationModel(); err != nil {
		t.Error(err)
		return
	}
	if entities, err = model.Entities(); err != nil {
		t.Error(err)
		return
	}
	if relation, ok = entities["User"].AllRelations.ByFullName["FamilyParent"]; !ok {
		t.Errorf("expected relation PersonFirstname")
		return
	}
	expected := "noModifiedVarName.Family.Parent"
	result := LinkRelationUF("noModifiedVarName", relation)
	if result != expected {
		t.Errorf("Expected %s and take %s", expected, result)
	}
}

func TestLinkRelationLF(t *testing.T) {
	var (
		model    *am.ApplicationModel
		relation *entitymodel.Relation
		entities entitymodel.Entities
		err      error
		ok       bool
	)
	t.Parallel()
	if model, err = testApplicationModel(); err != nil {
		t.Error(err)
		return
	}
	if entities, err = model.Entities(); err != nil {
		t.Error(err)
		return
	}
	if relation, ok = entities["User"].AllRelations.ByFullName["FamilyParent"]; !ok {
		t.Errorf("expected relation FamilyParent")
		return
	}
	expected := "NoModifiedVarName.family.parent"
	result := LinkRelationLF("NoModifiedVarName", relation)
	if result != expected {
		t.Errorf("Expected %s and take %s", expected, result)
	}
}
