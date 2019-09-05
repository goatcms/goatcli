package tfunc

import "fmt"

// ToError return string as error
func ToError(s string) (string, error) {
	return s, fmt.Errorf(s)
}
