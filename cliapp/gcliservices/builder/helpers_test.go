package builder

import (
	"bytes"

	"github.com/goatcms/goatcli/cliapp/common"
	"github.com/goatcms/goatcli/cliapp/common/am"
	"github.com/goatcms/goatcli/cliapp/common/gclivarutil"
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

func renderFile(fs filesystem.Filespace, deps RenderFileDeps, data, plainProperties, plainSecrets map[string]string) (err error) {
	var (
		ctx        app.IOContext
		appData    gcliservices.ApplicationData
		properties common.ElasticData
		secrets    common.ElasticData
	)
	if ctx, err = newEmptyIOContext(); err != nil {
		return
	}
	if appData, err = am.NewApplicationData(data); err != nil {
		return
	}
	if properties, err = gclivarutil.NewElasticData(plainProperties); err != nil {
		return
	}
	if secrets, err = gclivarutil.NewElasticData(plainSecrets); err != nil {
		return
	}
	if err = deps.BuilderService.Build(ctx, fs, appData, properties, secrets); err != nil {
		return err
	}
	if err = ctx.Scope().Wait(); err != nil {
		return err
	}
	return ctx.Scope().Trigger(gcliservices.BuildCommitevent, nil)
}

func newEmptyIOContext() (ctx app.IOContext, err error) {
	var (
		scp = scope.New(scope.Params{})
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
