# Product Requirements Document: packs

---

## 1. Executive Summary

### 1.1 Product Name
**packs** v1.0

### 1.2 One-Line Description
A TUI-first tool for discovering, installing, and sharing AI skills, prompts, and context — for humans, agents, and LLMs alike.

### 1.3 The Pitch
> "What's Homebrew for AI skills?"  
> "packs"

### 1.4 Problem Statement
AI skills and prompts are scattered across GitHub repos, random directories, and platform-specific stores. Developers and AI agents can't easily discover, install, or share them. It's 2026 and we're copy-pasting YAML configs like it's 2010.

### 1.5 Proposed Solution
A unified CLI + TUI + registry for AI agent skills. **Terminal-first** — install with `brew install packs`, launch interactively with `packs`, or go direct with `packs get commit-message`. 

The central registry at **packs.sh** is the npm of AI skills. But organizations with constraints can run their own registries backed by Postgres, S3, or Git — the CLI supports pluggable adapters.

### 1.6 Target Users
| User | How They Use It |
|------|-----------------|
| **Developers** | TUI to browse, `packs get` to fetch skills |
| **AI Agents** | `packs find --json` to discover, `packs get` to load |
| **LLMs** | Given CLI as a tool to find relevant context |
| **Enterprises** | Private registry with S3/Postgres adapter |

### 1.7 Success Metrics
| Metric | Target | Timeframe |
|--------|--------|-----------|
| Homebrew installs | 500 | Month 1 |
| Packs in registry | 100 | Launch |
| Daily active users | 50 | Month 1 |
| Community publishers | 20 | Month 3 |
| Enterprise adapters | 2 (S3, Postgres) | Launch |

---

## 2. Product Vision

### 2.1 Vision Statement
Become the default way to discover and install AI skills. When someone asks "how do I add X capability to my agent?", the answer is `packs get X`.

### 2.2 Design Principles

| Principle | What It Means |
|-----------|---------------|
| **Terminal-native** | TUI is primary, web is secondary |
| **Zero friction** | `brew install packs && packs` — that's it |
| **Works offline** | Local cache, sync when online |
| **Pipe-friendly** | `packs get X` outputs raw content to stdout |
| **Agent-friendly** | `--json` for structured output |
| **Portable** | Adapters let orgs bring their own storage |

### 2.3 What Makes This Click

| Feature | Why It Matters |
|---------|----------------|
| One command install | `packs get X` — no config editing |
| TUI browser | Visual discovery without leaving terminal |
| Pluggable registries | Enterprises can self-host |
| Versioned packs | Pin to stable, or ride latest |
| Stars/favorites | Community curation surfaces quality |

---

## 3. Pack Types

Three types of packs, each serving a distinct purpose:

### 3.1 Context Pack
**Purpose**: Dense knowledge about a subject that helps AI understand what it's working with.

| Attribute | Value |
|-----------|-------|
| File | `CONTEXT.md` |
| Size | 2,000 - 20,000 tokens |
| Use Case | "Give AI deep knowledge about a library or domain" |

**Examples**: 
- React Query patterns and best practices
- Drizzle ORM conventions
- Company API documentation
- Domain-specific terminology

**Structure**:
```markdown
# [Subject] Context

## Overview
What this is and when to reference it.

## Key Concepts
- Concept 1: explanation
- Concept 2: explanation

## API Reference
Core APIs, methods, signatures.

## Patterns & Best Practices
Recommended approaches.

## Common Pitfalls
Mistakes to avoid.

## Examples
Code samples demonstrating usage.
```

---

### 3.2 Skill Pack
**Purpose**: A specific capability that teaches AI how to perform a task. Procedural and actionable.

| Attribute | Value |
|-----------|-------|
| File | `SKILL.md` |
| Size | 500 - 5,000 tokens |
| Use Case | "Teach AI a repeatable technique" |

**Examples**:
- Generate conventional commit messages
- Create PR descriptions from branch diffs
- Write unit tests for a function
- Debug a failing CI build

