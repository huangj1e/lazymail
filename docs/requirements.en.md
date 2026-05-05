# LazyMail Requirements Specification

This document is the English requirements spec for LazyMail.

## 1. Project Name

LazyMail (working title)

Positioning: a terminal-native email workspace inspired by lazygit-style workflows.

## 2. Product Goals

Build a TUI email client that provides:

- Mouse-enabled interactions
- Panel-based UI layout
- IMAP/SMTP support
- Multi-account readiness
- Extensibility for AI automation

## 3. Core Design Principles

### 3.1 State-driven UI

Use Bubble Tea's state-driven loop:

```text
Model -> Update -> View
```

All UI updates must come from state transitions.

### 3.2 Panelized layout

```text
┌──────────────┬────────────────────────────┬────────────────────────────┐
│ Sidebar      │ Mail List                  │ Mail Viewer                │
│--------------│----------------------------│----------------------------│
│ Inbox        │ > Subject 1                │ From: xxx                  │
│ Sent         │   Subject 2                │ To: xxx                    │
│ Drafts       │   Subject 3                │                            │
│ Archive      │                            │ Mail content               │
└──────────────┴────────────────────────────┴────────────────────────────┘
```

### 3.3 Keyboard + mouse parity

Must support:

- Mouse click selection
- Mouse wheel scrolling
- Keyboard navigation (vim-style + arrow keys)

## 4. Technical Stack

Core stack:

- Go 1.22+
- Bubble Tea
- Bubbles
- Lip Gloss

Mail protocols:

- IMAP for receiving
- SMTP for sending

Suggested libraries:

- go-imap
- go-smtp or net/smtp

Local storage:

- SQLite for message caching

## 5. Module Architecture

```text
lazymail/
├── cmd/
│   └── lazymail/
│       └── main.go
│
├── internal/
│   ├── tui/
│   │   ├── model.go
│   │   ├── layout.go
│   │   ├── sidebar.go
│   │   ├── maillist.go
│   │   ├── viewer.go
│   │
│   ├── mail/
│   │   ├── imap.go
│   │   ├── smtp.go
│   │
│   ├── store/
│   │   └── sqlite.go
│   │
│   ├── config/
│   │   └── config.go
│   │
│   ├── domain/
│   │   └── mail.go
│   │
│   └── app/
│       └── service.go
│
└── config.yaml
```

## 6. Core Data Model

```go
type Mail struct {
    ID        string
    Subject   string
    From      string
    To        []string
    Date      time.Time
    Body      string
    IsRead    bool
    Folder    string
}
```

## 7. UI Components

### 7.1 Sidebar (left)

- Show accounts
- Show folders (Inbox/Sent/etc.)
- Support click-to-switch folders

### 7.2 Mail List (center)

- Show message list
- Support up/down keys, click selection, wheel scroll
- Highlight current selection

### 7.3 Mail Viewer (right)

- Show message details
- Support scrolling
- Support HTML-to-text fallback

### 7.4 Status Bar (bottom)

```text
[R] Reply  [F] Forward  [D] Delete  [A] Archive  [Q] Quit
```

## 8. Interaction Design

Keyboard mapping:

| Key | Action |
| --- | --- |
| Up/Down | Move selection |
| Enter | Open |
| Tab | Switch panel |
| q | Quit |
| r | Reply |
| d | Delete |

Mouse mapping:

| Action | Behavior |
| --- | --- |
| Click | Select |
| Double click | Open |
| Wheel | Scroll |
| Hover | Highlight |

## 9. Mail Features

Must-have:

- IMAP login
- Fetch message list
- Fetch message detail
- Mark read
- Delete message
- Send message via SMTP

Optional (Phase 2+):

- Search
- Labels
- Multi-account switching
- Local cache

## 10. State Model (Key)

```go
type Model struct {
    ActivePanel string

    Sidebar     SidebarModel
    MailList    MailListModel
    Viewer      ViewerModel

    Mails       []Mail
    SelectedIdx int

    Width  int
    Height int
}
```

## 11. Event Flow

```text
User Action
   ↓
Bubble Tea Msg
   ↓
Update()
   ↓
Mutate Model
   ↓
View()
   ↓
UI Update
```

## 12. Delivery Phases

### Phase 1 (Required)

- 3-panel layout
- Mock data
- Keyboard + mouse interaction
- Mail switching

### Phase 2

- IMAP integration
- Real mail retrieval

### Phase 3

- SMTP sending
- Delete/mark actions

### Phase 4

- SQLite cache
- Search

### Phase 5 (Extension)

- AI mail summary
- AI auto-reply draft
- RPA integration
