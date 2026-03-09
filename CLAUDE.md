# CLAUDE.md - Project Conventions

## Overview

skill-cli is a Go CLI tool for managing AI agent skills across platforms.

## Project Structure

```
skill-cli/
├── cmd/cli/main.go       # CLI entry point
├── internal/
│   ├── create/          # Skill scaffolding
│   ├── validate/        # Schema & best-practice validation
│   ├── discover/        # Find installed skills
│   ├── convert/         # Format conversion
│   └── formats/          # SKILL.md parsing
└── go.mod
```

## Commands

- `skill-cli create` — Create new skill from template
- `skill-cli validate` — Validate skill definition
- `skill-cli list` — List installed skills
- `skill-cli convert` — Convert between formats

## Development

```bash
go build ./...
go test ./...
go run ./cmd/cli --help
```

## Commit Convention

Use conventional commits:
- `feat(name): description`
- `fix(name): description`
- `docs: description`

## Testing

```bash
go test ./... -v
```
