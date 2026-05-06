package config

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestAccountResolvePasswordFromEnv(t *testing.T) {
	t.Setenv("LAZYMAIL_PASSWORD_TEST", "secret-token")

	account := Account{Name: "personal", PasswordEnv: "LAZYMAIL_PASSWORD_TEST"}
	password, err := account.ResolvePassword()
	if err != nil {
		t.Fatalf("ResolvePassword() error = %v", err)
	}
	if password != "secret-token" {
		t.Fatalf("ResolvePassword() = %q, want %q", password, "secret-token")
	}
}

func TestAccountResolvePasswordEnvMissing(t *testing.T) {
	account := Account{Name: "personal", PasswordEnv: "LAZYMAIL_PASSWORD_MISSING"}
	_, err := account.ResolvePassword()
	if err == nil {
		t.Fatal("ResolvePassword() error = nil, want missing env error")
	}
	if !strings.Contains(err.Error(), "LAZYMAIL_PASSWORD_MISSING") {
		t.Fatalf("ResolvePassword() error = %v, want env name in message", err)
	}
}

func TestSaveRejectsInvalidAccount(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())

	cfg := &Config{
		Accounts: []Account{{Name: "broken"}},
	}

	err := Save(cfg)
	if err == nil {
		t.Fatal("Save() error = nil, want validation error")
	}
	if !strings.Contains(err.Error(), "email is required") {
		t.Fatalf("Save() error = %v, want email validation error", err)
	}
}

func TestRunOnboardingStoresPasswordEnv(t *testing.T) {
	configHome := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", configHome)

	input := strings.NewReader(strings.Join([]string{
		"y",
		"alice@example.com",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
	}, "\n") + "\n")
	var output bytes.Buffer

	cfg, err := RunOnboarding(input, &output, &Config{})
	if err != nil {
		t.Fatalf("RunOnboarding() error = %v", err)
	}
	if len(cfg.Accounts) != 1 {
		t.Fatalf("RunOnboarding() accounts = %d, want 1", len(cfg.Accounts))
	}
	account := cfg.Accounts[0]
	if account.Password != "" {
		t.Fatalf("Password = %q, want empty", account.Password)
	}
	if account.PasswordEnv != "LAZYMAIL_PASSWORD_ALICE" {
		t.Fatalf("PasswordEnv = %q, want %q", account.PasswordEnv, "LAZYMAIL_PASSWORD_ALICE")
	}

	data, err := os.ReadFile(filepath.Join(configHome, "lazymail", "config.yaml"))
	if err != nil {
		t.Fatalf("ReadFile(config.yaml) error = %v", err)
	}
	content := string(data)
	if strings.Contains(content, "password:") {
		t.Fatalf("config.yaml unexpectedly contains plaintext password: %s", content)
	}
	if !strings.Contains(content, "password_env: LAZYMAIL_PASSWORD_ALICE") {
		t.Fatalf("config.yaml = %s, want password_env entry", content)
	}
}
