package config

import "testing"

const (
	projectTestSimpleData = `{"projectPath":"projectPathValue", "modules":[]}`
)

func TestProject(t *testing.T) {
	t.Parallel()
	c, err := NewProject([]byte(projectTestSimpleData))
	if err != nil {
		t.Error(err)
		return
	}
	if c.ProjectPath != "projectPathValue" {
		t.Errorf("incorrect ProjectPath value parsing (expected projectPathValue and take %s)", c.ProjectPath)
		return
	}
}
