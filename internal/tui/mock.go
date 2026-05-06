package tui

import (
	"fmt"
	"time"

	"lazymail/internal/domain"
)

func mockMailsForFolder(folder string) []domain.Mail {
	now := time.Now()
	return []domain.Mail{
		{
			ID:      fmt.Sprintf("%s-001", folder),
			Subject: "Welcome to LazyMail",
			From:    "hello@lazymail.dev",
			To:      []string{"you@example.com"},
			Date:    now.Add(-2 * time.Hour),
			Body:    "Thanks for trying LazyMail. This is a mock message used for Phase 1 interaction testing.",
			IsRead:  false,
			Folder:  folder,
		},
		{
			ID:      fmt.Sprintf("%s-002", folder),
			Subject: "Daily report draft",
			From:    "ops@example.com",
			To:      []string{"you@example.com"},
			Date:    now.Add(-6 * time.Hour),
			Body:    "Report summary:\n- Jobs OK\n- Latency stable\n- No incidents",
			IsRead:  true,
			Folder:  folder,
		},
		{
			ID:      fmt.Sprintf("%s-003", folder),
			Subject: "Product sync notes",
			From:    "team@example.com",
			To:      []string{"you@example.com"},
			Date:    now.Add(-24 * time.Hour),
			Body:    "Action items:\n1) Validate Bubble Tea panel focus flow\n2) Wire IMAP adapter in Phase 2\n3) Add compose modal in Phase 3",
			IsRead:  false,
			Folder:  folder,
		},
	}
}
