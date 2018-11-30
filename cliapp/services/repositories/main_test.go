package repositories

import "testing"

func TestReduceRepoURL(t *testing.T) {
	t.Parallel()
	result := reduceRepoURL("https://github.com/goatcms/goatcli.git")
	if result != "github.com/goatcms/goatcli" {
		t.Errorf("wrong reduceRepoURL result: %s", result)
	}
}
