package containerc

import (
	"strconv"
	"strings"
	"testing"

	"github.com/goatcms/goatcore/app/scope"
	"github.com/goatcms/goatcore/testbase"

	"github.com/goatcms/goatcore/goatnet"

	"github.com/goatcms/goatcore/app/gio"

	"github.com/goatcms/goatcore/varutil/goaterr"

	"github.com/goatcms/goatcore/filesystem"

	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/app/mockupapp"
	"github.com/goatcms/goatcore/app/modules/ocm/ocservices"
)

func TestBuildStory(t *testing.T) {
	t.Parallel()
	var (
		err         error
		mapp        *mockupapp.App
		bootstraper app.Bootstrap
		rootFS      filesystem.Filespace
		deps        struct {
			OCManager ocservices.Manager `dependency:"OCManager"`
		}
		registryPort int
		testScope    = scope.NewScope(scope.Params{})
	)
	defer testScope.Kill()
	if _, err = testbase.LoadDockerTestConfig(); err != nil {
		t.Skip(err.Error())
		return
	}
	if registryPort, err = goatnet.GetFreePort(); err != nil {
		t.Error(err)
		return
	}
	if rootFS, err = newBaseFS(t, "TestContainerImageBuildStory"); err != nil {
		t.Error(err)
		return
	}
	if mapp, bootstraper, err = newApp(mockupapp.MockupOptions{
		Args: []string{
			`goatcli`,
			`terminal`,
			`--strict=true`,
		},
		Input: strings.NewReader(`
			container:image --pip=<<ENDPIP
				build --steps="FROM alpine:latest"
				push --tls-verify=false --dest="host.docker.internal:` + strconv.Itoa(registryPort) + `/testuser/testrepo"
ENDPIP`),
		RootFilespace: rootFS,
	}); err != nil {
		t.Error(err)
		return
	}
	if err = mapp.DependencyProvider().InjectTo(&deps); err != nil {
		t.Error(err)
		return
	}
	fs := mapp.RootFilespace()
	if err = fs.MkdirAll(".goat", filesystem.DefaultUnixDirMode); err != nil {
		t.Error(err)
		return
	}
	// run registry
	go func() {
		if err = deps.OCManager.Run(ocservices.Container{
			IO: gio.NewIO(gio.IOParams{
				In:  gio.NewNilInput(),
				Out: gio.NewNilOutput(),
				Err: gio.NewNilOutput(),
				CWD: fs,
			}),
			Image: "registry:2",
			Ports: map[int]int{
				5000: registryPort,
			},
			Scope: testScope,
		}); err != nil {
			mapp.AppScope().AppendError(err)
		}
	}()
	// test
	if err = bootstraper.Run(); err != nil {
		t.Error(goaterr.Wrapf("Error:\n%s\nStdOut:\n%s\nStdErr:\n%s\n", err, err.Error(), mapp.OutputBuffer().String(), mapp.ErrorBuffer().String()))
		return
	}
	if err = mapp.AppScope().Wait(); err != nil {
		t.Error(goaterr.Wrapf("Error:\n%s\nStdOut:\n%s\nStdErr:\n%s\n", err, err.Error(), mapp.OutputBuffer().String(), mapp.ErrorBuffer().String()))
		return
	}
	result := mapp.OutputBuffer().String()
	if strings.Index(result, "build... ok") == -1 {
		t.Errorf("expected 'build... ok' and take: \n%s", result)
		return
	}
	if strings.Index(result, "login... ok") != -1 {
		t.Errorf("unexpected 'login... ok' (login only when --login parameter is allow) and take: \n%s", result)
		return
	}
	if strings.Index(result, "push... ok") == -1 {
		t.Errorf("expected 'push... ok' and take: \n%s", result)
		return
	}
}
