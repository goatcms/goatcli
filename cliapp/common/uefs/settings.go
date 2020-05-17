package uefs

import (
	"os"
	"path/filepath"

	"github.com/goatcms/goatcore/filesystem"
	"github.com/goatcms/goatcore/varutil"
)

// Settings provide host keystore
type Settings struct {
	Salt1 string `json:"salt1"`
	Salt2 string `json:"salt2"`
}

// NewSettings create new Settings instance
func NewSettings() Settings {
	return Settings{
		Salt1: varutil.RandString(45, varutil.StrongBytes),
		Salt2: varutil.RandString(45, varutil.StrongBytes),
	}
}

// loadSettings load settings from filepsace
func loadSettings(fs filesystem.Filespace, requireSafe bool) (settings Settings, err error) {
	var (
		json        string
		byteContent []byte
	)
	if requireSafe {
		if err = preventUnsafe(fs); err != nil {
			return settings, err
		}
	}
	if fs.IsFile(configPath) {
		if byteContent, err = fs.ReadFile(configPath); err != nil {
			return settings, err
		}
		if err = varutil.ObjectFromJSON(&settings, string(byteContent)); err != nil {
			return settings, err
		}
		return settings, nil
	}
	settings = NewSettings()
	if json, err = varutil.ObjectToJSON(&settings); err != nil {
		return settings, err
	}
	if err = fs.MkdirAll(filepath.Dir(configPath), filesystem.DefaultUnixDirMode); err != nil {
		return settings, err
	}
	if err = fs.WriteFile(configPath, []byte(json), filesystem.SafeFilePermissions); err != nil {
		return settings, err
	}
	return settings, nil
}

func preventUnsafe(fs filesystem.Filespace) (err error) {
	var info os.FileInfo
	if fs.IsFile(configPath) {
		if info, err = fs.Lstat(configPath); err != nil {
			return err
		}
		if info.Mode() != filesystem.SafeFilePermissions {
			return resetfs(fs)
		}
	} else {
		return resetfs(fs)
	}
	return nil
}

func resetfs(fs filesystem.Filespace) (err error) {
	if err = fs.Remove(configPath); err != nil {
		return err
	}
	return fs.RemoveAll(dataPath)
}
