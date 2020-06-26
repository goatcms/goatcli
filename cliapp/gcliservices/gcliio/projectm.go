package gcliio

import (
	"fmt"
	"os"
	"strings"

	"github.com/goatcms/goatcli/cliapp/common/am"
	"github.com/goatcms/goatcli/cliapp/common/config"
	"github.com/goatcms/goatcli/cliapp/common/prevents"
	"github.com/goatcms/goatcli/cliapp/common/uefs"
	"github.com/goatcms/goatcli/cliapp/gcliservices"
	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/dependency"
	"github.com/goatcms/goatcore/filesystem"
	"github.com/goatcms/goatcore/varutil"
	"github.com/goatcms/goatcore/varutil/goaterr"
)

// ProjectManager run application scripts
type ProjectManager struct {
	deps struct {
		Interactive       string                             `argument:"?interactive" ,command:"?interactive"`
		ProfileName       string                             `argument:"?profile" ,command:"?profile"`
		PropertiesService gcliservices.PropertiesService     `dependency:"PropertiesService"`
		SecretsService    gcliservices.SecretsService        `dependency:"SecretsService"`
		DataService       gcliservices.DataService           `dependency:"DataService"`
		Keystorage        gcliservices.GCLIKeystorageService `dependency:"GCLIKeystorage"`
	}
	interactive bool
}

// ProjectManagerFactory create new ProjectManager instance
func ProjectManagerFactory(dp dependency.Provider) (result interface{}, err error) {
	pm := &ProjectManager{}
	if err = dp.InjectTo(&pm.deps); err != nil {
		return nil, err
	}
	if pm.deps.ProfileName == "" {
		pm.deps.ProfileName = "main"
	}
	pm.interactive = strings.ToLower(pm.deps.Interactive) != "false"
	return gcliservices.GCLIProjectManager(pm), nil
}

// Project return goat cli project (contains data, secrets, properties and filespaces)
func (pm *ProjectManager) Project(ctx app.IOContext) (project *gcliservices.Project, err error) {
	var (
		propertiesDef []*config.Property
		secretsDef    []*config.Property
		isChanged     bool
		data          map[string]string
	)
	project = &gcliservices.Project{
		FS: ctx.IO().CWD(),
	}
	if err = pm.initLocal(project); err != nil {
		return nil, err
	}
	if err = prevents.RequireGoatProject(project.FS); err != nil {
		return nil, err
	}
	// load properties
	if propertiesDef, err = pm.deps.PropertiesService.ReadDefFromFS(project.FS); err != nil {
		return nil, err
	}
	if project.Properties, err = pm.deps.PropertiesService.ReadDataFromFS(project.FS); err != nil {
		return nil, err
	}
	if isChanged, err = pm.deps.PropertiesService.FillData(ctx, propertiesDef, project.Properties, map[string]string{}, pm.interactive); err != nil {
		return nil, err
	}
	if isChanged {
		if err = pm.deps.PropertiesService.WriteDataToFS(project.FS, project.Properties); err != nil {
			return nil, err
		}
	}
	// load data
	if data, err = pm.deps.DataService.ReadDataFromFS(project.FS); err != nil {
		return nil, err
	}
	project.Data = am.NewApplicationData(data)
	// load secrets
	if secretsDef, err = pm.deps.SecretsService.ReadDefFromFS(project.FS, project.Properties, project.Data); err != nil {
		return nil, err
	}
	if pm.HasSecretsProfile(project, pm.deps.ProfileName) {
		if err = pm.LoadSecretsProfile(ctx, project, pm.deps.ProfileName); err != nil {
			return nil, err
		}
	} else {
		if len(secretsDef) != 0 {
			if err = pm.CreateSecretsProfile(ctx, project, pm.deps.ProfileName); err != nil {
				return nil, err
			}
		}
	}
	return project, nil
}

func (pm *ProjectManager) initLocal(project *gcliservices.Project) (err error) {
	var (
		json    string
		content []byte
	)
	localFilePath := gcliservices.ProjectTmpDirPath + localFileName
	if !project.FS.IsFile(localFilePath) {
		project.Local.ID = varutil.RandString(25, varutil.AlphaNumericBytes)
		if err = project.FS.MkdirAll(gcliservices.ProjectTmpDirPath, filesystem.DefaultUnixDirMode); err != nil {
			return err
		}
		if json, err = varutil.ObjectToJSON(&project.Local); err != nil {
			return err
		}
		if err = project.FS.WriteFile(localFilePath, []byte(json), filesystem.SafeFilePermissions); err != nil {
			return err
		}
		return nil
	}
	if content, err = project.FS.ReadFile(localFilePath); err != nil {
		return err
	}
	return varutil.ObjectFromJSON(&project.Local, string(content))
}

// ListSecretProfiles return project secrets profiles list
func (pm *ProjectManager) ListSecretProfiles(project *gcliservices.Project) (result []string, err error) {
	var nodes []os.FileInfo
	if nodes, err = project.FS.ReadDir(secretProfileBasePath); err != nil {
		return nil, err
	}
	for _, row := range nodes {
		result = append(result, row.Name())
	}
	return result, nil
}

