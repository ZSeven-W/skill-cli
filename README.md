# skill-cli

CLI for creating, converting, discovering, and validating `SKILL.md` based skills.

## Validate skills

Run default validation (text output):

```bash
skill-cli validate ./my-skill
```

Run with JSON output:

```bash
skill-cli validate ./my-skill --format json
```

Run strict validation (best-practice warnings fail validation):

```bash
skill-cli validate ./my-skill --strict
```

Combine strict mode with JSON output:

```bash
skill-cli validate ./my-skill --strict --format json
```

Validation now includes:
- Frontmatter schema checks (`name`, `description`, optional `version`, `metadata`, `tags`)
- `SKILL.md` structure checks (heading, overview, usage)
- Best-practice checks (description quality, examples, referenced directories)
