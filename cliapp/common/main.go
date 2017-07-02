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
