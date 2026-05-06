package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// Account holds IMAP and SMTP settings for one email account.
type Account struct {
	Name        string `yaml:"name"`
	Email       string `yaml:"email"`
	IMAPHost    string `yaml:"imap_host"`
	IMAPPort    int    `yaml:"imap_port"`
	SMTPHost    string `yaml:"smtp_host"`
	SMTPPort    int    `yaml:"smtp_port"`
	Username    string `yaml:"username"`
	Password    string `yaml:"password,omitempty"`
	PasswordEnv string `yaml:"password_env,omitempty"`
	TLS         bool   `yaml:"tls"`
}

// Validate checks whether the account has the minimum required connection fields.
func (a Account) Validate() error {
	name := strings.TrimSpace(a.Name)
	if name == "" {
		return fmt.Errorf("config: account name is required")
	}
	if strings.TrimSpace(a.Email) == "" {
		return fmt.Errorf("config: account %q: email is required", name)
	}
	if strings.TrimSpace(a.IMAPHost) == "" {
		return fmt.Errorf("config: account %q: imap_host is required", name)
	}
	if a.IMAPPort <= 0 {
		return fmt.Errorf("config: account %q: imap_port must be positive", name)
	}
	if strings.TrimSpace(a.SMTPHost) == "" {
		return fmt.Errorf("config: account %q: smtp_host is required", name)
	}
	if a.SMTPPort <= 0 {
		return fmt.Errorf("config: account %q: smtp_port must be positive", name)
	}
	if strings.TrimSpace(a.Username) == "" {
		return fmt.Errorf("config: account %q: username is required", name)
	}
	if strings.TrimSpace(a.Password) == "" && strings.TrimSpace(a.PasswordEnv) == "" {
		return fmt.Errorf("config: account %q: password or password_env is required", name)
	}
	return nil
}

// ResolvePassword returns the credential used for authentication.
func (a Account) ResolvePassword() (string, error) {
	if envName := strings.TrimSpace(a.PasswordEnv); envName != "" {
		value, ok := os.LookupEnv(envName)
		if !ok || strings.TrimSpace(value) == "" {
			return "", fmt.Errorf("config: account %q: environment variable %s is not set", a.Name, envName)
		}
		return value, nil
	}

	password := strings.TrimSpace(a.Password)
	if password == "" {
		return "", fmt.Errorf("config: account %q: password is not configured", a.Name)
	}
	return password, nil
}

// Config is the root configuration structure.
type Config struct {
	Accounts []Account `yaml:"accounts"`
	Database string    `yaml:"database"` // path to SQLite file
}

// Validate checks the config for required fields and conflicting account names.
func (cfg Config) Validate() error {
	seenNames := make(map[string]struct{}, len(cfg.Accounts))
	for _, account := range cfg.Accounts {
		if err := account.Validate(); err != nil {
			return err
		}
		nameKey := strings.ToLower(strings.TrimSpace(account.Name))
		if _, exists := seenNames[nameKey]; exists {
			return fmt.Errorf("config: duplicate account name %q", account.Name)
		}
		seenNames[nameKey] = struct{}{}
	}
	return nil
}

// Default returns a Config with sensible defaults.
func Default() Config {
	return Config{
		Database: filepath.Join(configDir(), "mails.db"),
	}
}

// Path returns the absolute path of the LazyMail config file.
func Path() string {
	return filepath.Join(configDir(), "config.yaml")
}

// Load reads config from configDir/config.yaml, creating a default if missing.
func Load() (*Config, error) {
	dir := configDir()
	if err := os.MkdirAll(dir, 0700); err != nil {
		return nil, fmt.Errorf("config: mkdir: %w", err)
	}

	path := Path()
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		cfg := Default()
		if werr := writeDefault(path, cfg); werr != nil {
			return nil, werr
		}
		return &cfg, nil
	}
	if err != nil {
		return nil, fmt.Errorf("config: read: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("config: parse: %w", err)
	}
	if cfg.Database == "" {
		cfg.Database = Default().Database
	}
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func writeDefault(path string, cfg Config) error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("config: marshal: %w", err)
	}
	return os.WriteFile(path, data, 0600)
}

// Save persists the config to disk.
func Save(cfg *Config) error {
	if cfg == nil {
		return fmt.Errorf("config: save: nil config")
	}
	dir := configDir()
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("config: mkdir: %w", err)
	}
	if cfg.Database == "" {
		cfg.Database = Default().Database
	}
	if err := cfg.Validate(); err != nil {
		return err
	}
	return writeDefault(Path(), *cfg)
}

// configDir returns the OS-appropriate config directory for lazymail.
func configDir() string {
	if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
		return filepath.Join(xdg, "lazymail")
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "lazymail")
}
