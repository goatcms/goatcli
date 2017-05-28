package config

import "testing"

const (
	replaceTestDataSuffixAsString = `{"from":"fromvalue", "to":"tovalue", "suffix":"suffix1"}`
	replaceTestDataSuffixAsArray  = `{"from":"fromvalue", "to":"tovalue", "suffix":["suffix1", "suffix2"]}`
	replacesTestArray             = `[{"from":"fromvalue", "to":"tovalue", "suffix":"suffix1"},{"from":"fromvalue", "to":"tovalue", "suffix":["suffix1", "suffix2"]}]`
	replacesTestFailArray         = `[{"from":"fromvalue", "to":"tovalue", "suffix":"suffix1"},"error"]`
)

func TestReplaceSuffixAsString(t *testing.T) {
	t.Parallel()
	rc, err := NewReplace([]byte(replaceTestDataSuffixAsString))
	if err != nil {
		t.Error(err)
		return
	}
	if rc.From != "fromvalue" {
		t.Errorf("incorrect from value parsing (expected fromvalue and take %s)", rc.From)
		return
	}
	if rc.To != "tovalue" {
		t.Errorf("incorrect to value parsing (expected tovalue and take %s)", rc.To)
		return
	}
	if rc.Suffix == nil {
		t.Errorf("suffix is nil")
		return
	}
	if len(rc.Suffix) != 1 {
		t.Errorf("suffix shoult have one element (and it have %d)", len(rc.Suffix))
		return
	}
	if rc.Suffix[0] != "suffix1" {
		t.Errorf("incorrect from Suffix[0] (expected suffix1 and take %s)", rc.Suffix[0])
		return
	}
}

func TestReplaceSuffixAsArray(t *testing.T) {
	t.Parallel()
	rc, err := NewReplace([]byte(replaceTestDataSuffixAsArray))
	if err != nil {
		t.Error(err)
		return
	}
	if rc.From != "fromvalue" {
		t.Errorf("incorrect from value parsing (expected fromvalue and take %s)", rc.From)
		return
	}
	if rc.To != "tovalue" {
		t.Errorf("incorrect to value parsing (expected tovalue and take %s)", rc.To)
		return
	}
	if rc.Suffix == nil {
		t.Errorf("suffix is nil")
		return
	}
	if len(rc.Suffix) != 2 {
		t.Errorf("suffix shoult have two elements (and it have %d)", len(rc.Suffix))
		return
	}
	if rc.Suffix[0] != "suffix1" {
		t.Errorf("incorrect from Suffix[0] (expected suffix1 and take %s)", rc.Suffix[0])
		return
	}
	if rc.Suffix[1] != "suffix2" {
		t.Errorf("incorrect from Suffix[1] (expected suffix2 and take %s)", rc.Suffix[1])
		return
	}
}

func TestNewReplaces(t *testing.T) {
	t.Parallel()
	replaces, err := NewReplaces([]byte(replacesTestArray))
	if err != nil {
		t.Error(err)
		return
	}
	if len(replaces) != 2 {
		t.Errorf("replaces array should contains 2 elements (and it have %d)", len(replaces))
		return
	}
	if replaces[0].From != "fromvalue" {
		t.Errorf("incorrect from value parsing (expected fromvalue and take %s)", replaces[0].From)
		return
	}
	if replaces[0].To != "tovalue" {
		t.Errorf("incorrect to value parsing (expected tovalue and take %s)", replaces[0].To)
		return
	}
	if replaces[0].Suffix == nil {
		t.Errorf("suffix is nil")
		return
	}
	if len(replaces[0].Suffix) != 1 {
		t.Errorf("suffix shoult have one element (and it have %d)", len(replaces[0].Suffix))
		return
	}
	if replaces[0].Suffix[0] != "suffix1" {
		t.Errorf("incorrect from Suffix[0] (expected suffix1 and take %s)", replaces[0].Suffix[0])
		return
	}
}

func TestNewReplacesFail(t *testing.T) {
	t.Parallel()
	_, err := NewReplaces([]byte(replacesTestFailArray))
	if err.Error() != "NewReplaces array  must contains replace objects only" {
		t.Error(err)
		return
	}
}
