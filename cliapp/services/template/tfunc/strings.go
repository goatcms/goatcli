package tfunc

import "strings"

// InjectValues replace values in string
func InjectValues(str string, values map[string]string) string {
	for key, value := range values {
		str = strings.Replace(str, "{{$"+key+"}}", value, -1)
	}
	return str
}

// Replace returns a copy of the string s with the first n
// non-overlapping instances of old replaced by new.
func Replace(s, old, new string) string {
	return strings.Replace(s, old, new, -1)
}
