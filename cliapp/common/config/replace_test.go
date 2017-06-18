package config

import "testing"

const (
	replaceTestData       = `{"from":"fromvalue", "to":"tovalue", "pattern":"[A-Za-z0-9_/]*.go"}`
	replacesTestArray     = `[{"from":"fromvalue", "to":"tovalue", "pattern":"[A-Za-z0-9_/]*.go"},{"from":"fromvalue", "to":"tovalue", "pattern":"[A-Za-z0-9_/]*.css"}]`
	replacesTestFailArray = `[{"from":"fromvalue", "to":"tovalue", "pattern":"[A-Za-z0-9_/]*.go"},"error"]`
)

func TestReplace(t *testing.T) {
	t.Parallel()
	ti := NewTestStringInjector()
	rc, err := NewReplace([]byte(replaceTestData), ti)
	if err != nil {
		t.Error(err)
		return
	}
	if !rc.From.MatchString("fromvalue") {
		t.Errorf("fromvalue is not match")
		return
	}
	if rc.To != "tovalue" {
		t.Errorf("incorrect to value parsing (expected tovalue and take %s)", rc.To)
		return
	}
	if rc.Pattern.MatchString("path/to/main/file.go") != true {
		t.Errorf("path/to/main/file.go is *.go file")
		return
	}
	if rc.Pattern.MatchString("path/to/main/file.css") == true {
		t.Errorf("path/to/main/file.css is not *.go file")
		return
	}
}

func TestNewReplaces(t *testing.T) {
	t.Parallel()
	ti := NewTestStringInjector()
	replaces, err := NewReplaces([]byte(replacesTestArray), ti)
	if err != nil {
		t.Error(err)
		return
	}
	if len(replaces) != 2 {
		t.Errorf("replaces array should contains 2 elements (and it have %d)", len(replaces))
		return
	}
	if !replaces[0].From.MatchString("fromvalue") {
		t.Errorf("fromvalue is not match")
		return
	}
	if replaces[0].To != "tovalue" {
		t.Errorf("incorrect to value parsing (expected tovalue and take %s)", replaces[0].To)
		return
	}
	if replaces[0].Pattern.MatchString("path/to/main/file.go") != true {
		t.Errorf("path/to/main/file.go is *.go file")
		return
	}
	if replaces[0].Pattern.MatchString("path/to/main/file.css") == true {
		t.Errorf("path/to/main/file.css is not *.go file")
		return
	}
	if replaces[1].Pattern.MatchString("path/to/main/file.go") == true {
		t.Errorf("path/to/main/file.go is not *.css file")
		return
	}
	if replaces[1].Pattern.MatchString("path/to/main/file.css") != true {
		t.Errorf("path/to/main/file.css is *.css file")
		return
	}
}

func TestNewReplacesFail(t *testing.T) {
	t.Parallel()
	ti := NewTestStringInjector()
	_, err := NewReplaces([]byte(replacesTestFailArray), ti)
	if err.Error() != "NewReplaces array  must contains replace objects only" {
		t.Error(err)
		return
	}
}
