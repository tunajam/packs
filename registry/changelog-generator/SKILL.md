# Changelog Generator

Generate changelogs from git commit history.

## When to Use

- Preparing a release
- Updating CHANGELOG.md
- Creating release notes
- Summarizing changes since last tag

## Instructions

1. **Get commits since last tag:**
   ```bash
   git log $(git describe --tags --abbrev=0)..HEAD --oneline
   ```

2. **Parse conventional commits:**
   Group by type:
   - `feat:` → Features
   - `fix:` → Bug Fixes
   - `docs:` → Documentation
   - `refactor:` → Refactoring
   - `perf:` → Performance
   - `test:` → Tests
   - `chore:` → Maintenance

3. **Format as markdown:**

```markdown
## [1.2.0] - 2024-01-15

### Features
- Add user authentication (#42)
- Support dark mode (#45)

### Bug Fixes
- Fix login timeout issue (#43)
- Resolve memory leak in dashboard (#44)

### Breaking Changes
- Remove deprecated `v1/auth` endpoint
```

## Format Guidelines

### Header
```markdown
## [1.2.0] - 2024-01-15
```
- Version in brackets
- ISO date format

### Categories
Order by importance:
1. Breaking Changes (if any)
2. Features
3. Bug Fixes
4. Performance
5. Documentation
6. Other

### Entry Format
```markdown
- Brief description (#PR or commit)
```
- Start with verb (Add, Fix, Update, Remove)
- Link to PR or issue
- Keep under 80 chars

## Full Changelog Template

```markdown
# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/),
and this project adheres to [Semantic Versioning](https://semver.org/).

## [Unreleased]

## [1.2.0] - 2024-01-15

### Features
- Add SSO authentication support (#42)
- Enable bulk export for reports (#45)

### Bug Fixes
- Fix timezone handling in scheduler (#43)

### Breaking Changes
- Minimum Node.js version is now 18

## [1.1.0] - 2024-01-01

### Features
- Add dark mode support (#38)
```

## Automation Tips

### Git Commit Parsing

```bash
# Get commits grouped by type
git log --oneline --pretty=format:"%s" $LAST_TAG..HEAD | \
  grep -E "^(feat|fix|docs|refactor|perf|test|chore)" | \
  sort
```

### Version Bump Logic

Based on conventional commits:
- `BREAKING CHANGE:` → Major (1.0.0 → 2.0.0)
- `feat:` → Minor (1.0.0 → 1.1.0)
- `fix:` → Patch (1.0.0 → 1.0.1)

## Checklist

- [ ] Version follows semver
- [ ] Date is accurate
- [ ] All significant changes included
- [ ] Breaking changes clearly marked
- [ ] PRs/issues linked
- [ ] Entries are actionable (verb + what)
