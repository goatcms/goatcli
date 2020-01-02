package simpletf

import "testing"

func TestSequencerHasNext(t *testing.T) {
	t.Parallel()
	sequencer := NewSequencer("v1", "v2", "v3")
	result := []string{}
	for sequencer.HasNext() {
		result = append(result, sequencer.Next())
	}
	if len(result) != 3 {
		t.Errorf("Expected one result")
	}
	if result[0] != "v1" {
		t.Errorf("Expected v1")
	}
	if result[1] != "v2" {
		t.Errorf("Expected v2")
	}
	if result[2] != "v3" {
		t.Errorf("Expected v3")
	}
	afterAllValue := sequencer.Next()
	if afterAllValue != "v3" {
		t.Errorf("after end should return last element every time and take %s", afterAllValue)
	}
}

func TestSequencerNext(t *testing.T) {
	t.Parallel()
	sequencer := NewSequencer("v1", "v2", "v3")
	if v := sequencer.Next(); v != "v1" {
		t.Errorf("Expected v1 and take %s", v)
	}
	if v := sequencer.Next(); v != "v2" {
		t.Errorf("Expected v2 and take %s", v)
	}
	if v := sequencer.Next(); v != "v3" {
		t.Errorf("Expected v3 and take %s", v)
	}
	if v := sequencer.Next(); v != "v3" {
		t.Errorf("Expected v3 and take %s", v)
	}
	if v := sequencer.Next(); v != "v3" {
		t.Errorf("Expected v3 and take %s", v)
	}
}
