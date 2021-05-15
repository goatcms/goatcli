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
	shBuildScript = `
set -x
echo "build... started"
OUT_DIR="/tmp/containers/$TMP_NAME"
OUT_FILE="$OUT_DIR/image.tar"
DOCKERFILE_DIR="/tmp/$TMP_NAME"
DOCKERFILE_PATH="$DOCKERFILE_DIR/Dockerfile"
mkdir -p "$OUT_DIR"
mkdir -p "$DOCKERFILE_DIR"
echo "$STEPS" > "$DOCKERFILE_PATH"
# --cache=true --cache-dir="/cache" 
mkdir -p "/cache/test"
/kaniko/executor --context "/cwd" --dockerfile "$DOCKERFILE_PATH" --cache-dir="/cache" --no-push --destination "$TMP_NAME" --tarPath "$OUT_FILE"
if [ ! -f "$OUT_FILE" ]; then 
	echo "Expected $OUT_FILE file."
	exit 1
fi
echo "build... ok"
`
)

// RunBuild run build command
func RunBuild(a app.App, ctx app.IOContext) (err error) {
	var (
		deps struct {
			OCManager ocservices.Manager     `dependency:"OCManager"`
			Arguments gcliservices.Arguments `dependency:"GCLICoreArguments"`
			RootFS    filesystem.Filespace   `filespace:"root"`

			Steps string `command:"?steps"`
		}
		envs     = envs.NewEnvironments()
		childCtx app.IOContext
		tmpName  string
	)
	if err = a.DependencyProvider().InjectTo(&deps); err != nil {
		return err
	}
	if err = ctx.Scope().InjectTo(&deps); err != nil {
		return err
	}
	if tmpName, err = scope.GetString(ctx.Scope(), TmpImageNameKey); err != nil || tmpName == "" {
		return goaterr.Wrapf(err, "function must be run into image pipeline")
	}
	if deps.Steps == "" {
		return goaterr.Errorf("steps argument is required")
	}
	if err = goaterr.ToError(goaterr.AppendError(nil,
		envs.Set("STEPS", deps.Steps),
		envs.Set("TMP_NAME", tmpName),
	)); err != nil {
		return err
	}
	childCtx = gio.NewChildIOContext(ctx, gio.ChildIOContextParams{
		IO: gio.IOParams{
			In: gio.NewInput(strings.NewReader(shBuildScript)),
		},
	})
	defer childCtx.Scope().Close()
	return deps.OCManager.Run(ocservices.Container{
		IO:         childCtx.IO(),
		Image:      "gcr.io/kaniko-project/executor:debug",
		WorkDir:    "/cwd",
		Entrypoint: "sh",
		Envs:       envs,
		FSVolumes: map[string]ocservices.FSVolume{
			"/cwd": {
				Filespace: ctx.IO().CWD(),
			},
			"/cache": {
				Filespace: deps.RootFS,
				Path:      deps.Arguments.ContainerImageDir(),
			},
			"/tmp/containers": {
				Filespace: deps.RootFS,
				Path:      deps.Arguments.TmpContainerImageDir(),
			},
		},
	})
}
