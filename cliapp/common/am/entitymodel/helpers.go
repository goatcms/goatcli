package entitymodel

import (
	"regexp"
	"strings"
)

func arrayValue(data map[string]string, key string, defaultValue []string) (v []string) {
	var s string
	if s = data[key]; s == "" {
		return defaultValue
	}
	rx := regexp.MustCompile(`\s`)
	s = rx.ReplaceAllString(s, "")
	return strings.Split(s, ",")
}

func splitName(fullPath string) (path, name string) {
	if pos := strings.LastIndex(fullPath, "."); pos != -1 {
		return fullPath[:pos], fullPath[pos:]
	}
	return "", fullPath
}
