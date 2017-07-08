package config

import "testing"

const (
	testCollectionsJSON = `[{"type":"array", "key":"key1", "prompt":"prompt", "properties":[{"key":"key", "type":"alnum", "min":1, "max":22}]}]`
	testCollectionJSON  = `{"type":"map", "key":"key1", "prompt":"prompt","properties":[{"key":"key", "type":"alnum", "min":1, "max":22}]}`
)

func TestCollections(t *testing.T) {
	t.Parallel()
	c, err := NewCollections([]byte(testCollectionsJSON))
	if err != nil {
		t.Error(err)
		return
	}
	if len(c) != 1 {
		t.Errorf("modules array should contains 1 element (and it have %d)", len(c))
		return
	}
	if c[0].Type != "array" {
		t.Errorf("wrong Type value (expected array and take %s)", c[0].Type)
		return
	}
	if len(c[0].Properties) != 1 {
		t.Errorf("properties array should contains 1 element (and it have %d)", len(c[0].Properties))
		return
	}
}

func TestCollection(t *testing.T) {
	t.Parallel()
	c, err := NewCollection([]byte(testCollectionJSON))
	if err != nil {
		t.Error(err)
		return
	}
	if c.Type != "map" {
		t.Errorf("incorrect Type value parsing (expected map and take %s)", c.Type)
		return
	}
	if len(c.Properties) != 1 {
		t.Errorf("properties array should contains 1 element (and it have %d)", len(c.Properties))
		return
	}
}
