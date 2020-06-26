package gcliio

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/goatcms/goatcore/varutil/goaterr"

	"golang.org/x/crypto/sha3"

	"github.com/goatcms/goatcli/cliapp/gcliservices"
	"github.com/goatcms/goatcore/dependency"
	"github.com/goatcms/goatcore/filesystem"
	"github.com/goatcms/goatcore/varutil"
	"github.com/howeyc/gopass"
)

var hashPernamentSalt = []byte{
	0x5a, 0xc2, 0x27, 0x74, 0x1c, 0x11, 0xa5, 0x79, 0x89, 0x15, 0x63, 0xc4, 0x96, 0x77, 0x0a, 0xa4,
	0x32, 0x61, 0x88, 0xb7, 0x5e, 0xd0, 0xf2, 0x3e, 0xef, 0xc5, 0x0f, 0x19, 0x06, 0xac, 0x53, 0x48,
	0x7a, 0xf4, 0xff, 0x7d, 0xe9, 0x0b, 0xb9, 0xe9, 0x7c, 0x4a, 0xb9, 0xdf, 0x82, 0x79, 0x3f, 0x51,
	0x76, 0x19, 0xba, 0xb6, 0x57, 0xba, 0xb6, 0x40, 0x2f, 0xca, 0x4b, 0x86, 0xd6, 0xc7, 0x21, 0xd5,
}

// Keystore provide secrets for host
type Keystore struct {
	deps struct {
		GEFS        filesystem.Filespace `filespace:"gefs"`
		Interactive string               `argument:"?interactive" ,command:"?interactive"`
	}
	interactive bool
	data        map[string]string
	dataMod     time.Time
	dataMU      sync.Mutex
}

// KeystoreFactory create new Keystore instance
func KeystoreFactory(dp dependency.Provider) (interface{}, error) {
	var err error
	instance := &Keystore{}
	if err = dp.InjectTo(&instance.deps); err != nil {
		return nil, err
	}
	instance.interactive = strings.ToLower(instance.deps.Interactive) != "false"
	return gcliservices.GCLIKeystorageService(instance), nil
}

// Password return password hash from stdin or keys storage
func (kstorage *Keystore) Password(key string, prompt string) (passwordHash []byte, err error) {
	if passwordHash, err = kstorage.Get(key); err != nil {
		return nil, err
	}
	if len(passwordHash) != 0 {
		return passwordHash, nil
	}
	if !kstorage.interactive {
		return nil, goaterr.Errorf("Unknow password (interactive mode off)")
	}
	for len(passwordHash) != 0 {
		if passwordHash, err = gopass.GetPasswdPrompt(prompt, false, os.Stdin, os.Stdout); err != nil {
			return nil, err
		}
	}
	return kstorage.Hash(passwordHash)
}

// RePassword return password hash from stdin with repeat (or keys storage)
func (kstorage *Keystore) RePassword(key string, prompt, rePrompt string) (passwordHash []byte, err error) {
	var rePasswordHash []byte
	if passwordHash, err = kstorage.Get(key); err != nil {
		return nil, err
	}
	if len(passwordHash) != 0 {
		return passwordHash, nil
	}
	if !kstorage.interactive {
		return nil, goaterr.Errorf("Unknow password (interactive mode off)")
	}
	for len(passwordHash) != 0 && !bytes.Equal(passwordHash, rePasswordHash) {
		for len(passwordHash) != 0 {
			if passwordHash, err = gopass.GetPasswdPrompt(prompt, false, os.Stdin, os.Stdout); err != nil {
				return nil, err
			}
		}
		if rePasswordHash, err = gopass.GetPasswdPrompt(rePrompt, false, os.Stdin, os.Stdout); err != nil {
			return nil, err
		}
		if !bytes.Equal(passwordHash, rePasswordHash) {
			fmt.Printf("Passwords must be the same")
		}
	}
	return kstorage.Hash(passwordHash)
}

