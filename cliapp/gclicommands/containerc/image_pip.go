package containerc

import (
	"strings"
	"time"

	"github.com/goatcms/goatcli/cliapp/gclicommands/containerc/imagepip"
	"github.com/goatcms/goatcli/cliapp/gcliservices"
	"github.com/goatcms/goatcore/filesystem"

	"github.com/goatcms/goatcore/varutil"

	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/app/gio"
	"github.com/goatcms/goatcore/app/modules"
	"github.com/goatcms/goatcore/app/scope"
	"github.com/goatcms/goatcore/varutil/goaterr"
)

// RunContainerImagePip create new container context and run container pipeline (run container:image command)
func RunContainerImagePip(a app.App, ctx app.IOContext) (err error) {
	var (
		deps struct {
			Arguments gcliservices.Arguments `dependency:"GCLICoreArguments"`
			RootFS    filesystem.Filespace   `filespace:"root"`

			PIP string `command:"?pip"`
		}
		childCtx app.IOContext
		tmpName  string
		terminal modules.Terminal
	)
	// valid
	if err = a.DependencyProvider().InjectTo(&deps); err != nil {
		return err
	}
	if err = ctx.Scope().InjectTo(&deps); err != nil {
		return err
	}
	if deps.PIP == "" {
		return goaterr.Errorf("--pip argument is required")
	}
	if tmpName, err = scope.GetString(ctx.Scope(), imagepip.TmpImageNameKey); err == nil && tmpName != "" {
		return goaterr.Errorf("You can not run container:image in another container pipeline")
	}
	// prepare image name
	cTime := time.Now()
	tmpName = "gi" + cTime.Format("20060102150405") + varutil.RandString(10, varutil.LowerAlphaNumericBytes)
	// prepare scope
	childCtx = gio.NewChildIOContext(ctx, gio.ChildIOContextParams{
		Scope: scope.ChildParams{
			DataScope: scope.NewDataScope(map[interface{}]interface{}{
				imagepip.TmpImageNameKey: tmpName,
			}),
		},
		IO: gio.IOParams{
			In: gio.NewInput(strings.NewReader(deps.PIP)),
		},
	})
	defer func() {
		childCtx.Scope().Close()
		deps.RootFS.RemoveAll(deps.Arguments.TmpContainerImageDir() + "/" + tmpName)
	}()
	// prepare and run sub-terminal
	if terminal, err = newTerminal(a); err != nil {
		return err
	}
	return terminal.RunLoop(childCtx, "")
}
