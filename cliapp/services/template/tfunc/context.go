package tfunc

import "github.com/goatcms/goatcli/cliapp/services/builder/bcontext"

// NewContext build new build context
func NewContext(root *bcontext.BuildContext, from, to string) *bcontext.BuildContext {
	return &bcontext.BuildContext{
		From:       from,
		To:         to,
		Out:        root.Out,
		Data:       root.Data,
		Filesystem: root.Filesystem,
		Properties: root.Properties,
	}
}
