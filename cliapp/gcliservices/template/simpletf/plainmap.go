package simpletf

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

// FindRow row match by prefix, pattern, suffix and value. Key value must be equals to expectedValue. Return key without prefix and suffix.
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

// FindRows match all rows by prefix, pattern, suffix and value. Key value must be equals to expectedValue. Return keys without prefix and suffix.
func FindRows(prefix, pattern, suffix, expectedValue string, source map[string]string) (result []string) {
	var reg *regexp.Regexp
	result = []string{}
	reg = regexp.MustCompile(pattern)
	for key, value := range source {
		if !strings.HasPrefix(key, prefix) || !strings.HasSuffix(key, suffix) || value != expectedValue {
			continue
		}
		newkey := key[len(prefix) : len(key)-len(suffix)]
		if !reg.MatchString(newkey) {
			continue
		}
		result = append(result, newkey)
	}
	return result
}

// SubMap create new map from new one wand replace baseKey with newBaseKey
func SubMap(baseKey, newBaseKey string, source map[string]string) (result map[string]string) {
	result = map[string]string{}
	for key, value := range source {
		if !strings.HasPrefix(key, baseKey) {
			continue
		}
		newkey := newBaseKey + key[len(baseKey):len(key)]
		result[newkey] = value
	}
	return result
}
