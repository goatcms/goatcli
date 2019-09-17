package simpletf

import "testing"

func TestSum(t *testing.T) {
	t.Parallel()
	expectedResult := 50
	result := Sum(10, 15, 25)
	if result != expectedResult {
		t.Errorf("expected result is %v (and take %v)", expectedResult, result)
		return
	}
}

func TestMinus(t *testing.T) {
	t.Parallel()
	expectedResult := 5
	result := Sub(15, 10)
	if result != expectedResult {
		t.Errorf("expected result is %v (and take %v)", expectedResult, result)
		return
	}
}
