package simpletf

import (
	"math/rand"
	"sort"

	"github.com/goatcms/goatcore/varutil"
)

// Union return a set of unique and not empty values from collections
func Union(sources ...[]string) (result []string) {
	result = []string{}
	for _, source := range sources {
	UniqueReduceLoop:
		for _, value := range source {
			if value == "" {
				continue
			}
			for _, c := range result {
				if c == value {
					continue UniqueReduceLoop
				}
			}
			result = append(result, value)
		}
	}
	return result
}

// Unique return a set of unique and not empty values
func Unique(source []string) (result []string) {
	result = []string{}
UniqueReduceLoop:
	for _, value := range source {
		if value == "" {
			continue
		}
		for _, c := range result {
			if c == value {
				continue UniqueReduceLoop
			}
		}
		result = append(result, value)
	}
	return result
}

// Except return a set of unique values except values from other sources
func Except(base, except []string) (result []string) {
	except = Unique(except)
	base = Unique(base)
	result = []string{}
	for _, value := range base {
		if !varutil.IsArrContainStr(except, value) {
			result = append(result, value)
		}
	}
	return result
}

// Intersect return a set of unique values contains in all tables
func Intersect(base, intersect []string) (result []string) {
	intersect = Unique(intersect)
	base = Unique(base)
	result = []string{}
	for _, value := range base {
		if varutil.IsArrContainStr(intersect, value) {
			result = append(result, value)
		}
	}
	return result
}

// Sort order strings array
func Sort(base []string) (result []string) {
	result = append(base)
	sort.Strings(result)
	return result
}

// RandomValue return random value from a collection
func RandomValue(source ...string) string {
	return source[rand.Intn(len(source))]
}
