# skill-cli

> Cross-platform CLI for AI agent skills

Create, validate, discover, and convert skills for AI agents (OpenClaw, Claude Code, etc.)

![Go](https://img.shields.io/badge/Go-1.25+-00ADD8?style=flat&logo=go)
![MIT License](https://img.shields.io/badge/License-MIT-green.svg)

## Why skill-cli?

Each AI agent (Claude Code, OpenClaw, Cursor) has its own skill format. skill-cli provides a unified tool to manage skills across platforms.

## Features

- ✨ **Create** — Scaffold new skills from templates
- ✅ **Validate** — Full schema and best-practice validation
- 🔍 **Discover** — Find installed skills across platforms
- 🔄 **Convert** — Transform skills between formats

## Installation

```bash
go install github.com/ZSeven-W/skill-cli@latest
```

## Commands

### Create a skill

```bash
skill-cli create --name "My Skill" --description "Does useful things"
```

### Validate a skill

```bash
skill-cli validate ./my-skill
skill-cli validate ./my-skill --strict      # Treat warnings as errors
skill-cli validate ./my-skill --format json  # JSON output
```

### List installed skills

```bash
skill-cli list
```

### Convert between formats

```bash
skill-cli convert --from openclaw --to claude --input ./my-skill --output ./converted
```

## Supported Platforms

- **OpenClaw**: `~/.nvm/.../openclaw/skills/`
- **Claude Code**: `~/.claude/skills/`
- Custom paths via environment variables

## Validation Features

- Frontmatter schema checks (`name`, `description`, `version`, `tags`, `metadata`)
- SKILL.md structure validation (heading, Overview, Usage)
- Best-practice checks (description quality, examples, directory references)

## License

MIT