**Structure**:
```markdown
# Skill: [Name]

## When to Use
Situations where this skill applies.

## Prerequisites
What the AI needs access to.

## Instructions
1. Step one with details
2. Step two with details
3. ...

## Output Format
What the result should look like.

## Examples
### Example 1: [Scenario]
Input: ...
Output: ...

## Common Mistakes
- Mistake 1 and how to avoid
```

---

### 3.3 Prompt Pack
**Purpose**: Ready-to-use prompts for specific tasks. Copy-paste instructions.

| Attribute | Value |
|-----------|-------|
| File | `PROMPT.md` |
| Size | 100 - 2,000 tokens |
| Use Case | "Give AI precise instructions for a task" |

**Examples**:
- Refactor function to async/await
- Review code for security issues
- Explain code to a junior developer
- Generate API documentation

**Structure**:
```markdown
You are an expert [role].

## Task
[Clear description of what to do]

## Context
[Any relevant background]

## Constraints
- Constraint 1
- Constraint 2

## Output Format
[How the response should be structured]
```

---

### 3.4 Comparison

| Aspect | Context | Skill | Prompt |
|--------|---------|-------|--------|
| **Purpose** | Provide knowledge | Teach technique | Give instructions |
| **File** | CONTEXT.md | SKILL.md | PROMPT.md |
| **Size** | Large (2K-20K) | Medium (500-5K) | Small (100-2K) |
| **Reusability** | Reference material | Repeatable process | Task-specific |
| **Example** | "React Query API" | "How to write tests" | "Refactor this function" |

---

## 4. Pack Format

### 4.1 Directory Structure
```
my-pack/
├── pack.yaml       # Metadata
└── SKILL.md        # Content (or CONTEXT.md / PROMPT.md)
```

### 4.2 Metadata (`pack.yaml`)
```yaml
name: commit-message
version: 1.2.0
type: skill                    # skill | context | prompt
description: Generate conventional commit messages from staged changes
author: hunter
tags:
  - git
  - commits
  - conventional-commits
license: MIT
repository: https://github.com/tunajam/packs-official
```

### 4.3 Naming Rules
- **Unique names** — flat namespace, first come first served
- **Lowercase** — `commit-message` not `Commit-Message`
- **Hyphenated** — `my-pack` not `my_pack` or `mypack`
- **2-50 characters** — reasonable length

### 4.4 Versioning
- **SemVer** — `major.minor.patch`
- **Default to latest** — `packs get commit-message` gets newest
- **Pin version** — `packs get commit-message@1.2.0`
- **Version ranges** — `packs get commit-message@^1.0.0` (future)

---

## 5. Core User Flows

### 5.1 Flow 1: Interactive Discovery

```
$ packs

┌─────────────────────────────────────────────────────────────────┐
│  🎒 packs                                      [?] help [q] quit │
├─────────────────────────────────────────────────────────────────┤
│  [All]  Skills   Contexts   Prompts   ★ Starred                 │
├─────────────────────────────────────────────────────────────────┤
│  Search: _                                                       │
├─────────────────────────────────────────────────────────────────┤
│  > 📦 commit-message       1.2.0  ★ 892   Conventional commits  │
│    📦 pr-creator           1.0.0  ★ 654   Create PRs with gh    │
│    📦 react-query          2.1.0  ★ 1.2k  React Query patterns  │
│    📦 code-review          1.1.0  ★ 543   Review code for issues│
│                                                                  │
│  ↑↓ navigate  ⏎ view  c copy  g get  ★ star  / search  q quit   │
└─────────────────────────────────────────────────────────────────┘
```

User types `commit`:
```
├─────────────────────────────────────────────────────────────────┤
│  Search: commit_                                                 │
├─────────────────────────────────────────────────────────────────┤
│  > 📦 commit-message       1.2.0  ★ 892   Conventional commits  │
│    📦 commit-lint          1.0.0  ★ 234   Lint commit messages  │
│    📦 changelog-writer     0.9.0  ★ 187   Generate changelogs   │
└─────────────────────────────────────────────────────────────────┘
```

