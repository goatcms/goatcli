package tfunc

// UniqueReduce return a set of unique and not empty values
func UniqueReduce(source []string) (result []string) {
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
