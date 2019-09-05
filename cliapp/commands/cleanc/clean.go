package cleanc

import (
	"os"
	"strings"
	"time"

	"github.com/goatcms/goatcli/cliapp/services"
	"github.com/goatcms/goatcli/cliapp/services/vcs"
	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/filesystem"
	"github.com/goatcms/goatcore/varutil/goaterr"
)

// RunClean run clean command
func RunClean(a app.App, ctxScope app.Scope) (err error) {
	var (
		deps struct {
			AllFlag string `argument:"?interactive" ,command:"?interactive"`

			CurrentFS filesystem.Filespace `filespace:"current"`

			VCSService services.VCSService `dependency:"VCSService"`
			Input      app.Input           `dependency:"InputService"`
			Output     app.Output          `dependency:"OutputService"`
		}
		allMode  bool
		vcsData  services.VCSData
		saved    = vcs.NewGeneratedFiles(true)
		info     os.FileInfo
		response string
	)
	if err = goaterr.ToErrors(goaterr.AppendError(nil,
		a.DependencyProvider().InjectTo(&deps),
		ctxScope.InjectTo(&deps))); err != nil {
		return err
	}
	allMode = strings.ToLower(deps.AllFlag) == "y"
	// load data
	if vcsData, err = deps.VCSService.ReadDataFromFS(deps.CurrentFS); err != nil {
		return err
	}
	for _, row := range vcsData.VCSGeneratedFiles().All() {
		if !deps.CurrentFS.IsFile(row.Path) || vcsData.VCSIgnoredFiles().ContainsPath(row.Path) {
			continue
		}
		if info, err = deps.CurrentFS.Lstat(row.Path); err != nil {
			return err
		}
		generatedTime := row.ModTime.Format(time.RFC3339)
		filespaceTime := info.ModTime().Format(time.RFC3339)
		if allMode == false && (generatedTime != filespaceTime) {
			response = ""
			for response != "y" && response != "n" {
				deps.Output.Printf("File %s (generated at %s) was modified by user (at %s). Do you want remove it? (y/n)\n", row.Path, generatedTime, filespaceTime)
				if response, err = deps.Input.ReadLine(); err != nil {
					return err
				}
				response = strings.ToLower(response)
			}
			if response == "n" {
				saved.Add(row)
				deps.Output.Printf("Ignored %s\n", row.Path)
				response = ""
				for response != "y" && response != "n" {
					deps.Output.Printf("Do you want add %s to persist files? (y/n)\n", row.Path)
					if response, err = deps.Input.ReadLine(); err != nil {
						return err
					}
					response = strings.ToLower(response)
				}
				if response == "y" {
					vcsData.VCSIgnoredFiles().AddPath(row.Path)
					continue
				}
				continue
			}
		}
		if err = deps.CurrentFS.Remove(row.Path); err != nil {
			return err
		}
	}
	vcsData = vcs.NewData(saved, vcsData.VCSIgnoredFiles())
	if err = deps.VCSService.WriteDataToFS(deps.CurrentFS, vcsData); err != nil {
		return err
	}
	deps.Output.Printf("cleaned\n")
	return nil
}