Press `⏎` to view detail:
```
┌─────────────────────────────────────────────────────────────────┐
│  ← esc                              [c] copy  [g] get  [★] star │
├─────────────────────────────────────────────────────────────────┤
│  📦 commit-message v1.2.0                                       │
│  by hunter · ★ 892 · skill · MIT                                │
│  ─────────────────────────────────────────────────────────────  │
│                                                                  │
│  # Skill: Commit Message Generator                               │
│                                                                  │
│  Generate a conventional commit message for staged changes.      │
│                                                                  │
│  ## When to Use                                                  │
│  - Before any commit                                             │
│  - When you need a well-formatted conventional commit            │
│                                                                  │
│  ## Instructions                                                 │
│  1. Check staged changes: `git diff --cached`                    │
│  ...                                                             │
│                                                                  │
│  ↑↓ scroll  c copy  g get  ★ star  esc back                     │
└─────────────────────────────────────────────────────────────────┘
```

---

### 5.2 Flow 2: Direct CLI

```bash
# Search packs.sh registry
$ packs find commit
  📦 commit-message       1.2.0  ★ 892   Conventional commits
  📦 commit-lint          1.0.0  ★ 234   Lint commit messages

# Get from packs.sh (raw content to stdout)
$ packs get commit-message
# Skill: Commit Message Generator
...

# Get directly from any GitHub repo (no registry needed!)
$ packs get gh:hsbacot/packs/commit-message
# Skill: Commit Message Generator
...

# Get specific version
$ packs get commit-message@1.0.0

# Pipe to clipboard
$ packs get commit-message | pbcopy

# Info (metadata only)
$ packs info commit-message
Name:        commit-message
Version:     1.2.0
Type:        skill
Author:      hunter
Stars:       892
Description: Generate conventional commit messages from staged changes
Tags:        git, commits, conventional-commits
```

---

### 5.3 Flow 3: Agent/LLM Usage

```bash
# Agent searches for relevant skill
$ packs find "write tests" --json
[
  {"name": "test-writer", "version": "1.0.0", "type": "skill", "stars": 421},
  {"name": "jest-patterns", "version": "2.0.0", "type": "context", "stars": 312}
]

# Agent fetches the skill
$ packs get test-writer
[raw skill content]

# Agent uses skill to write tests
```

---

### 5.4 Flow 4: Publishing

```bash
$ mkdir my-awesome-skill && cd my-awesome-skill

$ packs init
  ? Pack name: my-awesome-skill
  ? Type: skill
  ? Description: Does something awesome
  ? Tags (comma-separated): awesome, cool
  ✓ Created pack.yaml
  ✓ Created SKILL.md template

$ vim SKILL.md  # Write your skill

$ packs publish
  ✓ Validated pack structure
  ✓ Authenticated as hunter
  ✓ Published my-awesome-skill@1.0.0
  ✓ Live at packs.sh/my-awesome-skill
```

---

### 5.5 Flow 5: Private Registry (Enterprise)

```bash
# Configure private registry
$ packs config set registry https://packs.acme.corp

# Or add as secondary (checked first)
$ packs config add-registry https://packs.acme.corp --priority 1

# Now searches check private registry first, then packs.sh
$ packs find "internal api"
  📦 acme-api-context     1.0.0  ★ 12   [acme.corp] Internal API docs
  📦 api-patterns         2.1.0  ★ 654  [packs.sh] General API patterns
```

---

## 6. CLI Specification

### 6.1 Installation
```bash
brew install packs
```

### 6.2 Commands

**MVP (no auth required):**
| Command | Description |
|---------|-------------|
| `packs` | Launch interactive TUI |
| `packs find <query>` | Search packs.sh registry |
| `packs get <name>[@version]` | Get pack from packs.sh |
| `packs get gh:<user>/<repo>/<pack>` | Get pack directly from GitHub |
| `packs info <name>` | Show pack metadata |
| `packs submit gh:<user>/<repo>/<pack>` | Submit GitHub pack for indexing |
| `packs config` | Manage configuration |
| `packs version` | Show CLI version |

**Post-MVP (auth required):**
| Command | Description |
|---------|-------------|
| `packs login` | Authenticate with GitHub |
| `packs logout` | Clear credentials |
| `packs whoami` | Show current user |
| `packs star <name>` | Star a pack |
| `packs starred` | List your starred packs |
| `packs init` | Initialize new pack in current directory |
| `packs validate` | Validate pack structure |
| `packs publish` | Publish pack to packs.sh |
| `packs unpublish <name>@<version>` | Remove a version (24h window) |

