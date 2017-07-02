package bcontext

import "github.com/goatcms/goatcore/filesystem"

// Options contains BuildContext options and private data
type Options struct {
	From, To   string
	FS         filesystem.Filespace
	Data       map[string]string
	Hash       string
	Properties PropertieOptions
}

// PropertieOptions contains properties
type PropertieOptions struct {
	Build   map[string]string
	Project map[string]string
	Module  map[string]string
}
