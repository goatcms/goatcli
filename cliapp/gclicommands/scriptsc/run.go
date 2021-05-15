package scriptsc

import (
	"time"

	"github.com/goatcms/goatcli/cliapp/common"
	"github.com/goatcms/goatcli/cliapp/common/gclivarutil"

	"github.com/goatcms/goatcli/cliapp/gcliservices"
	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/app/gio"
	"github.com/goatcms/goatcore/app/modules/pipelinem/pipservices"
	"github.com/goatcms/goatcore/filesystem"
	"github.com/goatcms/goatcore/varutil/goaterr"
)

// RunScript run script by name
func RunScript(a app.App, ctx app.IOContext) (err error) {
	var (
		deps struct {
			Name            string                       `command:"?$1"`
			ScriptsRunner   gcliservices.ScriptsRunner   `dependency:"ScriptsRunner"`
			GCLIInputs      gcliservices.GCLIInputs      `dependency:"GCLIInputs"`
			GCLIEnvironment gcliservices.GCLIEnvironment `dependency:"GCLIEnvironment"`
			NamespacesUnit  pipservices.NamespacesUnit   `dependency:"PipNamespacesUnit"`
		}
		propertiesData map[string]string
		secretsData    map[string]string
		properties     common.ElasticData
		secrets        common.ElasticData
		appData        gcliservices.ApplicationData
		ctxScope       = ctx.Scope()
		taskManager    pipservices.TasksManager
		scpNamespaces  pipservices.Namespaces

		now          = time.Now()
		cwd          = ctx.IO().CWD()
		fileBasename string

		logsFileWriter    filesystem.Writer
		summaryFileWriter filesystem.Writer
		errorFileWriter   filesystem.Writer

		logsFilePath    string
		summaryFilePath string
		errorFilePath   string

		scriptErr error
	)
	if err = a.DependencyProvider().InjectTo(&deps); err != nil {
		return err
	}
	if err = ctxScope.InjectTo(&deps); err != nil {
		return err
	}
	if deps.Name == "" {
		return goaterr.Errorf(FirstKeyParameterIsRequired)
	}
	// prepare paths
	fileBasename = ".goat/tmp/logs/scripts/" + now.Format("2006-01-02-15:04:05") + "-" + deps.Name
	logsFilePath = fileBasename + ".io.log"
	summaryFilePath = fileBasename + ".summary.log"
	errorFilePath = fileBasename + ".error.log"
	// load namespace
	if scpNamespaces, err = deps.NamespacesUnit.FromScope(ctx.Scope(), defaultNamespace); err != nil {
		return err
	}
	// load variables
	if propertiesData, secretsData, appData, err = deps.GCLIInputs.Inputs(ctx); err != nil {
		return err
	}
	if err = deps.GCLIEnvironment.LoadEnvs(ctxScope, propertiesData, secretsData); err != nil {
		return err
	}
	// run script
	if properties, err = gclivarutil.NewElasticData(propertiesData); err != nil {
		return err
	}
	if secrets, err = gclivarutil.NewElasticData(secretsData); err != nil {
		return err
	}
	if taskManager, err = deps.ScriptsRunner.Run(gcliservices.ScriptsContext{
		Scope:      ctx.Scope(),
		CWD:        ctx.IO().CWD(),
		Namespaces: scpNamespaces,
	}, cwd, deps.Name, properties, secrets, appData); err != nil {
		return err
	}
	if err = taskManager.StatusBroadcast().Add(ctx.IO().Out()); err != nil {
		return err
	}
	// if it is main lvl script
	if scpNamespaces.Task() == "" {
		// prepare log file
		if err = cwd.MkdirAll(".goat/tmp/logs/scripts", filesystem.DefaultUnixDirMode); err != nil {
			return err
		}
		if logsFileWriter, err = cwd.Writer(logsFilePath); err != nil {
			return err
		}
		defer logsFileWriter.Close()
		if err = taskManager.OBroadcast().Add(logsFileWriter); err != nil {
			return err
		}
		// wait for all tasks
		if scriptErr = taskManager.Wait(); scriptErr != nil {
			// save err
			if errorFileWriter, err = cwd.Writer(errorFilePath); err != nil {
				return err
			}
			defer errorFileWriter.Close()
			if jsonErr, ok := scriptErr.(goaterr.JSONError); ok {
				if _, err = errorFileWriter.Write([]byte(jsonErr.ErrorJSON())); err != nil {
					return err
				}
			} else {
				if _, err = errorFileWriter.Write([]byte(scriptErr.Error())); err != nil {
					return err
				}
			}
		}
		// save summary
		if summaryFileWriter, err = cwd.Writer(summaryFilePath); err != nil {
			return err
		}
		defer summaryFileWriter.Close()
		if err = taskManager.Summary(gio.NewOutput(summaryFileWriter)); err != nil {
			return err
		}
		if scriptErr != nil {
			return goaterr.Wrapf(scriptErr, "Script fail (see details at '%s' log files). ", errorFilePath)
		}
	}
	return nil
}
