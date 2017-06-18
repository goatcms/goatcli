package config

import "testing"

const (
	propertyTestData      = `{"key":"key", "type":"alnum", "min":1, "max":22}`
	propertyTestArrayData = `[{"key":"key", "type":"alnum", "min":1, "max":22},{}]`
)

func TestProperty(t *testing.T) {
	t.Parallel()
	property, err := NewProperty([]byte(propertyTestData))
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
	properties, err := NewProperties([]byte(propertyTestArrayData))
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
