# 🎒 packs

skills.sh compatible. Enterprise ready.

```bash
brew install packs
```

## Quick Start

```bash
# Browse packs in the TUI
packs

# Search for skills
packs find "commit message"

# Install from registry
packs get commit-message

# Install from GitHub
packs get @blader/humanizer

# Pipe to clipboard
packs get commit-message | pbcopy
```

## What Are Packs?

Packs are reusable instructions for AI agents. Three types:

| Type | File | Purpose |
|------|------|---------|
| **Skill** | `SKILL.md` | How to do X (procedural) |
| **Context** | `CONTEXT.md` | What is X (knowledge) |
| **Prompt** | `PROMPT.md` | Ready-to-use instructions |

## Commands

### `packs` — TUI Browser

Launch the interactive browser to discover and install packs.

```
  🎒 packs                                      [?] help [q] quit
  ────────────────────────────────────────────────────────────────
  [1] All  [2] Skills   [3] Contexts   [4] Prompts
  ────────────────────────────────────────────────────────────────

> 📦 commit-message          1.0.0  ★ 892   Generate conventional commits
  📦 humanizer               1.0.0  ★ 543   Remove AI patterns from writing
  📚 react-query             2.1.0  ★ 1247  React Query patterns and best...
```

**Keys:** `↑↓` navigate · `⏎` details · `g` quick install · `/` search · `1-4` filter

### `packs get <pack>` — Install

```bash
# From packs.sh registry
packs get commit-message
packs get commit-message@1.0.0      # specific version

# From GitHub (@ shorthand)
packs get @user/repo/skill
packs get @blader/humanizer         # root-level skill

# Custom install location
packs get commit-message -o ./skills/

# Output to stdout (for piping)
packs get commit-message | pbcopy
```

**Auto-detection:** Installs to the right place based on your agent:
- Claude Code → `~/.claude/skills/`
- Clawdbot → `./skills/`
- Codex → `~/.codex/skills/`
- Generic → `~/.packs/skills/`

### `packs find [query]` — Search

```bash
packs find                          # list popular
packs find "commit"                 # search by keyword
packs find --type context           # filter by type
packs find --json                   # JSON output for agents
```

### `packs info <pack>` — Details

```bash
packs info humanizer
packs info --json react-query
```

### `packs submit <ref>` — Publish

```bash
packs submit @myname/my-skills/commit-helper
```

Requires a `pack.yaml` in your GitHub repo:

```yaml
name: commit-helper
version: 1.0.0
type: skill
description: Generate helpful commit messages
author: myname
license: MIT
tags:
  - git
  - commits
```

## For AI Agents

Packs is designed to be used by AI agents, not just humans.

```bash
# JSON output for parsing
packs find --json "testing" | jq '.[0].name'

# Machine-readable help
packs get --help

# Non-interactive install
packs get commit-message --install
```

Example agent workflow:

```
User: "Help me write better commit messages"

Agent thinks: I should check if there's a skill for this
Agent runs: packs find --json "commit message"
Agent gets: [{"name": "commit-message", "description": "Generate conventional..."}]
Agent runs: packs get commit-message
Agent reads: ~/.claude/skills/commit-message/SKILL.md
Agent follows the skill instructions
```

## Configuration

```bash
packs config              # show current config
packs config path         # config file location
packs config reset        # reset to defaults
```

Config file: `~/.packs/config.yaml`

```yaml
registry: https://api.packs.sh
telemetry: true
# skills_dir: ~/.packs/skills  # override auto-detection
```

Environment variables:
- `PACKS_REGISTRY` — override registry URL
- `PACKS_SKILLS_DIR` — override skills directory
- `PACKS_NO_TELEMETRY=1` — disable telemetry

## Creating Packs

See the [pack creation guide](https://packs.sh/docs/creating-packs) or use the template:

```
my-skill/
├── pack.yaml       # metadata
├── SKILL.md        # instructions
└── README.md       # optional docs
```

## Links

- **Website:** [packs.sh](https://packs.sh)
- **Registry:** [github.com/tunajam/packs-registry](https://github.com/tunajam/packs-registry)
- **API:** [github.com/tunajam/packs-api](https://github.com/tunajam/packs-api)

## License

MIT
