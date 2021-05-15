package gcliio

import (
	"strings"

	"github.com/goatcms/goatcli/cliapp/common/naming"
	"github.com/goatcms/goatcli/cliapp/gcliservices"
	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/app/modules/commonm/commservices"
	"github.com/goatcms/goatcore/filesystem"
	"github.com/goatcms/goatcore/filesystem/filespace/diskfs"
	homedir "github.com/mitchellh/go-homedir"
)

const (
	sshPubPath    = ".ssh/id_rsa.pub"
	sshSecretPath = ".ssh/id_rsa"
)

// EnvironmentManager controll gcli environments
type EnvironmentManager struct {
	deps struct {
		EnvironmentsUnit commservices.EnvironmentsUnit `dependency:"CommonEnvironmentsUnit"`
	}
}

// EnvironmentFactory create new Environment instance
func EnvironmentFactory(dp app.DependencyProvider) (result interface{}, err error) {
	inst := &EnvironmentManager{}
	if err = dp.InjectTo(&inst.deps); err != nil {
		return nil, err
	}
	return gcliservices.GCLIEnvironment(inst), nil
}

// LoadEnvs load environments to scope
func (manager *EnvironmentManager) LoadEnvs(ctxScope app.Scope, propertiesData, secretsData map[string]string) (err error) {
	var (
		envs           commservices.Environments
		homePath       string
		homeFS         filesystem.Filespace
		sshPubBytes    []byte
		sshSecretBytes []byte
	)
	if envs, err = manager.deps.EnvironmentsUnit.Envs(ctxScope); err != nil {
		return err
	}
	// Convert secrets to envs
	for key, value := range secretsData {
		key = "SECRET_" + strings.ToUpper(naming.ToUnderscore(key))
		if err = envs.Set(key, value); err != nil {
			return err
		}
	}
	// Convert properties to envs
	for key, value := range propertiesData {
		key = "PROPERTY_" + strings.ToUpper(naming.ToUnderscore(key))
		if err = envs.Set(key, value); err != nil {
			return err
		}
	}
	// load certs
	if homePath, err = homedir.Dir(); err != nil {
		return err
	}
	if homeFS, err = diskfs.NewFilespace(homePath); err != nil {
		return err
	}
	sshCert := commservices.SSHCert{
		Public: secretsData["system.sshcert.public"],
		Secret: secretsData["system.sshcert.secret"],
	}
	if sshCert.Public == "" && homeFS.IsFile(sshPubPath) {
		if sshPubBytes, err = homeFS.ReadFile(sshPubPath); err != nil {
			return err
		}
		sshCert.Public = string(sshPubBytes)
	}
	if sshCert.Secret == "" && homeFS.IsFile(sshSecretPath) {
		if sshSecretBytes, err = homeFS.ReadFile(sshSecretPath); err != nil {
			return err
		}
		sshCert.Secret = string(sshSecretBytes)
	}
	envs.SetSSHCert(sshCert)
	return nil
}
