package tfunc

import "testing"

func TestDict(t *testing.T) {
	t.Parallel()
	m := Dict("Key1", int64(12), "Key2", "test")
	if len(m) != 2 {
		t.Errorf("result map should contains 2 elements (and take %v)", len(m))
		return
	}
	if m["Key1"].(int64) != 12 {
		t.Errorf("Key1 value should be equals to 12")
		return
	}
	if m["Key2"].(string) != "test" {
		t.Errorf("Key2 value should be equals to 'test'")
		return
	}
}
