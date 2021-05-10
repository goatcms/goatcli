package gcliservices

// Arguments return arguments
type Arguments interface {
	ContainerImageDir() string
	TmpContainerImageDir() string
}

// Core provide all nw services and arguments
type Core interface {
	Arguments() Arguments
}
