package vcs

import (
	"testing"

	"github.com/goatcms/goatcore/varutil"
	"github.com/goatcms/goatcore/workers"
)

func TestPersistedConcurrentlyRead(t *testing.T) {
	var (
		persisted = NewPersistedFiles(true)
	)
	t.Parallel()
	// create directories
	for ij := workers.AsyncTestReapeat; ij > 0; ij-- {
		path := varutil.RandString(10, varutil.AlphaNumericBytes)
		persisted.AddPath(path)
		for i := workers.MaxJob; i > 0; i-- {
			go (func() {
				persisted.AddPath(path)
			})()
			go (func() {
				if !persisted.ContainsPath(path) {
					t.Errorf("Should contains path")
					return
				}
			})()
		}
	}
}
