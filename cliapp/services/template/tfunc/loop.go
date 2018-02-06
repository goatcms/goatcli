package tfunc

// reapetNTimeEmptyStruct is struct who don't require allocate memmory
type reapetNTimeEmptyStruct struct{}

// RepeatNTimes return n elements array for template range
func RepeatNTimes(n int) []reapetNTimeEmptyStruct {
	return make([]reapetNTimeEmptyStruct, n)
}
