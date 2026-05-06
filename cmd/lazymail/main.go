package main

import (
	"fmt"
	"os"

	"lazymail/internal/app"
	"lazymail/internal/config"
	"lazymail/internal/tui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "lazymail: load config: %v\n", err)
		os.Exit(1)
	}

	if len(cfg.Accounts) == 0 {
		cfg, err = config.RunOnboarding(os.Stdin, os.Stdout, cfg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "lazymail: onboarding: %v\n", err)
			os.Exit(1)
		}
	}

	svc, err := app.New(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "lazymail: init service: %v\n", err)
		os.Exit(1)
	}
	defer svc.Close()

	// Sync in background before starting TUI (non-blocking: errors are logged).
	if len(cfg.Accounts) > 0 {
		go svc.Sync()
	}

	account := ""
	if len(cfg.Accounts) > 0 {
		account = cfg.Accounts[0].Name
	}

	model := tui.NewWithService(svc, account)
	p := tea.NewProgram(model, tea.WithAltScreen(), tea.WithMouseAllMotion())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "lazymail failed to start: %v\n", err)
		os.Exit(1)
	}
}
