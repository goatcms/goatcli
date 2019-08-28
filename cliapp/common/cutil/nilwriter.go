package cutil

// NilWriter is empty wrate (don't persist any data)
type NilWriter struct{}

// NewNilWriter create new NilWriter instance
func NewNilWriter() NilWriter {
	return NilWriter{}
}

// Write writes len(p) bytes from p to the nil data stream
func (w NilWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}

// Close the stream
func (w NilWriter) Close() (err error) {
	return nil
}