### 6.3 Global Flags

| Flag | Description |
|------|-------------|
| `--json`, `-j` | Output as JSON |
| `--registry`, `-r` | Override registry URL |
| `--no-cache` | Bypass local cache |
| `--offline` | Use cached data only |
| `--verbose`, `-v` | Verbose output |
| `--help`, `-h` | Show help |

### 6.4 Config Commands

```bash
# View all config
$ packs config list

# Set default registry
$ packs config set registry https://packs.acme.corp

# Add registry with priority (lower = checked first)
$ packs config add-registry https://packs.acme.corp --priority 1

# Remove registry
$ packs config remove-registry https://packs.acme.corp

# Set cache TTL
$ packs config set cache.ttl 1h
```

### 6.5 Output Modes

**Human mode** (default for TTY):
```bash
$ packs find commit
  📦 commit-message       1.2.0  ★ 892   Conventional commits
  📦 commit-lint          1.0.0  ★ 234   Lint commit messages
```

**JSON mode** (`--json` or non-TTY):
```bash
$ packs find commit --json
[
  {"name": "commit-message", "version": "1.2.0", "stars": 892, ...},
  {"name": "commit-lint", "version": "1.0.0", "stars": 234, ...}
]
```

**Raw mode** (`packs get` always outputs raw content):
```bash
$ packs get commit-message
# Skill: Commit Message Generator
...
```

---

## 7. TUI Specification

### 7.1 Technology
| Component | Choice |
|-----------|--------|
| Language | Go 1.21+ |
| TUI Framework | Bubble Tea |
| Styling | Lip Gloss |
| Search | Built-in fuzzy matching |

### 7.2 Splash Banner

On first run or `packs --version`:

```
                    __        
    ____  ____ _____/ /_______
   / __ \/ __ `/ __/ //_/ ___/
  / /_/ / /_/ / /_/ ,< (__  ) 
 / .___/\__,_/\__/_/|_/____/  
/_/                           
        
  Skills for AI agents. One command.
  
  v0.1.0 · packs.sh
```

### 7.3 Views

**List View** (main):
```
┌─────────────────────────────────────────────────────────────────┐
│  🎒 packs                                      [?] help [q] quit │
├─────────────────────────────────────────────────────────────────┤
│  [All]  Skills   Contexts   Prompts                             │
├─────────────────────────────────────────────────────────────────┤
│  Search: _                                                       │
├─────────────────────────────────────────────────────────────────┤
│  > 📦 commit-message       1.2.0  ★ 892   Conventional commits  │
│    📦 pr-creator           1.0.0  ★ 654   Create PRs with gh    │
│    📦 react-query          2.1.0  ★ 1.2k  React Query patterns  │
│                                                                  │
│  ↑↓ navigate  ⏎ view  c copy  g get  / search  q quit           │
└─────────────────────────────────────────────────────────────────┘
```

**Detail View**:
- Pack metadata header
- Scrollable content (full markdown)
- Action bar: copy, get, star

**Help View** (`?`):
- Full keyboard shortcut reference

### 7.3 Keyboard Shortcuts

| Key | List View | Detail View |
|-----|-----------|-------------|
| `↑/↓` or `j/k` | Navigate | Scroll |
| `⏎` | View detail | — |
| `/` | Focus search | — |
| `tab` | Switch type tab | — |
| `c` | Copy selected | Copy content |
| `g` | Get & exit | Get & exit |
| `★` or `s` | Toggle star | Toggle star |
| `i` | Show info | — |
| `?` | Help | Help |
| `esc` | Clear/quit | Back to list |
| `q` | Quit | Quit |

### 7.4 Search Behavior
- Fuzzy matching on name, description, tags
- Results update as you type (debounced 100ms)
- Filtered by current type tab
- Sorted by relevance, then stars

---

## 8. Technical Architecture

### 8.1 Stack

| Layer | Technology | Why |
|-------|------------|-----|
| Language | Go 1.21+ | Single binary, fast, portable |
| TUI | Bubble Tea | Best Go TUI framework |
| CLI | Cobra | Standard Go CLI framework |
| Config | Viper | Flexible config management |
| HTTP | net/http | Stdlib, no deps |
| Clipboard | golang-design/clipboard | Cross-platform |

### 8.2 Architecture Overview

```
┌─────────────────────────────────────────────────────────────────┐
│                          packs CLI                               │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  ┌────────────┐  ┌────────────┐  ┌────────────┐                 │
│  │    TUI     │  │    CLI     │  │   JSON     │                 │
│  │ (Bubble Tea)│ │  (Cobra)   │  │  Output    │                 │
│  └─────┬──────┘  └─────┬──────┘  └─────┬──────┘                 │
│        │               │               │                         │
│        └───────────────┼───────────────┘                         │
│                        │                                         │
│                 ┌──────▼──────┐                                  │
│                 │ Pack Service │                                  │
│                 └──────┬──────┘                                  │
│                        │                                         │
│           ┌────────────┼────────────┐                            │
│           │            │            │                            │
│    ┌──────▼─────┐ ┌────▼────┐ ┌────▼─────┐                      │
│    │   Cache    │ │ Registry │ │  Auth    │                      │
│    │  Manager   │ │ Resolver │ │ Manager  │                      │
│    └──────┬─────┘ └────┬────┘ └──────────┘                      │
│           │            │                                         │
└───────────┼────────────┼─────────────────────────────────────────┘
            │            │
            ▼            ▼
