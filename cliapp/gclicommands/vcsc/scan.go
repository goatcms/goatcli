package vcsc

import (
	"os"
	"strings"
	"time"

	"github.com/goatcms/goatcli/cliapp/gcliservices"
	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/filesystem"
	"github.com/goatcms/goatcore/varutil/goaterr"
)

// RunScan run vcs:scan command. It is looking for changes in generated files and add it do ignore if required
func RunScan(a app.App, ctx app.IOContext) (err error) {
	var (
		deps struct {
			CurrentFS       filesystem.Filespace    `filespace:"current"`
			InteractiveFlag string                  `argument:"?interactive" ,command:"?interactive"`
			VCSService      gcliservices.VCSService `dependency:"VCSService"`
		}
		vcsData         gcliservices.VCSData
		ignored         gcliservices.VCSIgnoredFiles
		interactiveMode bool
		info            os.FileInfo
	)
	if err = goaterr.ToErrors(goaterr.AppendError(nil,
		a.DependencyProvider().InjectTo(&deps),
		ctx.Scope().InjectTo(&deps))); err != nil {
		return err
	}
	interactiveMode = strings.ToLower(deps.InteractiveFlag) != "false"
	ctx.IO().Out().Printf("Scan generated files for changes...\n")
	if vcsData, err = deps.VCSService.ReadDataFromFS(deps.CurrentFS); err != nil {
		return err
	}
	ignored = vcsData.VCSIgnoredFiles()
	for _, row := range vcsData.VCSGeneratedFiles().All() {
		if !deps.CurrentFS.IsFile(row.Path) || vcsData.VCSIgnoredFiles().ContainsPath(row.Path) {
			continue
		}
		if info, err = deps.CurrentFS.Lstat(row.Path); err != nil {
			return err
		}
		generatedTime := row.ModTime.Format(time.RFC3339)
		filespaceTime := info.ModTime().Format(time.RFC3339)
		if generatedTime != filespaceTime {
			if !interactiveMode {
				ignored.AddPath(row.Path)
				continue
			}
			response := ""
			for response != "y" && response != "n" {
				ctx.IO().Out().Printf("File %s (generated at %s) was modified (at %s). Do you want persist it (by add to .goat/vcs/ignored)? (y/n)\n", row.Path, generatedTime, filespaceTime)
				if response, err = ctx.IO().In().ReadLine(); err != nil {
					return err
				}
				response = strings.ToLower(response)
			}
			if response == "y" {
				ignored.AddPath(row.Path)
				continue
			}
		}
	}
	return deps.VCSService.WriteDataToFS(deps.CurrentFS, vcsData)
}
