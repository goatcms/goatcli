package gclicore

import "github.com/goatcms/goatcli/cliapp/gcliservices"

// ArgumentsMockParams provide params
type ArgumentsMockParams struct {
	ContainerImageDir    string
	TmpContainerImageDir string
}

// ArgumentsMock provider
type ArgumentsMock struct {
	params ArgumentsMockParams
}

// NewArgumentsMock create new ArgumentsMock instance
func NewArgumentsMock(params ArgumentsMockParams) gcliservices.Arguments {
	return &ArgumentsMock{
		params: params,
	}
}

// ContainerImageDir return image directory (for cache remote images)
func (args *ArgumentsMock) ContainerImageDir() string {
	return args.params.ContainerImageDir
}

// TmpContainerImageDir return temporary image directory
func (args *ArgumentsMock) TmpContainerImageDir() string {
	return args.params.TmpContainerImageDir
}
