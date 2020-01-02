package cloner

import (
	"github.com/goatcms/goatcli/cliapp/common/config"
	"github.com/goatcms/goatcore/filesystem"
	"github.com/goatcms/goatcore/filesystem/fshelper"
)

const (
	// ReplaceConfigFile is path to replaces config file
	ReplaceConfigFile = ".goat/replaces.def.json"
)

func copy(sourcefs, destfs filesystem.Filespace, subPath string, replaces []*config.Replace) (err error) {
	var content []byte
	for i, replace := range replaces {
		if replace.Pattern.MatchString(subPath) {
			content, err = sourcefs.ReadFile(subPath)
			if err != nil {
				return err
			}
			content = replaceLoop(subPath, content, replaces[i:])
			return destfs.WriteFile(subPath, content, 0766)
		}
	}
	return fshelper.StreamCopy(sourcefs, destfs, subPath)
}

func replaceLoop(subPath string, content []byte, replaces []*config.Replace) []byte {
	for _, replace := range replaces {
		if replace.Pattern.MatchString(subPath) {
			content = replace.From.ReplaceAll(content, []byte(replace.To))
		}
	}
	return content
}
