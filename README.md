# LazyMail

Terminal-first email workspace with a focused, panel-based workflow.

LazyMail aims to bring the speed of modern terminal tools to everyday email handling: triage quickly, navigate confidently, and keep context in one place.

> Status: **Alpha (bootstrap complete, core TUI under active development)**

## Table of Contents

- [Why LazyMail](#why-lazymail)
- [Highlights](#highlights)
- [Project Status](#project-status)
- [Getting Started](#getting-started)
- [Development](#development)
- [Configuration](#configuration)
- [Architecture](#architecture)
- [Roadmap](#roadmap)
- [Contributing](#contributing)
- [FAQ](#faq)
- [License](#license)

## Why LazyMail

Email clients are often either too heavy for fast triage or too minimal for serious workflows. LazyMail is designed to sit in the middle:

- **Terminal-native speed** for high-volume inbox processing.
- **Keyboard-first ergonomics** with **mouse support** where it helps.
- **Panel-based layout** (sidebar, message list, viewer) to reduce context switching.
- **Extensible architecture** to support automation and intelligent workflows.

## Highlights

### Available now

- Go workspace bootstrapped and runnable.
- Clean project entrypoint: [cmd/lazymail/main.go](cmd/lazymail/main.go).
- Baseline VS Code Go tooling configuration in [.vscode/settings.json](.vscode/settings.json).
- Product and interaction notes in [docs/需求文档.md](docs/%E9%9C%80%E6%B1%82%E6%96%87%E6%A1%A3.md) and [docs/UI 交互细节.md](docs/UI%20%E4%BA%A4%E4%BA%92%E7%BB%86%E8%8A%82.md).

### Planned

- Bubble Tea multi-panel shell.
- IMAP receive + SMTP send.
- Local caching and fast search.
- Multi-account support.
- Better automation hooks.

## Project Status

This repository is in early-stage product development.

- The development environment is ready.
- The repository structure is stabilized for iterative delivery.
- Core runtime features are not production-ready yet.

If you want to contribute early, this is a great time to shape architecture and UX.

## Getting Started

### Requirements

- Go 1.26 or newer.

### Run locally

```bash
git clone <your-repository-url>
cd lazymail
go run ./cmd/lazymail
```

Expected output:

```text
LazyMail Go environment is ready.
```

## Development

### Common commands

Run:

```bash
go run ./cmd/lazymail
```

Build:

```bash
go build -o bin/lazymail ./cmd/lazymail
```

Format:

```bash
goimports -w .
```

Static analysis:

```bash
staticcheck ./...
```

### Recommended tools

- `gopls`
- `goimports`
- `staticcheck`
- `dlv`

## Configuration

Runtime config file support is planned. A likely initial format is YAML:

```yaml
accounts:
	- name: personal
		imap_host: imap.example.com
		imap_port: 993
		smtp_host: smtp.example.com
		smtp_port: 465
ui:
	mouse: true
	theme: default
```

This is a draft schema and may change before the first stable release.

## Architecture

Current and planned structure:

```text
cmd/lazymail/         # application entrypoint
internal/tui/         # Bubble Tea models and views
internal/mail/        # IMAP/SMTP adapters
internal/store/       # local cache and persistence
internal/config/      # config loading and validation
internal/domain/      # core domain models
internal/app/         # orchestration and services
```

Design principles:

- State-driven UI updates.
- Clear separation between protocol, domain, and presentation.
- Incremental delivery with testable boundaries.

## Roadmap

1. Build a navigable Bubble Tea multi-panel shell.
2. Implement folder switching and mock message list interactions.
3. Integrate IMAP fetch flow and SMTP send flow.
4. Add local storage, indexing, and search.
5. Add quality-of-life actions (reply, archive, quick actions).

## Contributing

Contributions are welcome during alpha.

1. Open an issue to discuss the idea or bug.
2. Keep pull requests focused and atomic.
3. Run formatting and static checks before submitting.

A dedicated contributing guide may be added as the project grows.

## FAQ

### Is LazyMail production-ready?

No. It is currently alpha-stage and under active design and implementation.

### Does it support multiple accounts today?

Not yet. Multi-account support is part of the planned milestones.

### Why terminal UI for email?

To optimize triage speed, reduce context switches, and support keyboard-centric workflows.

## License

Distributed under the terms in [LICENSE](LICENSE).
