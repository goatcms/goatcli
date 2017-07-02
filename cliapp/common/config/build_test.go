package config

import "testing"

const (
	testBuildJSON  = `{"from":"fromv","to":"tov","layout":"layoutv","view":"viewv", "properties":{"key1":"value1"}}`
	testBuildsJSON = `[{"from":"fromv","to":"tov","layout":"layoutv","view":"viewv", "properties":{"key1":"value1"}}]`
)

func TestBuilds(t *testing.T) {
	t.Parallel()
	c, err := NewBuilds([]byte(testBuildsJSON))
	if err != nil {
		t.Error(err)
		return
	}
	if len(c) != 1 {
		t.Errorf("build array should contains 1 element (and it have %d)", len(c))
		return
	}
	if c[0].From != "fromv" {
		t.Errorf("wrong From value (expected fromv and take %s)", c[0].From)
		return
	}
	if c[0].To != "tov" {
		t.Errorf("wrong To value (expected tov and take %s)", c[0].To)
		return
	}
	if c[0].Layout != "layoutv" {
		t.Errorf("wrong Layout value (expected layout and take %s)", c[0].Layout)
		return
	}
	if c[0].View != "viewv" {
		t.Errorf("wrong View value (expected viewv and take %s)", c[0].From)
		return
	}
	if len(c[0].Properties) != 1 {
		t.Errorf("expected one property")
		return
	}
}

func TestBuild(t *testing.T) {
	t.Parallel()
	c, err := NewBuild([]byte(testBuildJSON))
	if err != nil {
		t.Error(err)
		return
	}
	if c.From != "fromv" {
		t.Errorf("wrong From value (expected fromv and take %s)", c.From)
		return
	}
	if c.To != "tov" {
		t.Errorf("wrong To value (expected tov and take %s)", c.To)
		return
	}
	if c.Layout != "layoutv" {
		t.Errorf("wrong Layout value (expected layout and take %s)", c.Layout)
		return
	}
	if c.View != "viewv" {
		t.Errorf("wrong View value (expected viewv and take %s)", c.From)
		return
	}
	if len(c.Properties) != 1 {
		t.Errorf("expected one property")
		return
	}
}
