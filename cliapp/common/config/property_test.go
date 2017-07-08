package config

import "testing"

const (
	testPropertyJSON            = `{"key":"key", "prompt":"insert new value", "type":"alnum", "min":1, "max":22}`
	testPropertyBoolPatternJSON = `{"key":"key", "pattern": "(yes|no)", "prompt":"insert new value", "type":"line", "min":1, "max":4}`
	testPropertyArrayJSON       = `[{"key":"key", "prompt":"insert new value", "type":"alnum", "min":1, "max":22},{}]`
)

func TestProperty(t *testing.T) {
	t.Parallel()
	property, err := NewProperty([]byte(testPropertyJSON))
	if err != nil {
		t.Error(err)
		return
	}
	if property.Key != "key" {
		t.Errorf("incorrect key value parsing (expected key and take %s)", property.Key)
		return
	}
	if property.Type != "alnum" {
		t.Errorf("incorrect type value parsing (expected alnum and take %s)", property.Type)
		return
	}
	if property.Prompt != "insert new value" {
		t.Errorf("incorrect prompt value (expected 'insert new value' text and take '%s')", property.Prompt)
		return
	}
	if property.Min != 1 {
		t.Errorf("incorrect min value parsing (expected 1 and take %d)", property.Min)
		return
	}
	if property.Max != 22 {
		t.Errorf("incorrect max value parsing (expected 22 and take %d)", property.Max)
		return
	}
}

func TestNewProperties(t *testing.T) {
	t.Parallel()
	properties, err := NewProperties([]byte(testPropertyArrayJSON))
	if err != nil {
		t.Error(err)
		return
	}
	if len(properties) != 2 {
		t.Errorf("properties array should contains 2 elements (and it have %d)", len(properties))
		return
	}
	if properties[0].Key != "key" {
		t.Errorf("incorrect key value parsing (expected key and take %s)", properties[0].Key)
		return
	}
}

func TestNewPatternProperty(t *testing.T) {
	t.Parallel()
	property, err := NewProperty([]byte(testPropertyBoolPatternJSON))
	if err != nil {
		t.Error(err)
		return
	}
	if property.Pattern == nil {
		t.Errorf("pattern is defined (and can not be nil)")
		return
	}
	if property.Pattern.MatchString("yes") != true {
		t.Errorf("yes is correct pattern value")
		return
	}
	if property.Pattern.MatchString("no") != true {
		t.Errorf("no is correct pattern value")
		return
	}
	if property.Pattern.MatchString("nil") != false {
		t.Errorf("nil is incorrect pattern value")
		return
	}
}
