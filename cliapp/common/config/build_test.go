package config

import "testing"

const (
	testBuildJSON  = `{"from":"fromv", "afterBuild":"afterBuildV", "to":"tov","layout":"layoutv","template":"templatev", "properties":{"key1":"value1"}}`
	testBuildsJSON = `[{"from":"fromv", "afterBuild":"afterBuildV", "to":"tov","layout":"layoutv","template":"templatev", "properties":{"key1":"value1"}}]`
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
	if c[0].Template != "templatev" {
		t.Errorf("wrong Template value (expected templatev and take %s)", c[0].Template)
		return
	}
	if c[0].AfterBuild != "afterBuildV" {
		t.Errorf("wrong AfterBuild value (expected 'afterBuildV' and take %s)", c[0].AfterBuild)
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
	if c.Template != "templatev" {
		t.Errorf("wrong Template value (expected templatev and take %s)", c.Template)
		return
	}
	if c.AfterBuild != "afterBuildV" {
		t.Errorf("wrong AfterBuild value (expected 'afterBuildV' and take %s)", c.AfterBuild)
		return
	}
	if len(c.Properties) != 1 {
		t.Errorf("expected one property")
		return
	}
}
