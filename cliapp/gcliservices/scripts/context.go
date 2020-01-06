package scripts

// Context is script context data
type Context struct {
	PlainData  map[string]string
	Properties TaskProperties
	AM         interface{}
}

// TaskProperties contains task properties
type TaskProperties struct {
	Project map[string]string
	Secrets map[string]string
}
