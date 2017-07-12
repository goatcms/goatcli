package tfunc

import "testing"

func TestToUnderscore(t *testing.T) {
	t.Parallel()
	if ToUnderscore("FirstAndLast") != "first_and_last" {
		t.Errorf("ToUnderscore('FirstAndLast') expected 'first_and_last' result")
		return
	}
	if ToUnderscore("FirstANDLast") != "first_and_last" {
		t.Errorf("ToUnderscore('FirstANDLast') expected 'first_and_last' result")
		return
	}
}

func TestToCamelCase(t *testing.T) {
	t.Parallel()
	if ToCamelCase("First_and_Last") != "FirstAndLast" {
		t.Errorf("ToCamelCase('First_and_Last') = %s expected 'FirstAndLast' result", ToCamelCase("First_and_Last"))
		return
	}
	if ToCamelCase("FirstAndLast") != "FirstAndLast" {
		t.Errorf("ToCamelCase('FirstAndLast') = %s expected 'FirstAndLast' result", ToCamelCase("FirstAndLast"))
		return
	}
	if ToCamelCase("firstAndLast") != "firstAndLast" {
		t.Errorf("ToCamelCase('firstAndLast') = %s expected 'firstAndLast' result", ToCamelCase("firstAndLast"))
		return
	}
}

func TestToCamelCaseLF(t *testing.T) {
	t.Parallel()
	if ToCamelCaseLF("First_and_Last") != "firstAndLast" {
		t.Errorf("ToCamelCaseLF('First_and_Last') = %s expected 'firstAndLast' result", ToCamelCaseLF("First_and_Last"))
		return
	}
	if ToCamelCaseLF("FirstAndLast") != "firstAndLast" {
		t.Errorf("ToCamelCaseLF('FirstAndLast') = %s expected 'firstAndLast' result", ToCamelCaseLF("FirstAndLast"))
		return
	}
	if ToCamelCaseLF("firstAndLast") != "firstAndLast" {
		t.Errorf("ToCamelCaseLF('firstAndLast') = %s expected 'firstAndLast' result", ToCamelCaseLF("firstAndLast"))
		return
	}
}

func TestToCamelCaseUF(t *testing.T) {
	t.Parallel()
	if ToCamelCaseUF("First_and_Last") != "FirstAndLast" {
		t.Errorf("ToCamelCaseUF('First_and_Last') = %s expected 'FirstAndLast' result", ToCamelCaseUF("First_and_Last"))
		return
	}
	if ToCamelCaseUF("FirstAndLast") != "FirstAndLast" {
		t.Errorf("ToCamelCaseUF('FirstAndLast') = %s expected 'FirstAndLast' result", ToCamelCaseUF("FirstAndLast"))
		return
	}
	if ToCamelCaseUF("firstAndLast") != "FirstAndLast" {
		t.Errorf("ToCamelCaseUF('firstAndLast') = %s expected 'FirstAndLast' result", ToCamelCaseUF("firstAndLast"))
		return
	}
}

func TestToUpperFirst(t *testing.T) {
	t.Parallel()
	if ToUpperFirst("upperFirst") != "UpperFirst" {
		t.Errorf("ToUpperFirst('upperFirst') = %s expected 'UpperFirst' result", ToUpperFirst("upperFirst"))
		return
	}
}

func TestToLowerFirst(t *testing.T) {
	t.Parallel()
	if ToLowerFirst("LOWER") != "lOWER" {
		t.Errorf("ToLowerFirst('LOWER') = %s expected 'lOWER' result", ToCamelCaseUF("LOWER"))
		return
	}
}
