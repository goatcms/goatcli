package simpletf

// Sequencer is a helper to iterate over vlues
type Sequencer struct {
	values []string
	index  int
}

// NewSequencer return new Sequencer instance
func NewSequencer(values ...string) *Sequencer {
	return &Sequencer{
		values: values,
		index:  0,
	}
}

// HasNext return true if has more elements
func (sequencer *Sequencer) HasNext() bool {
	return sequencer.index < len(sequencer.values)
}

// Next increment index and return next value (or last value if it is finich)
func (sequencer *Sequencer) Next() string {
	current := sequencer.value()
	sequencer.index++
	return current
}

// Value return current value
func (sequencer *Sequencer) value() string {
	if sequencer.index >= len(sequencer.values) {
		return sequencer.values[len(sequencer.values)-1]
	}
	return sequencer.values[sequencer.index]
}
