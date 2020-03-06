package scripts

import (
	"sync"

	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/varutil/goaterr"
)

// History run application scripts
type History struct {
	runned map[string]bool
	mu     sync.RWMutex
}

// NewHistory return new history instance
func NewHistory() (instance *History) {
	return &History{}
}

// ScopeHistory create new History instance (and set it to the scope) or return exist singleton from the scope
func ScopeHistory(scp app.Scope) (history *History, err error) {
	var inst interface{}
	locker := scp.LockData()
	if inst, err = locker.Get(historyInstanceKey); err != nil {
		return nil, goaterr.ToError(goaterr.AppendError(nil, err, locker.Commit()))
	}
	if inst != nil {
		history = inst.(*History)
	} else {
		history = NewHistory()
		locker.Set(historyInstanceKey, history)
	}
	if err = locker.Commit(); err != nil {
		return nil, err
	}
	return history, nil
}

// Mark script as runned
func (history *History) Mark(name string) {
	history.mu.Lock()
	defer history.mu.Unlock()
	history.runned[name] = true
}

// IsRuned return true if script is runned
func (history *History) IsRuned(name string) bool {
	history.mu.RLock()
	defer history.mu.RUnlock()
	return history.runned[name]
}
