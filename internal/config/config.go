package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Account holds IMAP and SMTP settings for one email account.
type Account struct {
	Name     string `yaml:"name"`
	Email    string `yaml:"email"`
	IMAPHost string `yaml:"imap_host"`
	IMAPPort int    `yaml:"imap_port"`
	SMTPHost string `yaml:"smtp_host"`
	SMTPPort int    `yaml:"smtp_port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	TLS      bool   `yaml:"tls"`
}

// Config is the root configuration structure.
type Config struct {
	Accounts []Account `yaml:"accounts"`
	Database string    `yaml:"database"` // path to SQLite file
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
