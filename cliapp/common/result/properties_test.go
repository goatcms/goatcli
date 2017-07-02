package result

import "testing"

func TestInjectToString(t *testing.T) {
	t.Parallel()
	propertiesResult := NewPropertiesResult(map[string]string{
		"key": "value",
	})
	result, err := propertiesResult.InjectToString("a {{key}} b")
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
	propertiesResult := NewPropertiesResult(map[string]string{
		"key": "value",
	})
	if _, err := propertiesResult.InjectToString("{{incorrectkey}}"); err == nil {
		t.Errorf("InjectToString should return error when use unknow key")
		return
	}
}

func TestGet(t *testing.T) {
	var (
		value string
		err   error
	)
	t.Parallel()
	propertiesResult := NewPropertiesResult(map[string]string{
		"key": "value",
	})
	if value, err = propertiesResult.Get("key"); err != nil {
		t.Error(err)
		return
	}
	if value != "value" {
		t.Errorf("incorrect value %s", value)
		return
	}
}