┌──────────────────┐  ┌────────────────────────────────────────────┐
│   Local Cache    │  │           Registry Adapters                │
│  ~/.packs/cache  │  │                                            │
│                  │  │  ┌──────────┐ ┌──────────┐ ┌──────────┐   │
│  - index.json    │  │  │ packs.sh │ │ Postgres │ │    S3    │   │
│  - packs/        │  │  │   API    │ │ Adapter  │ │ Adapter  │   │
└──────────────────┘  │  └──────────┘ └──────────┘ └──────────┘   │
                      │                                            │
                      │  ┌──────────┐                              │
                      │  │   Git    │                              │
                      │  │ Adapter  │                              │
                      │  └──────────┘                              │
                      └────────────────────────────────────────────┘
```

### 8.3 Adapter Interface

```go
// PackStore defines the interface for pack storage backends
type PackStore interface {
    // Search packs by query with optional filters
    Search(ctx context.Context, query string, opts SearchOpts) ([]PackMeta, error)
    
    // Get pack metadata by name (latest version)
    GetMeta(ctx context.Context, name string) (*PackMeta, error)
    
    // Get pack metadata by name and version
    GetMetaVersion(ctx context.Context, name, version string) (*PackMeta, error)
    
    // Get pack content (the actual skill/context/prompt)
    GetContent(ctx context.Context, name, version string) ([]byte, error)
    
    // List all versions of a pack
    ListVersions(ctx context.Context, name string) ([]string, error)
    
    // Publish a new pack version
    Publish(ctx context.Context, pack *Pack) error
    
    // Unpublish a pack version (within 24h window)
    Unpublish(ctx context.Context, name, version string) error
    
    // Star/unstar a pack
    Star(ctx context.Context, name string, star bool) error
    
    // Get user's starred packs
    Starred(ctx context.Context) ([]PackMeta, error)
}

type SearchOpts struct {
    Type     string   // skill, context, prompt, or empty for all
    Tags     []string // filter by tags
    Author   string   // filter by author
    Limit    int      // max results
    Offset   int      // pagination
    Sort     string   // relevance, stars, newest, name
}
```

### 8.4 Adapters

**HTTP Adapter** (packs.sh and compatible APIs):
```go
type HTTPAdapter struct {
    baseURL    string
    httpClient *http.Client
    authToken  string
}
```

**Postgres Adapter** (self-hosted):
```go
type PostgresAdapter struct {
    db *sql.DB
}
```

**S3 Adapter** (file-based):
```go
type S3Adapter struct {
    client *s3.Client
    bucket string
    prefix string
}
```

**Git Adapter** (repo-based):
```go
type GitAdapter struct {
    repoURL string
    branch  string
    client  *github.Client  // or generic git
}
```

### 8.7 GitHub Direct Fetch (`gh:` prefix)

For `packs get gh:user/repo/path/to/pack`, parse and fetch:

```go
// Parse: gh:user/repo/path/to/pack → user, repo, path
// Examples:
//   gh:hsbacot/packs/commit-message → hsbacot, packs, commit-message
//   gh:hsbacot/packs/skills/commit-message → hsbacot, packs, skills/commit-message
//   gh:acme/agent-packs/support/triage → acme, agent-packs, support/triage

