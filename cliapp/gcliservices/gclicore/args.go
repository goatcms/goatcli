package gclicore

import (
	"os"

	"github.com/goatcms/goatcli/cliapp/gcliservices"
	"github.com/goatcms/goatcore/dependency"
)

// ArgumentsProvider provider
type ArgumentsProvider struct {
	deps struct {
		ContainerImageDir    string `argument:"?container-image-dir" ,config:"?container.image.dir"`
		TmpContainerImageDir string `argument:"?container-image-tmp" ,config:"?container.image.tmp"`
	}
}

// ArgumentsFactory create new Arguments instance
func ArgumentsFactory(dp dependency.Provider) (in interface{}, err error) {
	instance := &ArgumentsProvider{}
	if err = dp.InjectTo(&instance.deps); err != nil {
		return nil, err
	}
	// prepare ContainerImageDir
	if instance.deps.ContainerImageDir == "" {
		instance.deps.ContainerImageDir = os.Getenv("GCLI_CONTAINER_IMAGE_DIR")
	}
	if instance.deps.ContainerImageDir == "" {
		instance.deps.ContainerImageDir = "/var/lib/virtual-containers"
	}
	// prepare TmpContainerImageDir
	if instance.deps.TmpContainerImageDir == "" {
		instance.deps.TmpContainerImageDir = os.Getenv("GCLI_CONTAINER_IMAGE_TMP")
	}
	if instance.deps.TmpContainerImageDir == "" {
		instance.deps.TmpContainerImageDir = "/tmp/containers"
	}
	return gcliservices.Arguments(instance), nil
}

// ContainerImageDir return image directory
func (args *ArgumentsProvider) ContainerImageDir() string {
	return args.deps.ContainerImageDir
}

// TmpContainerImageDir return image directory
func (args *ArgumentsProvider) TmpContainerImageDir() string {
	return args.deps.TmpContainerImageDir
}
