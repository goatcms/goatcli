package naming

import (
	"regexp"
	"strings"
)

var under = regexp.MustCompile("(^[^A-Z]*|[A-Z]*)([A-Z][^A-Z]+|$)")
var cc = regexp.MustCompile("[\\s\\t!@#$%^&*()_\\+-=\\.]+[a-zA-Z0-9]{1}")

// ToUnderscore convert string to underscore
func ToUnderscore(name string) (result string) {
	var wasBig = false
	for i := 0; i < len(name); i++ {
		c := name[i]
		if c >= 'a' && c <= 'z' {
			result += string(c)
			wasBig = false
			continue
		}
		if !wasBig && len(result) > 0 && !strings.HasSuffix(result, "_") {
			result += "_"
		}
		if c >= 'A' && c <= 'Z' {
			result += strings.ToLower(string(c))
			wasBig = true
		} else {
			wasBig = false
		}
	}
	return result
}

// ToCamelCase convert string to CamelCase
func ToCamelCase(s string) string {
	var a []string
	lastIndex := 0
	for _, sub := range cc.FindAllStringIndex(s, -1) {
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
