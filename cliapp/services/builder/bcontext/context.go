package bcontext

// BuildContext contains helpers and context data
type BuildContext struct {
	From, To    string
	Out         *Out
	Data        map[string]string
	Filesystem  *Filesystem
	Propertsies PropertieOptions
}

// NewBuildContext create new build context
func NewBuildContext(options *Options) *BuildContext {
	return &BuildContext{
		From: options.From,
		To:   options.To,
		Out: &Out{
			hash:       options.Hash,
			isFileOpen: false,
		},
		Data: options.Data,
		Filesystem: &Filesystem{
			fs: options.FS,
		},
		Propertsies: options.Properties,
	}
}
