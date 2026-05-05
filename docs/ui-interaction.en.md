# LazyMail UI Interaction Details

This document defines interaction behavior for LazyMail's TUI.

## 1. Scope

This document covers:

- Layout and panel focus rules
- Keyboard behavior
- Mouse behavior
- Status bar behavior
- Loading/empty/error state behavior

## 2. Layout

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

Panel responsibilities:

- Sidebar: account/folder navigation
- Mail List: list browsing and selection
- Mail Viewer: detail reading and scrolling

## 3. Focus Model

- Exactly one panel is focused at any time.
- Focused panel has clear visual highlight.
- `Tab` switches focus to the next panel.
- `Shift+Tab` switches focus to the previous panel.
- Mouse click moves focus to clicked panel.

## 4. Keyboard Interactions

### 4.1 Global

| Key | Action |
| --- | --- |
| q | Quit |
| Tab | Next panel |
| Shift+Tab | Previous panel |
| Esc | Close modal/overlay |

### 4.2 Sidebar

| Key | Action |
| --- | --- |
| j / Down | Move down |
| k / Up | Move up |
| Enter | Open folder |

### 4.3 Mail List

| Key | Action |
| --- | --- |
| j / Down | Move down |
| k / Up | Move up |
| Enter | Open selected mail |
| r | Reply (planned) |
| d | Delete selected mail |
| a | Archive selected mail |
| u | Toggle read/unread |

### 4.4 Mail Viewer

| Key | Action |
| --- | --- |
| j / Down | Scroll down |
| k / Up | Scroll up |
| PageDown | Scroll one page down |
| PageUp | Scroll one page up |
| Home | Jump top |
| End | Jump bottom |

## 5. Mouse Interactions

| Gesture | Target | Behavior |
| --- | --- | --- |
| Single click | Sidebar item | Select folder |
| Single click | Mail row | Select mail |
| Double click | Mail row | Open mail in viewer |
| Wheel | Focused panel | Scroll content |
| Hover | Interactive items | Highlight only |

Behavior constraints:

- Hover never changes data state.
- Click on selected item is idempotent.
- Scrolling should not break selection state.

## 6. Status Bar

Status bar always shows:

- Primary shortcuts
- Current account
- Current folder
- Connectivity state
- Selection indicator (for example 12/340)

Example:

```text
[Tab] Focus  [Enter] Open  [R] Reply  [D] Delete  [Q] Quit   account:work  state:online  12/340
```

## 7. Loading, Empty, Error States

### 7.1 Loading

- Use spinner or skeleton rows.
- Keep non-destructive actions available.

### 7.2 Empty

- Empty folder: show `No messages in this folder`.
- Empty search: show `No results` and clear-filter hint.

### 7.3 Error

- Show concise error message in status bar.
- Provide retry action when possible.
- Preserve current focus in transient failures.

## 8. Modal Behavior (Planned)

Compose modal:

- Open with `c` for new mail or `r` for reply.
- `Tab` switches fields.
- `Ctrl+Enter` sends.
- `Esc` triggers discard confirmation when draft is dirty.

Confirm modal:

- Required for delete and quit-with-unsaved-draft.
- `Enter` confirms.
- `Esc` cancels.

## 9. Interaction Acceptance Checklist

- Keyboard-only path can complete core flow.
- Mouse interactions can complete selection/open/scroll flow.
- Focus indication is always unambiguous.
- Status bar reflects real-time app state.
- Empty/loading/error patterns are consistent.
