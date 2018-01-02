package tfunc

import (
	"regexp"
	"strings"
)

var camel = regexp.MustCompile("(^[^A-Z]*|[A-Z]*)([A-Z][^A-Z]+|$)")
var under = regexp.MustCompile("[\\s\\t!@#$%^&*()_\\+-=.]+[a-zA-Z0-9]{1}")

// ToUnderscore convert string to underscore
func ToUnderscore(s string) string {
	var a []string
	for _, sub := range camel.FindAllStringSubmatch(s, -1) {
		if sub[1] != "" {
			a = append(a, sub[1])
		}
		if sub[2] != "" {
			a = append(a, sub[2])
		}
	}
	return strings.ToLower(strings.Join(a, "_"))
}

// ToCamelCase convert string to CamelCase
func ToCamelCase(s string) string {
	var a []string
	lastIndex := 0
	for _, sub := range under.FindAllStringIndex(s, -1) {
		a = append(a, s[lastIndex:sub[0]])
		a = append(a, strings.ToUpper(s[sub[1]-1:sub[1]]))
		lastIndex = sub[1]
	}
	a = append(a, s[lastIndex:])
	return strings.Join(a, "")
}

// ToCamelCaseLF convert string to CamelCase (with first letter lower)
func ToCamelCaseLF(s string) string {
	return ToLowerFirst(ToCamelCase(s))
}

// ToCamelCaseUF convert string to CamelCase (with first letter upper)
func ToCamelCaseUF(s string) string {
	return ToUpperFirst(ToCamelCase(s))
}

// ToLowerFirst make first letter lower
func ToLowerFirst(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToLower(s[:1]) + s[1:]
}

// ToUpperFirst make first letter upper
func ToUpperFirst(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToUpper(s[:1]) + s[1:]
}