// HasSecretsProfile return true if project contains secret profile
func (pm *ProjectManager) HasSecretsProfile(project *gcliservices.Project, profileName string) bool {
	return project.FS.IsDir(secretProfileBasePath + profileName)
}

// CreateSecretsProfile create new secrets profile
func (pm *ProjectManager) CreateSecretsProfile(ctx app.IOContext, project *gcliservices.Project, profileName string) (err error) {
	var passwordHash []byte
	propmpt := fmt.Sprintf("Insert passowrd for new '%s' profile: ", profileName)
	rePropmpt := fmt.Sprintf("Repeat passowrd for '%s' profile: ", profileName)
	if passwordHash, err = pm.deps.Keystorage.RePassword(project.Local.ID+passwordHashPK, propmpt, rePropmpt); err != nil {
		return err
	}
	if err = pm.createSecretsProfile(ctx, project, profileName, passwordHash); err != nil {
		return err
	}
	data := map[string][]byte{}
	data[project.Local.ID+passwordHashPK] = passwordHash
	data[project.Local.ID+lastProfilePK] = []byte(profileName)
	return pm.deps.Keystorage.Sets(data)
}

func (pm *ProjectManager) createSecretsProfile(ctx app.IOContext, project *gcliservices.Project, profileName string, passwordHash []byte) (err error) {
	var (
		bPath      = secretProfileBasePath + profileName
		fs         filesystem.Filespace
		secretsDef []*config.Property
	)
	if project.FS.IsDir(bPath) {
		return goaterr.Errorf("Project profile %s exists", profileName)
	}
	if err = project.FS.MkdirAll(bPath, filesystem.DefaultUnixDirMode); err != nil {
		return err
	}
	if fs, err = project.FS.Filespace(bPath); err != nil {
		pm.RemoveSecretProfile(project, profileName)
		return err
	}
	if project.EncryptFS, err = uefs.BuildEFS(fs, false, false, passwordHash); err != nil {
		pm.RemoveSecretProfile(project, profileName)
		return err
	}
	project.Secrets = map[string]string{}
	if secretsDef, err = pm.deps.SecretsService.ReadDefFromFS(project.EncryptFS, project.Properties, project.Data); err != nil {
		return err
	}
	if _, err = pm.deps.SecretsService.FillData(ctx, secretsDef, project.Secrets, map[string]string{}, pm.interactive); err != nil {
		return err
	}
	if err = pm.deps.SecretsService.WriteDataToFS(fs, project.Secrets); err != nil {
		return err
	}
	return nil
}

// LoadSecretsProfile load project secrets profile by name
func (pm *ProjectManager) LoadSecretsProfile(ctx app.IOContext, project *gcliservices.Project, profileName string) (err error) {
	var (
		passwordHash []byte
	)
	propmpt := fmt.Sprintf("Insert passowrd for '%s' profile: ", profileName)
	if passwordHash, err = pm.deps.Keystorage.Password(project.Local.ID+passwordHashPK, propmpt); err != nil {
		return err
	}
	if err = pm.loadSecretsProfile(ctx, project, profileName, passwordHash); err != nil {
		return err
	}
	return pm.deps.Keystorage.Set(project.Local.ID+lastProfilePK, []byte(profileName))
}

// LoadSecretsProfile load project secrets profile by name
func (pm *ProjectManager) loadSecretsProfile(ctx app.IOContext, project *gcliservices.Project, profileName string, passwordHash []byte) (err error) {
	var (
		bPath      = secretProfileBasePath + profileName
		fs         filesystem.Filespace
		secretsDef []*config.Property
		isChanged  bool
	)
	if !project.FS.IsDir(bPath) {
		return goaterr.Errorf("Project profile %s is not exist", profileName)
	}
	if fs, err = project.FS.Filespace(bPath); err != nil {
		pm.RemoveSecretProfile(project, profileName)
		return err
	}
	if project.EncryptFS, err = uefs.BuildEFS(fs, false, false, passwordHash); err != nil {
		pm.RemoveSecretProfile(project, profileName)
		return err
	}
	if secretsDef, err = pm.deps.SecretsService.ReadDefFromFS(project.EncryptFS, project.Properties, project.Data); err != nil {
		return err
	}
	if project.Secrets, err = pm.deps.SecretsService.ReadDataFromFS(project.EncryptFS); err != nil {
		return err
	}
	if isChanged, err = pm.deps.SecretsService.FillData(ctx, secretsDef, project.Secrets, map[string]string{}, pm.interactive); err != nil {
		return err
	}
	if isChanged {
		if err = pm.deps.SecretsService.WriteDataToFS(fs, project.Secrets); err != nil {
			return err
		}
	}
	return nil
}

// RemoveSecretProfile remove secrets profile
func (pm *ProjectManager) RemoveSecretProfile(project *gcliservices.Project, profileName string) (err error) {
	var bPath = secretProfileBasePath + profileName
	return project.FS.RemoveAll(bPath)
}
