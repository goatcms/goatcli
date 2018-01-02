package tfunc

import (
	"testing"

	"github.com/goatcms/goatcore/varutil"
)

func TestValuesFor(t *testing.T) {
	t.Parallel()
	soure := map[string]string{
		"model.key1.name": "v1",
		"any.key1.name":   "err",
		"model.key2.name": "v2",
	}
	result := ValuesFor(`^model.[A-Za-z0-9]+.name$`, soure)
	if len(result) != 2 {
		t.Errorf("expected result array should contains 2 elements (and result is contain %v element)", len(result))
		return
	}
	if !varutil.IsArrContainStr(result, "v1") {
		t.Errorf("result should contains 'v1' (and it is %v)", result)
		return
	}
	if !varutil.IsArrContainStr(result, "v2") {
		t.Errorf("result should contains 'v2' (and it is %v)", result)
		return
	}
}

func TestFindRow(t *testing.T) {
	t.Parallel()
	soure := map[string]string{
		"model.key1.name": "v1",
		"any.key1.name":   "err",
		"model.key2.name": "v2",
	}
	result := FindRow("model.", "^[A-Za-z0-9]+$", ".name", "v1", soure)
	if result != "key1" {
		t.Errorf("expeccted key1 for value 'v1' (and result is %v)", result)
		return
	}
	result = FindRow("model.", "^[A-Za-z0-9]+$", ".name", "v2", soure)
	if result != "key2" {
		t.Errorf("expeccted key2 for value 'v2' (and result is %v)", result)
		return
	}
	result = FindRow("model.", "^[A-Za-z0-9]+$", ".name", "no_exist_key", soure)
	if result != "" {
		t.Errorf("expeccted empty string for value 'no_exist_key' (and result is %v)", result)
		return
	}
}
