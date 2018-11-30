package config

import "testing"

func TestDependencies(t *testing.T) {
	t.Parallel()
	var (
		deps []*Dependency
		err  error
	)
	if deps, err = NewDependencies([]byte(`[{
		"repo": "RepoValue1",
		"branch": "BranchValue1",
		"rev": "RevValue1",
		"dest": "DestValue1",
	}, {
		"repo": "RepoValue2",
		"branch": "BranchValue2",
		"rev": "RevValue2",
		"dest": "DestValue2",
	}]`)); err != nil {
		t.Error(err)
		return
	}
	if len(deps) != 2 {
		t.Errorf("dependencies array should contains 2 elements (and it have %d)", len(deps))
		return
	}
	if deps[0].Repo != "RepoValue1" {
		t.Errorf("wrong repo value")
		return
	}
	if deps[0].Branch != "BranchValue1" {
		t.Errorf("wrong branch value")
		return
	}
	if deps[0].Rev != "RevValue1" {
		t.Errorf("wrong rev value")
		return
	}
	if deps[0].Dest != "DestValue1" {
		t.Errorf("wrong dest value")
		return
	}
}
