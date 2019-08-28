package tfunc

import (
	"regexp"
	"sort"
	"strings"
)

// Keys return sub keys for prefix
func Keys(plainmap map[string]string, prefix string) (keys []string) {
	m := map[string]bool{}
	for k := range plainmap {
		if strings.HasPrefix(k, prefix) {
			k = k[len(prefix):]
			i := strings.Index(k, ".")
			if i == -1 {
				i = len(k)
			}
			k = k[:i]
			m[k] = true
		}
	}
	keys = make([]string, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	return keys
}

// Regexp return regexp from string
func Regexp(pattern string) (reg *regexp.Regexp, err error) {
	if reg, err = regexp.Compile(pattern); err != nil {
		return nil, err
	}
	return reg, nil
}

// StrainMap return map filtred by regexp
func StrainMap(data map[string]string, reg *regexp.Regexp) (result map[string]string, err error) {
	result = make(map[string]string)
	for key, value := range data {
		if reg.MatchString(key) {
			result[key] = value
		}
	}
	return result, nil
}
