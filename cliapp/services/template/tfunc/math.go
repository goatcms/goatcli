package tfunc

// Sum add all values
func Sum(rows ...int64) (value int64) {
	for _, row := range rows {
		value += row
	}
	return value
}

// Sum add all values
func Minus(base, val int64) int64 {
	return base - val
}
