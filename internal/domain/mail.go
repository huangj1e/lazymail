package domain

import "time"

// Mail is the core domain model used by list and viewer panels.
type Mail struct {
	ID      string
	Subject string
	From    string
	To      []string
	Date    time.Time
	Body    string
	IsRead  bool
	Folder  string
}
