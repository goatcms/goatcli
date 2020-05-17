package uefs

import (
	"github.com/goatcms/goatcore/filesystem"
	"github.com/goatcms/goatcore/filesystem/filespace/encryptfs"
	"github.com/goatcms/goatcore/filesystem/filespace/encryptfs/cipherfs/extcfs"
)

// BuildEFS create new encrypted filesystem
func BuildEFS(fs filesystem.Filespace, hostOnly, requireSafe bool, secret []byte) (efs filesystem.Filespace, err error) {
	var (
		settings Settings
	)
	if settings, err = loadSettings(fs, requireSafe); err != nil {
		return nil, err
	}
	if efs, err = fs.Filespace(dataPath); err != nil {
		return nil, err
	}
	if efs, err = encryptfs.NewEncryptFS(efs, encryptfs.Settings{
		Salt:     []byte(settings.Salt1),
		Secret:   secret,
		HostOnly: hostOnly,
		Cipher:   extcfs.NewDefaultCipher(),
	}); err != nil {
		return nil, err
	}
	return efs, nil
}
