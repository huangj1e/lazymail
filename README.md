# LazyMail

LazyMail is a terminal-first email workspace inspired by the panel-based workflow of tools like lazygit.

The project is currently in early bootstrap stage. The Go development environment and repository structure are ready, and the next milestone is a usable TUI shell.

## Vision

- Fast email triage from the terminal
- Keyboard-first workflow with mouse support
- Multi-panel interface (sidebar, mail list, viewer)
- Extensible architecture for future automation (including AI-assisted flows)

## Current Status

- Go module initialized
- Minimal runnable entrypoint at `cmd/lazymail/main.go`
- VS Code Go tooling settings added (`gopls`, `goimports`, `staticcheck`)
- Product notes are available under the `docs/` directory

## Requirements

- Go 1.26+

## Quick Start

```bash
git clone <your-repo-url>
cd lazymail
go run ./cmd/lazymail
```

Expected output:

```text
LazyMail Go environment is ready.
```

## Development Commands

Run app:

```bash
go run ./cmd/lazymail
```

Format code:

```bash
goimports -w .
```

Static analysis:

```bash
staticcheck ./...
```

Build binary:

```bash
go build -o bin/lazymail ./cmd/lazymail
```

## Planned Architecture

```text
cmd/lazymail/         # entrypoint
internal/tui/         # Bubble Tea UI models and panels
internal/mail/        # IMAP/SMTP adapters
internal/store/       # local cache (SQLite)
internal/config/      # configuration loading
internal/domain/      # core domain models
internal/app/         # app orchestration
```

## Roadmap

1. Build a Bubble Tea multi-panel layout shell.
2. Add folder switching and mock mail list navigation.
3. Integrate IMAP read flow and SMTP send flow.
4. Add local caching and search.
5. Add quality-of-life features (shortcuts, mouse gestures, status bar actions).

## License

This project is licensed under the terms of the [LICENSE](LICENSE) file.
