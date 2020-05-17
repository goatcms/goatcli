package gcliio

import (
	"encoding/base64"
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
	if passwordHash, err = gopass.GetPasswdPrompt(prompt, false, os.Stdin, os.Stdout); err != nil {
		return nil, err
	}
	hash := sha3.New512()
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
	kstorage.dataMU.Lock()
	defer kstorage.dataMU.Unlock()
	if err = kstorage.checkAndLoadData(); err != nil {
		return nil, err
	}
	return base64.StdEncoding.DecodeString(kstorage.data[key])
}

func (kstorage *Keystore) checkAndLoadData() (err error) {
	var info os.FileInfo
	if kstorage.data == nil {
		return kstorage.loadData()
	}
	if !kstorage.deps.GEFS.IsFile(keystoreFileName) {
		kstorage.data = map[string]string{}
		kstorage.dataMod = time.Now()
	}
	if info, err = kstorage.deps.GEFS.Lstat(keystoreFileName); err != nil {
		return err
	}
	if !info.ModTime().Equal(kstorage.dataMod) {
		return kstorage.loadData()
	}
	return nil
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
