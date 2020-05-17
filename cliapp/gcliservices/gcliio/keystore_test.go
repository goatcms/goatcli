package gcliio

import (
	"bytes"
	"testing"

	"github.com/goatcms/goatcli/cliapp/gcliservices"
	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/app/mockupapp"
)

func TestKeystoreSetGetAndOverwriteStory(t *testing.T) {
	t.Parallel()
	var (
		err       error
		mapp      app.App
		expected  = []byte{1, 2, 3}
		expected2 = []byte{4, 5, 6}
		result    []byte
		deps      struct {
			Keystorage gcliservices.GCLIKeystorageService `dependency:"GCLIKeystorage"`
		}
	)
	if mapp, err = newMockupApp(mockupapp.MockupOptions{}); err != nil {
		t.Error(err)
		return
	}
	if err = mapp.DependencyProvider().InjectTo(&deps); err != nil {
		t.Error(err)
		return
	}
	if err = deps.Keystorage.Set("key", expected); err != nil {
		t.Error(err)
		return
	}
	if result, err = deps.Keystorage.Get("key"); err != nil {
		t.Error(err)
		return
	}
	if !bytes.Equal(expected, result) {
		t.Errorf(" expected equels result and expected and take %v and %v", result, expected)
		return
	}
	if err = deps.Keystorage.Sets(map[string][]byte{
		"key": expected2,
	}); err != nil {
		t.Error(err)
		return
	}
	if result, err = deps.Keystorage.Get("key"); err != nil {
		t.Error(err)
		return
	}
	if !bytes.Equal(expected2, result) {
		t.Errorf("After overwrite expected equels result and expected and take %v and %v", result, expected2)
		return
	}
}

func TestKeystoreSetAndRemoveStory(t *testing.T) {
	t.Parallel()
	var (
		err      error
		mapp     app.App
		expected = []byte{1, 2, 3}
		result   []byte
		deps     struct {
			Keystorage gcliservices.GCLIKeystorageService `dependency:"GCLIKeystorage"`
		}
	)
	if mapp, err = newMockupApp(mockupapp.MockupOptions{}); err != nil {
		t.Error(err)
		return
	}
	if err = mapp.DependencyProvider().InjectTo(&deps); err != nil {
		t.Error(err)
		return
	}
	if err = deps.Keystorage.Set("key", expected); err != nil {
		t.Error(err)
		return
	}
	if result, err = deps.Keystorage.Get("key"); err != nil {
		t.Error(err)
		return
	}
	if !bytes.Equal(expected, result) {
		t.Errorf(" expected equels result and expected and take %v and %v", result, expected)
		return
	}
	if err = deps.Keystorage.Remove("key"); err != nil {
		t.Error(err)
		return
	}
	if result, err = deps.Keystorage.Get("key"); err != nil {
		t.Error(err)
		return
	}
	if result != nil {
		t.Errorf("After remove expected nil and take %v", result)
		return
	}
}

func TestKeystoreSetAndPasswordStory(t *testing.T) {
	t.Parallel()
	var (
		err      error
		mapp     app.App
		expected = []byte{1, 2, 3}
		result   []byte
		deps     struct {
			Keystorage gcliservices.GCLIKeystorageService `dependency:"GCLIKeystorage"`
		}
	)
	if mapp, err = newMockupApp(mockupapp.MockupOptions{
		Args: []string{"interactive=false"},
	}); err != nil {
		t.Error(err)
		return
	}
	if err = mapp.DependencyProvider().InjectTo(&deps); err != nil {
		t.Error(err)
		return
	}
	if err = deps.Keystorage.Set("key", expected); err != nil {
		t.Error(err)
		return
	}
	if result, err = deps.Keystorage.Password("key", "insert value"); err != nil {
		t.Error(err)
		return
	}
	if !bytes.Equal(expected, result) {
		t.Errorf(" expected equels result and expected and take %v and %v", result, expected)
		return
	}
}

func TestKeystoreUnknowPasswordStory(t *testing.T) {
	t.Parallel()
	var (
		err  error
		mapp app.App
		deps struct {
			Keystorage gcliservices.GCLIKeystorageService `dependency:"GCLIKeystorage"`
		}
	)
	if mapp, err = newMockupApp(mockupapp.MockupOptions{
		Args: []string{"interactive=false"},
	}); err != nil {
		t.Error(err)
		return
	}
	if err = mapp.DependencyProvider().InjectTo(&deps); err != nil {
		t.Error(err)
		return
	}
	if _, err = deps.Keystorage.Password("key", "insert value"); err == nil {
		t.Errorf("expected error for unknow key and non interactive mode")
		return
	}
}
