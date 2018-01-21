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

func TestFindRows(t *testing.T) {
	t.Parallel()
	var result []string
	soure := map[string]string{
		"model.key1.name": "v1",
		"any.key1.name":   "err",
		"model.key2.name": "v2",
		"model.key3.name": "v1",
	}
	result = FindRows("model.", "^[A-Za-z0-9]+$", ".name", "v1", soure)
	if len(result) != 2 {
		t.Errorf("result should contains 2 elements and it is %v", result)
		return
	}
	if !varutil.IsArrContainStr(result, "key1") {
		t.Errorf("result should contains 'key1' and it is %v", result)
		return
	}
	if !varutil.IsArrContainStr(result, "key3") {
		t.Errorf("result should contains 'key3' and it is %v", result)
		return
	}
}

func TestSubMap(t *testing.T) {
	t.Parallel()
	var result map[string]string
	soure := map[string]string{
		"model.key1.name1": "v1",
		"model.key1.name2": "v2",
		"any.key1.name":    "err",
		"model.key2.name":  "v2",
		"model.key3.name":  "v1",
	}
	result = SubMap("model.key1.", "newmap.", soure)
	if len(result) != 2 {
		t.Errorf("result should contains 2 elements and it is %v", result)
		return
	}
	if result["newmap.name1"] != "v1" {
		t.Errorf("newmap.name1 key should contains 'v1' value and it is %v (%v)", result["newmap.name1"], result)
		return
	}
	if result["newmap.name2"] != "v2" {
		t.Errorf("newmap.name2 key should contains 'v2' value and it is %v (%v)", result["newmap.name2"], result)
		return
	}
}
