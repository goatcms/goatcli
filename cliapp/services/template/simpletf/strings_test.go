package simpletf

import "testing"

func TestInjectValues(t *testing.T) {
	t.Parallel()
	expectedResult := "/example/path/value1/value2"
	result := InjectValues("/example/path/{{$key1}}/{{$key2}}", map[string]string{
		"key1": "value1",
		"key2": "value2",
	})
	if result != expectedResult {
		t.Errorf("expected result is '%v' (and take '%v')", expectedResult, result)
		return
	}
}

func TestReplace(t *testing.T) {
	t.Parallel()
	expectedResult := "/example/path/value1"
	result := Replace("/example/path/{{$key1}}", "{{$key1}}", "value1")
	if result != expectedResult {
		t.Errorf("expected result is '%v' (and take '%v')", expectedResult, result)
		return
	}
}
