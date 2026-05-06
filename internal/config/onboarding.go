package config

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// RunOnboarding guides the user through adding the first mail account.
// If the user skips setup, the original config is returned unchanged.
func RunOnboarding(in io.Reader, out io.Writer, cfg *Config) (*Config, error) {
	if cfg == nil {
		defaultCfg := Default()
		cfg = &defaultCfg
	}

	reader := bufio.NewReader(in)
	fmt.Fprintln(out, "Welcome to LazyMail.")
	fmt.Fprintln(out, "No mail account is configured yet.")

	shouldSetup, err := askBool(reader, out, "Configure an email account now", true)
	if err != nil {
		return nil, err
	}
	if !shouldSetup {
		fmt.Fprintf(out, "Skipping setup. You can edit %s later.\n", Path())
		return cfg, nil
	}

	fmt.Fprintln(out, "")
	fmt.Fprintln(out, "Account Setup")
	fmt.Fprintln(out, "Enter to accept the value shown in brackets.")

	email, err := askString(reader, out, "Email address", "")
	if err != nil {
		return nil, err
	}
	if email == "" || !strings.Contains(email, "@") {
		return nil, fmt.Errorf("config: onboarding: invalid email address")
	}

	defaults := guessServerDefaults(email)
	accountNameDefault := strings.Split(email, "@")[0]

	name, err := askString(reader, out, "Account name", accountNameDefault)
	if err != nil {
		return nil, err
	}
	username, err := askString(reader, out, "Login username", email)
	if err != nil {
		return nil, err
	}
	password, err := askString(reader, out, "Password or app password", "")
	if err != nil {
		return nil, err
	}
	imapHost, err := askString(reader, out, "IMAP host", defaults.IMAPHost)
	if err != nil {
		return nil, err
	}
	imapPort, err := askInt(reader, out, "IMAP port", defaults.IMAPPort)
	if err != nil {
		return nil, err
	}
	smtpHost, err := askString(reader, out, "SMTP host", defaults.SMTPHost)
	if err != nil {
		return nil, err
	}
	smtpPort, err := askInt(reader, out, "SMTP port", defaults.SMTPPort)
	if err != nil {
		return nil, err
	}
	tlsEnabled, err := askBool(reader, out, "Use TLS", true)
	if err != nil {
		return nil, err
	}

	account := Account{
		Name:     name,
		Email:    email,
		IMAPHost: imapHost,
		IMAPPort: imapPort,
		SMTPHost: smtpHost,
		SMTPPort: smtpPort,
		Username: username,
		Password: password,
		TLS:      tlsEnabled,
	}

	updated := *cfg
	updated.Accounts = append([]Account(nil), cfg.Accounts...)
	updated.Accounts = append(updated.Accounts, account)
	if updated.Database == "" {
		updated.Database = Default().Database
	}
	if err := Save(&updated); err != nil {
		return nil, err
	}

	fmt.Fprintf(out, "Saved account %q to %s\n", account.Name, Path())
	return &updated, nil
}

type serverDefaults struct {
	IMAPHost string
	IMAPPort int
	SMTPHost string
	SMTPPort int
}

func guessServerDefaults(email string) serverDefaults {
	domain := ""
	if parts := strings.Split(email, "@"); len(parts) == 2 {
		domain = strings.ToLower(strings.TrimSpace(parts[1]))
	}

	switch domain {
	case "gmail.com":
		return serverDefaults{IMAPHost: "imap.gmail.com", IMAPPort: 993, SMTPHost: "smtp.gmail.com", SMTPPort: 465}
	case "qq.com":
		return serverDefaults{IMAPHost: "imap.qq.com", IMAPPort: 993, SMTPHost: "smtp.qq.com", SMTPPort: 465}
	case "163.com":
		return serverDefaults{IMAPHost: "imap.163.com", IMAPPort: 993, SMTPHost: "smtp.163.com", SMTPPort: 465}
	case "outlook.com", "hotmail.com", "live.com":
		return serverDefaults{IMAPHost: "outlook.office365.com", IMAPPort: 993, SMTPHost: "smtp.office365.com", SMTPPort: 587}
	case "icloud.com", "me.com", "mac.com":
		return serverDefaults{IMAPHost: "imap.mail.me.com", IMAPPort: 993, SMTPHost: "smtp.mail.me.com", SMTPPort: 587}
	default:
		if domain == "" {
			return serverDefaults{IMAPHost: "imap.example.com", IMAPPort: 993, SMTPHost: "smtp.example.com", SMTPPort: 465}
		}
		return serverDefaults{
			IMAPHost: "imap." + domain,
			IMAPPort: 993,
			SMTPHost: "smtp." + domain,
			SMTPPort: 465,
		}
	}
}

func askString(reader *bufio.Reader, out io.Writer, label, defaultValue string) (string, error) {
	prompt := label
	if defaultValue != "" {
		prompt += fmt.Sprintf(" [%s]", defaultValue)
	}
	prompt += ": "
	if _, err := fmt.Fprint(out, prompt); err != nil {
		return "", err
	}
	line, err := reader.ReadString('\n')
	if err != nil && err != io.EOF {
		return "", err
	}
	value := strings.TrimSpace(line)
	if value == "" {
		return defaultValue, nil
	}
	return value, nil
}

func askInt(reader *bufio.Reader, out io.Writer, label string, defaultValue int) (int, error) {
	value, err := askString(reader, out, label, fmt.Sprintf("%d", defaultValue))
	if err != nil {
		return 0, err
	}
	var parsed int
	if _, err := fmt.Sscanf(value, "%d", &parsed); err != nil || parsed <= 0 {
		return 0, fmt.Errorf("config: onboarding: invalid integer for %s", label)
	}
	return parsed, nil
}

func askBool(reader *bufio.Reader, out io.Writer, label string, defaultValue bool) (bool, error) {
	defaultText := "y"
	if !defaultValue {
		defaultText = "n"
	}
	value, err := askString(reader, out, label+" [y/n]", defaultText)
	if err != nil {
		return false, err
	}
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "y", "yes":
		return true, nil
	case "n", "no":
		return false, nil
	default:
		return false, fmt.Errorf("config: onboarding: invalid yes/no for %s", label)
	}
}