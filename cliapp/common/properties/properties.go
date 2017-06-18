package properties

import (
	"fmt"
	"regexp"
	"strings"
)

// Properties is a properties set
type Properties struct {
	data map[string]string
}

// NewProperties parse json and return Property array instance
func NewProperties(data map[string]string) *Properties {
	return &Properties{
		data: data,
	}
}

// InjectToString inject properties values to string
// All property key is start with '{{' and end with '}}'
// For example: "property value: {{my_property_key}}"
func (properties *Properties) InjectToString(s string) (out string, err error) {
	pattern := regexp.MustCompile("\\{\\{[A-Za-z_]+\\}\\}")
	indexes := pattern.FindAllStringIndex(s, -1)
	parts := make([]string, len(indexes)*2+1)
	last := 0
	for i, index := range indexes {
		parts[i*2] = s[last:index[0]]
		key := s[index[0]+2 : index[1]-2]
		value, ok := properties.data[key]
		if !ok {
			return "", fmt.Errorf("unknow property for '%s' key", key)
		}
		parts[i*2+1] = value
		last = index[1]
	}
	parts[len(indexes)*2] = s[last:len(s)]
	return strings.Join(parts, ""), nil
}
