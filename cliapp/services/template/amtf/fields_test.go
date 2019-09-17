package amtf

import (
	"testing"

	"github.com/goatcms/goatcli/cliapp/common/am"
	"github.com/goatcms/goatcli/cliapp/common/am/entitymodel"
)

func TestLinkFieldUF(t *testing.T) {
	var (
		model    *am.ApplicationModel
		field    *entitymodel.Field
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
	if field, ok = entities["User"].AllFields.ByFullName["PersonFirstname"]; !ok {
		t.Errorf("expected field PersonFirstname")
		return
	}
	expected := "noModifiedVarName.Person.Firstname"
	result := LinkFieldUF("noModifiedVarName", field)
	if result != expected {
		t.Errorf("Expected %s and take %s", expected, result)
	}
}

func TestLinkFieldLF(t *testing.T) {
	var (
		model    *am.ApplicationModel
		field    *entitymodel.Field
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
	if field, ok = entities["User"].AllFields.ByFullName["PersonFirstname"]; !ok {
		t.Errorf("expected field PersonFirstname")
		return
	}
	expected := "NoModifiedVarName.person.firstname"
	result := LinkFieldLF("NoModifiedVarName", field)
	if result != expected {
		t.Errorf("Expected %s and take %s", expected, result)
	}
}