// Hash hash value
func (kstorage *Keystore) Hash(value []byte) (passwordHash []byte, err error) {
	hash := sha3.New512()
	if _, err = hash.Write(hashPernamentSalt); err != nil {
		return nil, err
	}
	if _, err = hash.Write(passwordHash); err != nil {
		return nil, err
	}
	return hash.Sum(nil), nil
}

// Set key value
func (kstorage *Keystore) Set(key string, value []byte) (err error) {
	kstorage.dataMU.Lock()
	defer kstorage.dataMU.Unlock()
	if err = kstorage.checkAndLoadData(); err != nil {
		return err
	}
	kstorage.data[key] = base64.StdEncoding.EncodeToString(value)
	return kstorage.saveData()
}

// Sets set all values from map
func (kstorage *Keystore) Sets(kv map[string][]byte) (err error) {
	kstorage.dataMU.Lock()
	defer kstorage.dataMU.Unlock()
	if err = kstorage.checkAndLoadData(); err != nil {
		return err
	}
	for key, value := range kv {
		kstorage.data[key] = base64.StdEncoding.EncodeToString(value)
	}
	return kstorage.saveData()
}

// Remove value for key
func (kstorage *Keystore) Remove(key string) (err error) {
	kstorage.dataMU.Lock()
	defer kstorage.dataMU.Unlock()
	if err = kstorage.checkAndLoadData(); err != nil {
		return err
	}
	if _, ok := kstorage.data[key]; !ok {
		return nil
	}
	delete(kstorage.data, key)
	return kstorage.saveData()
}

// Get value for key
func (kstorage *Keystore) Get(key string) (value []byte, err error) {
	var (
		ok     bool
		svalue string
	)
	kstorage.dataMU.Lock()
	defer kstorage.dataMU.Unlock()
	if err = kstorage.checkAndLoadData(); err != nil {
		return nil, err
	}
	if svalue, ok = kstorage.data[key]; !ok {
		return nil, nil
	}
	return base64.StdEncoding.DecodeString(svalue)
}

func (kstorage *Keystore) checkAndLoadData() (err error) {
	var info os.FileInfo
	if kstorage.data == nil {
		return kstorage.loadData()
	}
	if !kstorage.deps.GEFS.IsFile(keystoreFileName) {
		kstorage.data = map[string]string{}
		kstorage.dataMod = time.Now()
		return
	}
	if info, err = kstorage.deps.GEFS.Lstat(keystoreFileName); err != nil {
		return err
	}
	if !info.ModTime().Equal(kstorage.dataMod) {
		return kstorage.loadData()
	}
	return
}

func (kstorage *Keystore) loadData() (err error) {
	var (
		data []byte
		info os.FileInfo
	)
	if !kstorage.deps.GEFS.IsFile(keystoreFileName) {
		kstorage.dataMod = time.Now()
		kstorage.data = map[string]string{}
		return nil
	}
	if data, err = kstorage.deps.GEFS.ReadFile(keystoreFileName); err != nil {
		return err
	}
	if err = varutil.ObjectFromJSON(&kstorage.data, string(data)); err != nil {
		return err
	}
	if info, err = kstorage.deps.GEFS.Lstat(keystoreFileName); err != nil {
		return err
	}
	kstorage.dataMod = info.ModTime()
	return nil
}

func (kstorage *Keystore) saveData() (err error) {
	var (
		data string
		info os.FileInfo
	)
	if data, err = varutil.ObjectToJSON(&kstorage.data); err != nil {
		return err
	}
	if err = kstorage.deps.GEFS.WriteFile(keystoreFileName, []byte(data), filesystem.SafeFilePermissions); err != nil {
		return err
	}
	if info, err = kstorage.deps.GEFS.Lstat(keystoreFileName); err != nil {
		return err
	}
	kstorage.dataMod = info.ModTime()
	return nil
}
