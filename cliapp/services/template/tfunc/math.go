package tfunc

import (
	"fmt"
	"strconv"
)

// Sum add all values
func Sum(rows ...interface{}) (value int) {
	var (
		i   int
		err error
	)
	for _, row := range rows {
		switch v := row.(type) {
		case int:
			value += v
		case string:
			if i, err = strconv.Atoi(v); err != nil {
				panic(err)
			}
			value += i
		default:
			panic(fmt.Sprint("incorrect template sum type for ", row))
		}
	}
	return value
}

// Sub subtraction two values
func Sub(base, val int) int {
	return base - val
}
