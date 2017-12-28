package tfunc

import "testing"

func TestUniqueReduce(t *testing.T) {
	t.Parallel()
	soure := []string{"v1", "v2", "v2", "v3", "v3", "v4"}
	result := UniqueReduce(soure)
	if len(result) != 4 {
		t.Errorf("expected result array should contains 4 elements (and result is contain %v element)", len(result))
		return
	}
	if result[0] != "v1" {
		t.Errorf("result[0] shound be equals to 'v1' (and it is %v)", result[0])
		return
	}
	if result[1] != "v2" {
		t.Errorf("result[1] shound be equals to 'v2' (and it is %v)", result[1])
		return
	}
	if result[2] != "v3" {
		t.Errorf("result[2] shound be equals to 'v3' (and it is %v)", result[2])
		return
	}
	if result[3] != "v4" {
		t.Errorf("result[3] shound be equals to 'v4' (and it is %v)", result[3])
		return
	}
}
