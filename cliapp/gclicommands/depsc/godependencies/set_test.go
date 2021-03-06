package godependencies

import (
	"testing"

	"github.com/goatcms/goatcli/cliapp/common/config"
)

func TestSetAddDependenciesManyTimes(t *testing.T) {
	t.Parallel()
	var (
		set = NewSet()
	)
	set.Add([]*config.Dependency{
		&config.Dependency{
			Repo:   "http://github.com/goatcms/goatcli",
			Branch: "master",
			Rev:    "",
			Dest:   "vendor/github.com/goatcms/goatcli",
		},
		&config.Dependency{
			Repo:   "http://github.com/goatcms/goatcore",
			Branch: "master",
			Rev:    "",
			Dest:   "vendor/github.com/goatcms/goatcore",
		},
	})
	set.Add([]*config.Dependency{
		&config.Dependency{
			Repo:   "http://github.com/goatcms/goatcms",
			Branch: "master",
			Rev:    "",
			Dest:   "vendor/github.com/goatcms/goatcms",
		},
	})
	deps := set.Dependencies()
	if len(deps) != 3 {
		t.Errorf("Result should contains all added dependencies. Expected 3 and take %v", len(deps))
	}
}

func TestSetAddDuplicatedDependencies(t *testing.T) {
	t.Parallel()
	var (
		set = NewSet()
		err error
	)
	if set.Add([]*config.Dependency{
		&config.Dependency{
			Repo:   "http://github.com/goatcms/goatcli",
			Branch: "master",
			Rev:    "",
			Dest:   "vendor/github.com/goatcms/goatcli",
		},
	}); err != nil {
		t.Error(err)
		return
	}
	if err = set.Add([]*config.Dependency{
		&config.Dependency{
			Repo:   "http://github.com/goatcms/goatcli",
			Branch: "master",
			Rev:    "",
			Dest:   "vendor/github.com/goatcms/goatcli",
		},
	}); err == nil {
		t.Errorf("Expected error when duplicate destination path")
		return
	}
	deps := set.Dependencies()
	if len(deps) != 1 {
		t.Errorf("Set should remove duplicated dependency. Expected 1 dependency and take %v", len(deps))
	}
}

/*
func TestSetAddGOPath(t *testing.T) {
	t.Parallel()
	var (
		err error
		set = NewSet()
	)
	set.Add([]*config.Dependency{
		&config.Dependency{
			Repo:   "http://github.com/goatcms/goatcli",
			Branch: "master",
			Rev:    "",
			Dest:   "vendor/github.com/goatcms/goatcli",
		},
		&config.Dependency{
			Repo:   "http://github.com/goatcms/goatcore",
			Branch: "master",
			Rev:    "",
			Dest:   "vendor/github.com/goatcms/goatcore",
		},
		&config.Dependency{
			Repo:   "http://github.com/goatcms/goatcms",
			Branch: "master",
			Rev:    "",
			Dest:   "vendor/github.com/goatcms/goatcms",
		},
	})
	deps := set.Dependencies()
	if len(deps) != 3 {
		t.Errorf("Result should contains all added dependencies. Expected 3 and take %v", len(deps))
	}
}

func TestSetAddDuplicatedGOPath(t *testing.T) {
	var (
		err error
		set = NewSet()
	)
	t.Parallel()
	if _, err = set.AddSource("http://github.com/goatcms/goatcli"); err != nil {
		t.Error(err)
		return
	}
	if _, err = set.AddSource("http://github.com/goatcms/goatcli"); err == nil {
		t.Errorf("AddSource duplicate should return an error")
		return
	}
	deps := set.Dependencies()
	if len(deps) != 1 {
		t.Errorf("Set should remove duplicated dependency. Expected 1 dependency and take %v", len(deps))
	}
}

func TestAddSourceReturnRow(t *testing.T) {
	var (
		err error
		set = NewSet()
		row *SetRow
	)
	t.Parallel()
	if row, err = set.AddSource("http://github.com/goatcms/goatcli"); err != nil {
		t.Error(err)
		return
	}
	if row == nil {
		t.Errorf("Set.AddSource should return new row")
	}
}*/

func TestSetImported(t *testing.T) {
	t.Parallel()
	var (
		err    error
		set    = NewSet()
		row    *SetRow
		result bool
	)
	if err = set.Add([]*config.Dependency{
		&config.Dependency{
			Repo:   "http://github.com/goatcms/goatcms",
			Branch: "master",
			Rev:    "",
			Dest:   "vendor/github.com/goatcms/goatcms",
		},
	}); err != nil {
		t.Error(err)
		return
	}
	if row = set.Row("vendor/github.com/goatcms/goatcms"); row == nil {
		t.Errorf("Row must be defined")
		return
	}
	row.SetImported(true)
	if result = set.Row("vendor/github.com/goatcms/goatcms").Imported; result == false {
		t.Errorf("Expected importet flag equels to true and take %v", result)
	}
}

func TestSetImportedWrongPath(t *testing.T) {
	t.Parallel()
	var (
		set = NewSet()
		row *SetRow
	)
	if row = set.Row("github.com/goatcms/goatcms"); row != nil {
		t.Errorf("Expected undefined row equals to nil")
	}
}

func TestImportedResolveStory(t *testing.T) {
	t.Parallel()
	var (
		set = NewSet()
	)
	if set.Resolve("github.com/goatcms/goatcms") == true {
		t.Errorf("should return false for undefined path")
	}
	set.SetResolve("github.com/goatcms/goatcms", true)
	if set.Resolve("github.com/goatcms/goatcms") == false {
		t.Errorf("should return true for resolved path")
	}
	// Last case is unused now
	// We don't support unresolved path
	set.SetResolve("github.com/goatcms/goatcms", false)
	if set.Resolve("github.com/goatcms/goatcms") == true {
		t.Errorf("should return true for unresolved path")
	}
}
