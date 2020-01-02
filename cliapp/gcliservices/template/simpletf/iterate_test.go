package simpletf

import (
	"regexp"
	"testing"

	"github.com/goatcms/goatcore/varutil"
)

func TestKeys(t *testing.T) {
	t.Parallel()
	m := map[string]string{
		"basekey.key1.subkey":      "somevalue1",
		"basekey.key2.subkey":      "somevalue2",
		"basekey.key3.subkey":      "somevalue3",
		"basekey.key4":             "somevalue4",
		"otherbasekey.key5.subkey": "somevalue5",
	}
	k := Keys(m, "basekey.")
	if len(k) != 4 {
		t.Errorf("expected result array contains 4 elements (and take %v)", k)
		return
	}
	if !varutil.IsArrContainStr(k, "key1") {
		t.Errorf("expected key1 in result array")
		return
	}
	if !varutil.IsArrContainStr(k, "key2") {
		t.Errorf("expected key2 in result array")
		return
	}
	if !varutil.IsArrContainStr(k, "key3") {
		t.Errorf("expected key3 in result array")
		return
	}
	if !varutil.IsArrContainStr(k, "key4") {
		t.Errorf("expected key3 in result array")
		return
	}
}

func TestStrainMap(t *testing.T) {
	var (
		result map[string]string
		err    error
	)
	t.Parallel()
	m := map[string]string{
		"basekey.key1.subkey":      "somevalue1",
		"basekey.key2.subkey":      "somevalue2",
		"basekey.key3.subkey":      "somevalue3",
		"basekey.key4":             "somevalue4",
		"otherbasekey.key5.subkey": "somevalue5",
	}
	if result, err = StrainMap(m, regexp.MustCompile("^basekey.[a-z0-9.]+$")); err != nil {
		t.Error(err)
		return
	}
	if len(result) != 4 {
		t.Errorf("expected result map contains 4 elements (and take %v)", result)
		return
	}
}
