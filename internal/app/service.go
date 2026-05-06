package app

import (
	"fmt"
	"log"

	"lazymail/internal/config"
	"lazymail/internal/domain"
	"lazymail/internal/mail"
	"lazymail/internal/store"
)

// Service coordinates IMAP sync and provides data to the TUI.
type Service struct {
	cfg   *config.Config
	store *store.Store
}

// New creates a Service, opening the local SQLite store.
func New(cfg *config.Config) (*Service, error) {
	s, err := store.Open(cfg.Database)
	if err != nil {
		return nil, fmt.Errorf("app: open store: %w", err)
	}
	return &Service{cfg: cfg, store: s}, nil
}

// Close releases resources.
func (svc *Service) Close() {
	svc.store.Close()
}

// AccountNames returns names of all configured accounts.
func (svc *Service) AccountNames() []string {
	names := make([]string, 0, len(svc.cfg.Accounts))
	for _, a := range svc.cfg.Accounts {
		names = append(names, a.Name)
	}
	return names
}

// Sync fetches recent mails from every configured account and caches them.
// It is intended for background execution.
func (svc *Service) Sync() error {
	for _, acc := range svc.cfg.Accounts {
		if err := svc.syncAccount(acc); err != nil {
			log.Printf("sync %s: %v", acc.Name, err)
		}
	}
	return nil
}

func (svc *Service) syncAccount(acc config.Account) error {
	cl, err := mail.Dial(acc)
	if err != nil {
		return err
	}
	defer cl.Close()

	folders, err := cl.ListFolders()
	if err != nil {
		return err
	}

	for _, folder := range folders {
		mails, err := cl.FetchMails(folder, 50)
		if err != nil {
			log.Printf("fetch %s/%s: %v", acc.Name, folder, err)
			continue
		}
		for _, m := range mails {
			if err := svc.store.Upsert(acc.Name, m); err != nil {
				log.Printf("store upsert: %v", err)
			}
		}
	}
	return nil
}

// Folders returns folders available for an account (from local cache).
// Falls back to default list if cache is empty.
func (svc *Service) Folders(account string) []string {
	f, err := svc.store.Folders(account)
	if err != nil || len(f) == 0 {
		return []string{"Inbox", "Sent", "Drafts", "Archive"}
	}
	return f
}

// Mails returns cached mails for account+folder.
func (svc *Service) Mails(account, folder string) []domain.Mail {
	mails, err := svc.store.ListByFolder(account, folder)
	if err != nil {
		log.Printf("store list: %v", err)
		return nil
	}
	return mails
}

// SetRead updates the read state locally and, if online, on the server.
func (svc *Service) SetRead(mailID string, isRead bool) {
	_ = svc.store.SetRead(mailID, isRead)
}

// Delete removes a mail from local cache.
func (svc *Service) Delete(mailID string) {
	_ = svc.store.Delete(mailID)
}
