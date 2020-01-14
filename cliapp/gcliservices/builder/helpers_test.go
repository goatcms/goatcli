package builder

import (
	"bytes"

	"github.com/goatcms/goatcli/cliapp/common/am"
	"github.com/goatcms/goatcli/cliapp/gcliservices"
	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/app/gio"
	"github.com/goatcms/goatcore/app/scope"
	"github.com/goatcms/goatcore/filesystem"
	"github.com/goatcms/goatcore/filesystem/filespace/memfs"
)

type RenderFileDeps struct {
	BuilderService gcliservices.BuilderService `dependency:"BuilderService"`
}

func renderFile(fs filesystem.Filespace, deps RenderFileDeps, data, properties, secrets map[string]string) (err error) {
	var (
		ctx app.IOContext
	)
	if ctx, err = newEmptyIOContext(); err != nil {
		return err
	}
	appData := am.NewApplicationData(data)
	if err = deps.BuilderService.Build(ctx, fs, appData, properties, secrets); err != nil {
		return err
	}
	if err = ctx.Scope().Wait(); err != nil {
		return err
	}
	return ctx.Scope().Trigger(app.CommitEvent, nil)
}

func newEmptyIOContext() (ctx app.IOContext, err error) {
	var (
		scp = scope.NewScope(scope.Params{})
		io  app.IO
		cwd filesystem.Filespace
		in  = gio.NewInput(new(bytes.Buffer))
		out = gio.NewOutput(new(bytes.Buffer))
	)
	if cwd, err = memfs.NewFilespace(); err != nil {
		return nil, err
	}
	io = gio.NewIO(gio.IOParams{
		In:  in,
		Out: out,
		Err: out,
		CWD: cwd,
	})
	return gio.NewIOContext(scp, io), nil
}
