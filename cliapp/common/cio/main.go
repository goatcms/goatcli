package cio

import "regexp"

var (
	numericReg = regexp.MustCompile("^[0-9]+$")
	alphaReg   = regexp.MustCompile("^[A-Za-z]+$")
	alnumReg   = regexp.MustCompile("^[A-Za-z0-9]+$")
)
