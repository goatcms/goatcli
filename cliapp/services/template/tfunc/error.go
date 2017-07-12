package tfunc

import "fmt"

// ToError return string as error
func ToError(s string) error {
	return fmt.Errorf(s)
}
