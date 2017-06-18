package common

// StringInjector inject data / values to string
type StringInjector interface {
	InjectToString(string) (string, error)
}