func getFromGitHub(user, repo, path string) ([]byte, error) {
    // First fetch pack.yaml to get the content file type
    packYaml := fetchFile(user, repo, path + "/pack.yaml")
    contentFile := determineContentFile(packYaml) // SKILL.md, CONTEXT.md, or PROMPT.md
    
    // If gh CLI is installed, use it (handles auth, private repos, SSO)
    if ghInstalled() {
        return exec("gh", "api", 
            fmt.Sprintf("/repos/%s/%s/contents/%s/%s", user, repo, path, contentFile))
    }
    
    // Fallback: raw.githubusercontent.com (public repos only)
    return httpGet(fmt.Sprintf(
        "https://raw.githubusercontent.com/%s/%s/main/%s/%s", 
        user, repo, path, contentFile))
}
```

**Benefits of `gh` CLI wrapper:**
- Already authenticated (no separate GitHub login for packs)
- Handles rate limiting gracefully
- Private repo access if user has permissions
- SSO/org access works automatically
- We don't reinvent GitHub API auth

### 8.5 Pack Resolution

Two ways to get a pack:

**Registry lookup** (default):
```bash
packs get commit-message
  → Searches packs.sh registry
  → Returns indexed pack
```

**Direct GitHub fetch** (`gh:` prefix):
```bash
# Simple (pack at repo root or /packs/)
packs get gh:hsbacot/packs/commit-message

# Subdirectories (any depth)
packs get gh:hsbacot/packs/skills/commit-message
packs get gh:hsbacot/packs/contexts/react-query
packs get gh:acme/agent-packs/support/ticket-triage

# Translates to GitHub path:
# gh:user/repo/path/to/pack
# → github.com/user/repo/path/to/pack/pack.yaml
# → github.com/user/repo/path/to/pack/SKILL.md (or CONTEXT.md, PROMPT.md)
```

Works with any public GitHub repo, any directory structure.

**Resolution order:**
1. If `gh:` prefix → direct GitHub fetch
2. If configured private registries → check those first
3. Otherwise → packs.sh (default)

**The pack lifecycle:**
```
1. Build pack locally
2. Push to GitHub: gh:yourname/packs/my-skill
3. Share: "packs get gh:yourname/packs/my-skill"
4. Gets popular? Publish to packs.sh
5. Now "packs get my-skill" works for everyone
```

This is like Go modules — `go get github.com/user/repo` works directly, but popular packages get indexed at pkg.go.dev.

### 8.6 Data Models

**PackMeta** (metadata only, for listings):
```go
type PackMeta struct {
    Name        string    `json:"name"`
    Version     string    `json:"version"`
    Type        string    `json:"type"`        // skill, context, prompt
    Description string    `json:"description"`
    Author      string    `json:"author"`
    Stars       int       `json:"stars"`
    Tags        []string  `json:"tags"`
    License     string    `json:"license"`
    Repository  string    `json:"repository"`
    CreatedAt   time.Time `json:"createdAt"`
    UpdatedAt   time.Time `json:"updatedAt"`
    Registry    string    `json:"registry"`    // which registry it came from
}
```

**Pack** (full pack with content):
```go
type Pack struct {
    PackMeta
    Content     string `json:"content"`      // the actual markdown content
    ContentHash string `json:"contentHash"`  // SHA256 for integrity
}
```

**Config** (`~/.packs/config.yaml`):
```yaml
# User auth
auth:
  token: ghp_xxxx  # from `packs login`
  user: hunter

