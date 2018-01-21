package tfunc

import (
	"regexp"
	"strings"
)

// ValuesFor return values for all keys matches to pattern
func ValuesFor(keyPattern string, source map[string]string) (result []string) {
	var reg *regexp.Regexp
	reg = regexp.MustCompile(keyPattern)
	result = []string{}
	for key, value := range source {
		if !reg.MatchString(key) {
			continue
		}
		result = append(result, value)
	}
	return result
}

func FindRow(prefix, pattern, suffix, expectedValue string, source map[string]string) (result string) {
	var reg *regexp.Regexp
	reg = regexp.MustCompile(pattern)
	for key, value := range source {
		if !strings.HasPrefix(key, prefix) || !strings.HasSuffix(key, suffix) || value != expectedValue {
			continue
		}
		result = key[len(prefix) : len(key)-len(suffix)]
		if !reg.MatchString(result) {
			continue
		}
		return result
	}
	return ""
}

func FindRows(prefix, pattern, suffix, expectedValue string, source map[string]string) (result []string) {
	var reg *regexp.Regexp
	result = []string{}
	reg = regexp.MustCompile(pattern)
	for key, value := range source {
		if !strings.HasPrefix(key, prefix) || !strings.HasSuffix(key, suffix) || value != expectedValue {
			continue
		}
		key := key[len(prefix) : len(key)-len(suffix)]
		if !reg.MatchString(key) {
			continue
		}
		result = append(result, key)
	}
	return result
}

func SubMap(baseKey, newBaseKey string, source map[string]string) (result map[string]string) {
	result = map[string]string{}
	for key, value := range source {
		if !strings.HasPrefix(key, baseKey) {
			continue
		}
		key := newBaseKey + key[len(baseKey):len(key)]
		result[key] = value
	}
	return result
}
