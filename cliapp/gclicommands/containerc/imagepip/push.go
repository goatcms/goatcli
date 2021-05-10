package imagepip

import (
	"strings"

	"github.com/goatcms/goatcli/cliapp/gcliservices"
	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/app/gio"
	"github.com/goatcms/goatcore/app/modules/commonm/commservices/envs"
	"github.com/goatcms/goatcore/app/modules/ocm/ocservices"
	"github.com/goatcms/goatcore/app/scope"
	"github.com/goatcms/goatcore/filesystem"
	"github.com/goatcms/goatcore/varutil/goaterr"
)

const (
	shPushScript = `
set -x
IMAGE_PATH="/tmp/containers/$TMP_NAME/image.tar"
if [ ! -f "$IMAGE_PATH" ]; then
    echo "You must build image befor push"
fi
echo "push... started"
if [ ! -z "$LOGIN" ]; then
	crane auth login -u "$LOGIN" -p "$PASSWORD" "$REGISTRY"
fi
if [ "$TLSVERIFY" == "true" ]; then
	crane push "$IMAGE_PATH" "$REGISTRY/$IMAGE"
else
	crane push --insecure "$IMAGE_PATH" "$REGISTRY/$IMAGE"
fi
echo "push... ok"
`
)

// RunPush run push command
func RunPush(a app.App, ctx app.IOContext) (err error) {
	var (
		deps struct {
			GCLIInputs gcliservices.GCLIInputs `dependency:"GCLIInputs"`
			OCManager  ocservices.Manager      `dependency:"OCManager"`
			Arguments  gcliservices.Arguments  `dependency:"GCLICoreArguments"`
			RootFS     filesystem.Filespace    `filespace:"root"`

			LoginKey    string `command:"?login"`
			PassowrdKey string `command:"?password"`
			Dest        string `command:"?dest"`
			TLSVerify   string `command:"?tls-verify"`
		}
		envs     = envs.NewEnvironments()
		childCtx app.IOContext
		tmpName  string

		login       string
		passowrd    string
		secretsData map[string]string
		ok          bool

		destRegistry string
		destImage    string
	)
	if err = a.DependencyProvider().InjectTo(&deps); err != nil {
		return err
	}
	if err = ctx.Scope().InjectTo(&deps); err != nil {
		return err
	}
	if deps.Dest == "" {
		return goaterr.Errorf("--dest argument is required")
	}
	if deps.TLSVerify == "" {
		deps.TLSVerify = "true"
	}
	if deps.TLSVerify != "true" && deps.TLSVerify != "false" {
		return goaterr.Errorf("--tls-verify allow true or false values only")
	}
	if deps.LoginKey != "" || deps.PassowrdKey != "" {
		if deps.LoginKey == "" {
			return goaterr.Errorf("(--login and --password are required) --login argument is required. Should contains your login key (contains in secrets)")
		}
		if deps.PassowrdKey == "" {
			return goaterr.Errorf("(--login and --password are required) --password argument is required. Should contains your password key (contains in secrets)")
		}
		if _, secretsData, _, err = deps.GCLIInputs.Inputs(ctx); err != nil {
			return err
		}
		if login, ok = secretsData[deps.LoginKey]; !ok {
			return goaterr.Errorf("Login (for key '%s') is not defined", deps.LoginKey)
		}
		if login == "" {
			return goaterr.Errorf("Login (for key '%s') can not be empty", deps.LoginKey)
		}
		if passowrd, ok = secretsData[deps.PassowrdKey]; !ok {
			return goaterr.Errorf("Password (for key '%s') is not defined", deps.PassowrdKey)
		}
		if passowrd == "" {
			return goaterr.Errorf("Password (for key '%s') can not be empty", deps.PassowrdKey)
		}
	}
	pos := strings.Index(deps.Dest, "/")
	if pos == -1 {
		return goaterr.Errorf("Expected destination (--dest) like 'REGISTRY/IMAGE:VERSION' (version is optional) like 'localhost:5000/nginx-template'")
	}
	destRegistry = deps.Dest[:pos]
	destImage = deps.Dest[pos+1:]
	if destRegistry == "" {
		return goaterr.Errorf("Expected destination (--dest) like 'REGISTRY/IMAGE:VERSION'. REGISTRY can not be empty. ")
	}
	if destImage == "" {
		return goaterr.Errorf("Expected destination (--dest) like 'REGISTRY/IMAGE:VERSION'. IMAGE can not be empty. ")
	}
	if tmpName, err = scope.GetString(ctx.Scope(), TmpImageNameKey); err != nil || tmpName == "" {
		return goaterr.Wrapf("ImageName is required. Function must be run into image pipeline.", err)
	}
	if err = goaterr.ToError(goaterr.AppendError(nil,
		envs.Set("REGISTRY", destRegistry),
		envs.Set("IMAGE", destRegistry),
		envs.Set("LOGIN", login),
		envs.Set("PASSWORD", passowrd),
		envs.Set("TLSVERIFY", deps.TLSVerify),
		envs.Set("TMP_NAME", tmpName),
	)); err != nil {
		return err
	}
	childCtx = gio.NewChildIOContext(ctx, gio.ChildIOContextParams{
		IO: gio.IOParams{
			In: gio.NewInput(strings.NewReader(shPushScript)),
		},
	})
	defer childCtx.Scope().Close()
	return deps.OCManager.Run(ocservices.Container{
		IO:         childCtx.IO(),
		Image:      "gcr.io/go-containerregistry/crane:debug",
		WorkDir:    "/cwd",
		Entrypoint: "sh",
		Envs:       envs,
		FSVolumes: map[string]ocservices.FSVolume{
			"/tmp/containers": {
				Filespace: deps.RootFS,
				Path:      deps.Arguments.TmpContainerImageDir(),
			},
		},
	})
}
