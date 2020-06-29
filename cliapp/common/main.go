package common

// StringInjector inject data / values to string
type StringInjector interface {
	InjectToString(string) (string, error)
}

// PropertiesResult provide properties data
type PropertiesResult interface {
	StringInjector
	Get(string) (string, error)
}

// ElasticData provide plain data and tree data structures
type ElasticData struct {
	Plain map[string]string
	Tree  map[string]interface{}
}
