package properties

import "testing"

func TestInjectToString(t *testing.T) {
	t.Parallel()
	properties := NewProperties(map[string]string{
		"key": "value",
	})
	result, err := properties.InjectToString("a {{key}} b")
	if err != nil {
		t.Error(err)
		return
	}
	if result != "a value b" {
		t.Errorf("incorrect result '%s' (expected 'a value b')", result)
		return
	}
}

func TestInjectToStringError(t *testing.T) {
	t.Parallel()
	properties := NewProperties(map[string]string{
		"key": "value",
	})
	if _, err := properties.InjectToString("{{incorrectkey}}"); err == nil {
		t.Errorf("InjectToString should return error when use unknow key")
		return
	}
}
