package buildc

import (
	"os"

	"github.com/goatcms/goatcore/filesystem"
)

func pathsList(basePath string, fs filesystem.Filespace) (paths []string, err error) {
	var (
		fileInfos []os.FileInfo
		subPaths  []string
	)
	if fileInfos, err = fs.ReadDir(basePath); err != nil {
		return nil, err
	}
	for _, info := range fileInfos {
		p := basePath + "/" + info.Name()
		if info.IsDir() {
			if subPaths, err = pathsList(p, fs); err != nil {
				return nil, err
			}
			paths = append(paths, subPaths...)
		} else {
			paths = append(paths, p)
		}
	}
	return paths, nil
}
