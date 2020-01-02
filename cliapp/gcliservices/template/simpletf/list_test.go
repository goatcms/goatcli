package simpletf

import (
	"testing"

	"github.com/goatcms/goatcore/varutil"
)

func TestUnique(t *testing.T) {
	t.Parallel()
	source := []string{"v1", "v2", "v2", "v3", "v3", "v4"}
	result := Unique(source)
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

func TestUnion(t *testing.T) {
	t.Parallel()
	source1 := []string{"v1", "v2", "v2", "v3"}
	source2 := []string{"v3", "v4"}
	result := Union(source1, source2)
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

func TestExcept(t *testing.T) {
	t.Parallel()
	source1 := []string{"v1", "v2", "v2", "v3"}
	source2 := []string{"v3", "v3", "v4"}
	result := Except(source1, source2)
	if len(result) != 2 {
		t.Errorf("expected result array should contains 2 elements (and result is contain %v element)", len(result))
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
}

func TestIntersect(t *testing.T) {
	t.Parallel()
	source1 := []string{"v1", "v2", "v2", "v3"}
	source2 := []string{"v3", "v3", "v4"}
	result := Intersect(source1, source2)
	if len(result) != 1 {
		t.Errorf("expected result array should contains 1 elements (and result is contain %v element)", len(result))
		return
	}
	if result[0] != "v3" {
		t.Errorf("result[0] shound be equals to 'v3' (and it is %v)", result[0])
		return
	}
}

func TestSort(t *testing.T) {
	t.Parallel()
	source1 := []string{"v2", "v1", "v4", "v3"}
	result := Sort(source1)
	if len(result) != 4 {
		t.Errorf("expected result array should contains 4 elements (and result is contain %v element)", len(result))
		return
	}
	if result[0] != "v1" {
		t.Errorf("result[0] shound be equals to 'v1' (and it is %v)", result[0])
		return
	}
	if result[1] != "v2" {
		t.Errorf("result[0] shound be equals to 'v2' (and it is %v)", result[1])
		return
	}
	if result[2] != "v3" {
		t.Errorf("result[0] shound be equals to 'v3' (and it is %v)", result[2])
		return
	}
	if result[3] != "v4" {
		t.Errorf("result[0] shound be equals to 'v4' (and it is %v)", result[3])
		return
	}
}

func TestRandomValue(t *testing.T) {
	t.Parallel()
	expectedValues := []string{"v1", "v2", "v2", "v3"}
	result := RandomValue(expectedValues...)
	if !varutil.IsArrContainStr(expectedValues, result) {
		t.Errorf("should return value from parameters %v and return %v", expectedValues, result)
		return
	}
}