# Registries (checked in priority order, lower = first)
registries:
  - url: https://packs.acme.corp
    priority: 1
    auth: bearer xxx  # optional per-registry auth
  # packs.sh is implicit at priority 99

# Cache settings
cache:
  dir: ~/.packs/cache
  ttl: 1h
  max_size: 100MB

# UI preferences
ui:
  color: auto  # auto, always, never
  pager: less
```

---

## 9. packs.sh Registry API

### 9.1 Overview
The central registry at packs.sh provides:
- Pack discovery and search
- Version management
- User authentication (GitHub OAuth)
- Stars/favorites
- Download statistics

### 9.2 Tech Stack
| Layer | Technology |
|-------|------------|
| Backend | Go or Node.js |
| Database | Postgres |
| Storage | S3 (pack content) |
| Auth | GitHub OAuth |
| Hosting | Vercel/Fly.io |

### 9.3 API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/api/packs` | Search/list packs |
| `GET` | `/api/packs/:name` | Get latest version |
| `GET` | `/api/packs/:name/:version` | Get specific version |
| `GET` | `/api/packs/:name/versions` | List all versions |
| `GET` | `/api/packs/:name/:version/content` | Get raw content |
| `POST` | `/api/packs` | Publish new pack |
| `DELETE` | `/api/packs/:name/:version` | Unpublish (24h window) |
| `POST` | `/api/packs/:name/star` | Star pack |
| `DELETE` | `/api/packs/:name/star` | Unstar pack |
| `GET` | `/api/user/starred` | Get starred packs |
| `GET` | `/api/auth/github` | GitHub OAuth flow |

### 9.4 Search Query Params
```
GET /api/packs?q=commit&type=skill&tags=git&sort=stars&limit=20
```

| Param | Description |
|-------|-------------|
| `q` | Search query (name, description, tags) |
| `type` | Filter: skill, context, prompt |
| `tags` | Filter by tags (comma-separated) |
| `author` | Filter by author |
| `sort` | relevance, stars, newest, name |
| `limit` | Max results (default 50) |
| `offset` | Pagination offset |

---

## 10. Web Interface (Secondary)

### 10.1 Purpose
Browse and discover packs in a browser. Secondary to CLI but useful for:
- Browsing without installing CLI
- Sharing links to packs
- Reading longer content comfortably

### 10.2 Pages

| Page | URL | Description |
|------|-----|-------------|
| Home | packs.sh | Featured packs, search |
| Browse | packs.sh/browse | Full searchable list |
| Pack Detail | packs.sh/:name | Pack content, versions |
| User Profile | packs.sh/u/:username | User's published packs |
| Publish Guide | packs.sh/publish | How to publish |
| Docs | packs.sh/docs | CLI documentation |

### 10.3 Design Notes
- Clean, minimal (think pkg.go.dev)
- Fast (SSR or static)
- Copy button prominent
- CLI commands shown: `packs get pack-name`

---

## 11. MVP Scope (4 weeks)

### Phase 1: Ship Fast (Week 1-2)
**No auth required — read-only + GitHub direct**

- [ ] Project setup (Go, Cobra, Viper, Bubble Tea)
- [ ] `packs get gh:user/repo/pack` — direct GitHub fetch via `gh` CLI or raw
- [ ] `packs get <name>` — fetch from packs.sh (read-only API)
- [ ] `packs find <query>` — search packs.sh (no auth needed)
- [ ] `packs info <name>` — metadata display
- [ ] `--json` flag for all commands
- [ ] Local cache layer

### Phase 2: TUI (Week 3-4)
- [ ] Bubble Tea list view with fuzzy search
- [ ] Type tabs (All, Skills, Contexts, Prompts)
- [ ] Detail view with scrolling
- [ ] Keyboard navigation (j/k, /, enter, q)
- [ ] Copy to clipboard
- [ ] `packs` launches TUI by default

### Phase 3: Launch
- [ ] Homebrew formula
- [ ] Seed packs.sh with 50+ quality packs (manual/script)
- [ ] Documentation (README, basic docs)
- [ ] Ship it 🚀

