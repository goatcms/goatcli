package bcontext

import "regexp"

// Data is build data access provider
type Data struct {
	data map[string]string
}

// All return data map (all plainmap keys)
func (d *Data) All() map[string]string {
	return d.data
}

// Filter return keys match to regexp
func (d *Data) Filter(r string) (result map[string]string, err error) {
	var reg *regexp.Regexp
	result = make(map[string]string)
	if reg, err = regexp.Compile(r); err != nil {
		return nil, err
	}
	for key, value := range d.data {
		if reg.MatchString(key) {
			result[key] = value
		}
	}
	return result, nil
}
