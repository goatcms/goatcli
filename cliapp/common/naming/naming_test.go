package naming

import "testing"

func TestToUnderscore(t *testing.T) {
	t.Parallel()
	result := ToUnderscore("FirstAndLast")
	if ToUnderscore("FirstAndLast") != "first_and_last" {
		t.Errorf("ToUnderscore('FirstAndLast') expected 'first_and_last' result and take %s", result)
		return
	}
	result = ToUnderscore("FirstJSON")
	if result != "first_json" {
		t.Errorf("ToUnderscore('FirstANDLast') expected 'first_json' result and take %s", result)
		return
	}
}

func TestToUnderscoreFromDot(t *testing.T) {
	t.Parallel()
	result := ToUnderscore("First.And.Last")
	if result != "first_and_last" {
		t.Errorf("ToUnderscore('First.And.Last') expected 'first_and_last' result and take %s", result)
		return
	}
	result = ToUnderscore("First.AND.Last")
	if result != "first_and_last" {
		t.Errorf("ToUnderscore('First.AND.Last') expected 'first_and_last' result and take %s", result)
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
	if ToCamelCase("first.and.last") != "firstAndLast" {
		t.Errorf("ToCamelCase('first.and.last') = %s expected 'firstAndLast' result", ToCamelCase("first.and.last"))
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
	result := ToCamelCaseUF("First_and_Last")
	if result != "FirstAndLast" {
		t.Errorf("ToCamelCaseUF('First_and_Last') = %s expected 'FirstAndLast' result", result)
		return
	}
	result = ToCamelCaseUF("FirstAndLast")
	if result != "FirstAndLast" {
		t.Errorf("ToCamelCaseUF('FirstAndLast') = %s expected 'FirstAndLast' result", result)
		return
	}
	result = ToCamelCaseUF("firstAndLast")
	if result != "FirstAndLast" {
		t.Errorf("ToCamelCaseUF('firstAndLast') = %s expected 'FirstAndLast' result", result)
		return
	}
}

func TestToUpperFirst(t *testing.T) {
	t.Parallel()
	result := ToUpperFirst("upperFirst")
	if result != "UpperFirst" {
		t.Errorf("ToUpperFirst('upperFirst') = %s expected 'UpperFirst' result", result)
		return
	}
}

func TestToLowerFirst(t *testing.T) {
	t.Parallel()
	result := ToLowerFirst("LOWER")
	if result != "lOWER" {
		t.Errorf("ToLowerFirst('LOWER') = %s expected 'lOWER' result", result)
		return
	}
}