### Post-MVP (Auth + Write Operations)
- [ ] GitHub OAuth (`packs login`)
- [ ] `packs publish` — upload to packs.sh
- [ ] `packs star` / `packs starred`
- [ ] `packs init` / `packs validate`
- [ ] packs.sh write API
- [ ] Web UI at packs.sh
- [ ] Private registry adapters
- [ ] MCP server integration
- [ ] Download analytics

---

## 12. Telemetry & Indexing

### 12.1 Usage Telemetry
Every `packs get` pings packs.sh asynchronously (non-blocking):

```go
// Async, doesn't slow down the get
go func() {
    http.Post("https://packs.sh/api/telemetry", json, {
        "pack": "commit-message",       // or "gh:user/repo/pack"  
        "source": "registry",           // or "github"
        "version": "1.2.0",
        "cli_version": "0.1.0",
        "os": "darwin",
        "arch": "arm64",
    })
}()
```

- **Opt-out:** `packs config set telemetry false`
- **Purpose:** Download counts, popular packs, usage patterns
- **Privacy:** No user identification, just aggregate stats

### 12.2 Index Submission
Anyone can submit a GitHub pack for indexing (no auth):

```bash
$ packs submit gh:hsbacot/packs/commit-message

  ✓ Fetched pack.yaml from GitHub
  ✓ Validated pack structure
  ✓ Indexed at packs.sh
  
  Your pack is now discoverable:
    packs find commit
    packs get commit-message
```

**How it works:**
1. CLI sends submit request to packs.sh
2. packs.sh fetches `pack.yaml` + content from GitHub
3. Validates structure (name, type, description, content file exists)
4. Adds to index with `source: github`, `github_ref: user/repo/pack`
5. Periodic re-sync to catch updates (or webhook later)

**Why no auth:**
- It's a public GitHub repo anyway
- Submitting just adds to index, doesn't claim ownership
- Auth needed later for: editing metadata, claiming ownership, stars

### 12.3 The Pack Lifecycle

```
┌─────────────────────────────────────────────────────────────────┐
│  1. CREATE                                                       │
│     Push pack to GitHub: gh:yourname/packs/my-skill             │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│  2. SHARE                                                        │
│     Anyone can get it: packs get gh:yourname/packs/my-skill     │
│     (works immediately, no indexing needed)                      │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│  3. INDEX (optional)                                             │
│     Submit for discovery: packs submit gh:yourname/packs/my-skill│
│     Now appears in: packs find, TUI, packs.sh website           │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│  4. PUBLISH (post-MVP, optional)                                 │
│     Claim & get short name: packs publish                        │
│     Now: packs get my-skill (no gh: prefix needed)              │
└─────────────────────────────────────────────────────────────────┘
```

---

## 13. Security & Privacy

| Concern | Mitigation |
|---------|------------|
| Malicious packs | Review process for featured, community flagging |
| Token storage | OS keychain preferred, file permissions 600 |
| Content injection | Validate markdown, sanitize on web |
| Supply chain | Content hashing, version immutability |
| Rate limiting | Per-user limits on publish, standard API limits |

---

## 14. Design Decisions

| Question | Decision |
|----------|----------|
| Primary interface | TUI (Bubble Tea), CLI secondary |
| Distribution | `brew install packs` |
| Registry backend | packs.sh with Postgres + S3 |
| GitHub direct | `gh:user/repo/pack` via `gh` CLI or raw fetch |
| Auth for MVP | None — read-only + GitHub direct |
| Index submission | `packs submit` — no auth needed |
| Telemetry | Opt-out, async, anonymous |
| Pack format | Directory with `pack.yaml` + content file |
| Versioning | SemVer, default to latest |
| Namespace | Flat, unique names |

---

## 15. Success Criteria

### Launch (Week 6)
- [ ] `brew install packs` works
- [ ] TUI is fast and intuitive
- [ ] 50+ quality packs seeded
- [ ] Docs complete

### Month 1
- [ ] 500 installs
- [ ] 10 community-contributed packs
- [ ] Featured on Hacker News

### Month 3
- [ ] 20 community publishers
- [ ] 1 enterprise using private registry
- [ ] Default skill source for at least one agent ecosystem

---

*Built with 🎒 by Hunter & Fred · packs.sh*
