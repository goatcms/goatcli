package config

import "testing"

const (
	testDataSetJSON  = `{"name":"data_name", "properties":[{"key":"key", "type":"alnum", "min":1, "max":22}]}`
	testDataSetsJSON = `[{"name":"data_name", "properties":[{"key":"key", "type":"alnum", "min":1, "max":22}]}]`
)

func TestDataSets(t *testing.T) {
	t.Parallel()
	c, err := NewDataSets([]byte(testDataSetsJSON))
	if err != nil {
		t.Error(err)
		return
	}
	if len(c) != 1 {
		t.Errorf("modules array should contains 1 element (and it have %d)", len(c))
		return
	}
	if c[0].Name != "data_name" {
		t.Errorf("wrong Name value (expected data_name and take %s)", c[0].Name)
		return
	}
	if len(c[0].Properties) != 1 {
		t.Errorf("expected one property")
		return
	}
}

func TestDataSet(t *testing.T) {
	t.Parallel()
	c, err := NewDataSet([]byte(testDataSetJSON))
	if err != nil {
		t.Error(err)
		return
	}
	if c.Name != "data_name" {
		t.Errorf("incorrect Name value parsing (expected data_name and take %s)", c.Name)
		return
	}
	if len(c.Properties) != 1 {
		t.Errorf("modules array should contains 1 element (and it have %d)", len(c.Properties))
		return
	}
}
