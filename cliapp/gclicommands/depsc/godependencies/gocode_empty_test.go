package godependencies

import (
	"testing"
)

func TestEmptyFindMultilineImports(t *testing.T) {
	t.Parallel()
	var (
		result []string
		err    error
	)
	if result, err = FindImports(`
		import()
		`); err != nil {
		t.Error(err)
		return
	}
	if len(result) != 0 {
		t.Errorf("Expected 0 results")
		return
	}
}
