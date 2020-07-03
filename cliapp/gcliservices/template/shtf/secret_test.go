package shtf

import "testing"

func TestSecretEnv(t *testing.T) {
	t.Parallel()
	expectedResult := "$SECRET_USER_TOM_PASSWORD"
	result := SecretEnv("user.tom.password")
	if result != expectedResult {
		t.Errorf("expected result is %v (and take %v)", expectedResult, result)
		return
	}
}
