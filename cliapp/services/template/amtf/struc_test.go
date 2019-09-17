package amtf

import (
	"testing"

	"github.com/goatcms/goatcli/cliapp/common/am"
	"github.com/goatcms/goatcli/cliapp/common/am/entitymodel"
)

func TestStructClassName(t *testing.T) {
	var (
		model    *am.ApplicationModel
		struc    *entitymodel.Structure
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
	if struc, ok = entities["User"].RootStructure.Structures.ByName["Person"]; !ok {
		t.Errorf("expected structure Person")
		return
	}
	expected := "UserPerson"
	result := StructClassName(struc)
	if result != expected {
		t.Errorf("Expected %s and take %s", expected, result)
	}
}
