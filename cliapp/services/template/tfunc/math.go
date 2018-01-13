package tfunc

// Sum add all values
func Sum(rows ...int) (value int) {
	for _, row := range rows {
		value += row
	}
	return value
}

// Sub subtraction two values
func Sub(base, val int) int {
	return base - val
}
