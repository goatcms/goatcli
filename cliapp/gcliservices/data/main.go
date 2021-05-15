package data

import (
	"os"
	"strings"

	"github.com/goatcms/goatcore/filesystem"
	"github.com/goatcms/goatcore/varutil/goaterr"
	"github.com/goatcms/goatcore/varutil/plainmap"
)

const (
	// DataDefPath is path to data definitions
	DataDefPath = ".goat/data.def.json"
	// BaseDataDefPath is base directory for data files
	BaseDataDefPath = ".goat/data.def/"
	// DataDefSuffix is definition file extensions
	DataDefSuffix = ".json"
	// BaseDataPath is base directory for data files
	BaseDataPath = ".goat/data/"
)

func mergeMap(dest, src map[string]string) (err error) {
	for key, value := range src {
		if old, ok := dest[key]; ok {
			return goaterr.Errorf("duplicate key %s in source (%s) and destination (%s) map ", key, value, old)
		}
		dest[key] = value
	}
	return nil
}

func readDataFromFS(data map[string]string, fs filesystem.Filespace, path string) (err error) {
	var (
		nodes    []os.FileInfo
		json     []byte
		filedata map[string]string
	)
	if nodes, err = fs.ReadDir(path); err != nil {
		return err
	}
	for _, node := range nodes {
		if node.IsDir() {
			if err = readDataFromFS(data, fs, path+node.Name()+"/"); err != nil {
				return err
			}
		} else {
			nodeName := node.Name()
			if !strings.HasSuffix(nodeName, ".json") {
				// skip no *.json file
				continue
			}
			if json, err = fs.ReadFile(path + nodeName); err != nil {
				return err
			}
			if filedata, err = plainmap.JSONToPlainStringMap(json); err != nil {
				return err
			}
			if err = mergeMap(data, filedata); err != nil {
				return err
			}
		}
	}
	return nil
}
