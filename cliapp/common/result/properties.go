package result

import (
	"fmt"
	"regexp"
	"strings"
)

// PropertiesResult is a properties set
type PropertiesResult struct {
	data map[string]string
}

// NewPropertiesResult parse json and return Property array instance
func NewPropertiesResult(data map[string]string) *PropertiesResult {
	return &PropertiesResult{
		data: data,
	}
}

// InjectToString inject properties values to string
// All property key is start with '{{' and end with '}}'
// For example: "property value: {{my_property_key}}"
func (result *PropertiesResult) InjectToString(s string) (out string, err error) {
	pattern := regexp.MustCompile("\\{\\{[A-Za-z_]+\\}\\}")
	indexes := pattern.FindAllStringIndex(s, -1)
	parts := make([]string, len(indexes)*2+1)
	last := 0
	for i, index := range indexes {
		parts[i*2] = s[last:index[0]]
		key := s[index[0]+2 : index[1]-2]
		value, ok := result.data[key]
		if !ok {
			return "", fmt.Errorf("unknow property for '%s' key", key)
		}
		parts[i*2+1] = value
		last = index[1]
	}
	parts[len(indexes)*2] = s[last:len(s)]
	return strings.Join(parts, ""), nil
}

// Get return property value for key
func (result *PropertiesResult) Get(key string) (string, error) {
	value, ok := result.data[key]
	if !ok {
		return "", fmt.Errorf("Unknow property value for %s key", key)
	}
	return value, nil
}
