package store

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"lazymail/internal/domain"

	_ "modernc.org/sqlite"
)

// Store is a SQLite-backed local mail cache.
type Store struct {
	db *sql.DB
}

// Open opens (or creates) the SQLite database at the given path.
func Open(path string) (*Store, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("store: open: %w", err)
	}
	s := &Store{db: db}
	if err := s.migrate(); err != nil {
		db.Close()
		return nil, err
	}
	return s, nil
}

// Close releases the database connection.
func (s *Store) Close() error {
	return s.db.Close()
}

func (s *Store) migrate() error {
	_, err := s.db.Exec(`CREATE TABLE IF NOT EXISTS mails (
		id      TEXT PRIMARY KEY,
		account TEXT NOT NULL,
		folder  TEXT NOT NULL,
		subject TEXT NOT NULL,
		from_addr TEXT NOT NULL,
		to_addrs TEXT NOT NULL,
		date    INTEGER NOT NULL,
		body    TEXT NOT NULL,
		is_read INTEGER NOT NULL DEFAULT 0
	)`)
	return err
}

// Upsert inserts or replaces a mail record.
func (s *Store) Upsert(account string, m domain.Mail) error {
	_, err := s.db.Exec(`INSERT OR REPLACE INTO mails
		(id, account, folder, subject, from_addr, to_addrs, date, body, is_read)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		m.ID, account, m.Folder, m.Subject, m.From,
		strings.Join(m.To, ","), m.Date.Unix(), m.Body, boolToInt(m.IsRead),
	)
	return err
}

// ListByFolder returns all mails for an account+folder, sorted newest first.
func (s *Store) ListByFolder(account, folder string) ([]domain.Mail, error) {
	rows, err := s.db.Query(`SELECT id, subject, from_addr, to_addrs, date, body, is_read, folder
		FROM mails WHERE account=? AND folder=? ORDER BY date DESC`, account, folder)
	if err != nil {
		return nil, fmt.Errorf("store: list: %w", err)
	}
	defer rows.Close()

	var mails []domain.Mail
	for rows.Next() {
		var m domain.Mail
		var toStr string
		var ts int64
		var isRead int
		if err := rows.Scan(&m.ID, &m.Subject, &m.From, &toStr, &ts, &m.Body, &isRead, &m.Folder); err != nil {
			return nil, err
		}
		m.To = strings.Split(toStr, ",")
		m.Date = time.Unix(ts, 0)
		m.IsRead = isRead == 1
		mails = append(mails, m)
	}
	return mails, rows.Err()
}

// SetRead updates the is_read flag for a mail by ID.
func (s *Store) SetRead(id string, isRead bool) error {
	_, err := s.db.Exec(`UPDATE mails SET is_read=? WHERE id=?`, boolToInt(isRead), id)
	return err
}

// Delete removes a mail by ID.
func (s *Store) Delete(id string) error {
	_, err := s.db.Exec(`DELETE FROM mails WHERE id=?`, id)
	return err
}

// Folders returns the distinct folder names for an account.
func (s *Store) Folders(account string) ([]string, error) {
	rows, err := s.db.Query(`SELECT DISTINCT folder FROM mails WHERE account=? ORDER BY folder`, account)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var folders []string
	for rows.Next() {
		var f string
		if err := rows.Scan(&f); err != nil {
			return nil, err
		}
		folders = append(folders, f)
	}
	return folders, rows.Err()
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
