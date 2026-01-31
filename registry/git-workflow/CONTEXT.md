# Git Workflow

Common Git branching strategies and workflows.

## Branch Naming

```
<type>/<ticket>-<description>

feat/PROJ-123-add-login
fix/PROJ-456-broken-api
chore/PROJ-789-update-deps
refactor/PROJ-012-auth-module
```

## GitHub Flow (Simple)

Best for: Continuous deployment, small teams.

```
main ─────●─────●─────●─────●─────●─────
          │           ↑           ↑
          └─ feat/x ──┘           │
                    └─ feat/y ────┘
```

1. Branch from `main`
2. Make changes
3. Open PR
4. Review & merge
5. Deploy from `main`

## Git Flow (Structured)

Best for: Scheduled releases, larger teams.

```
main    ─────●─────────────────●─────────
             ↑                 ↑
develop ──●──┴──●──●──●──●──●──┴──●──●──
          │     ↑  ↑  ↑        ↑
feature   └─────┘  │  │        │
feature            └──┘        │
release                  └─────┘
```

Branches:
- `main`: Production releases
- `develop`: Integration branch
- `feature/*`: New features
- `release/*`: Preparing releases
- `hotfix/*`: Emergency production fixes

## Common Commands

### Starting Work
```bash
# Create feature branch
git checkout main
git pull
git checkout -b feat/PROJ-123-description

# Or from existing branch
git checkout -b fix/PROJ-456-bug develop
```

### During Development
```bash
# Stage and commit
git add -p                    # Interactive staging
git commit -m "feat: add login form"

# Sync with main
git fetch origin
git rebase origin/main        # Keep history linear
# Or merge:
git merge origin/main
```

### Preparing PR
```bash
# Squash commits (optional)
git rebase -i HEAD~3         # Interactive rebase last 3 commits

# Push branch
git push -u origin feat/PROJ-123-description

# Create PR
gh pr create --draft
```

### After Merge
```bash
# Delete local branch
git checkout main
git pull
git branch -d feat/PROJ-123-description

# Delete remote branch
git push origin --delete feat/PROJ-123-description
```

## Rebase vs Merge

### Rebase (Linear History)
```bash
git checkout feature
git rebase main
git push --force-with-lease
```

✅ Clean, linear history
❌ Rewrites history (don't use on shared branches)

### Merge (Preserve History)
```bash
git checkout feature
git merge main
```

✅ Preserves exact history
❌ Creates merge commits

## Useful Aliases

```bash
# .gitconfig
[alias]
  co = checkout
  br = branch
  ci = commit
  st = status
  lg = log --oneline --graph --decorate
  unstage = reset HEAD --
  last = log -1 HEAD
  amend = commit --amend --no-edit
```

## Recovery

### Undo Last Commit (Keep Changes)
```bash
git reset --soft HEAD~1
```

### Undo Last Commit (Discard Changes)
```bash
git reset --hard HEAD~1
```

### Recover Deleted Branch
```bash
git reflog                    # Find the commit
git checkout -b recovered-branch abc123
```

### Fix Wrong Branch
```bash
# Move commits to correct branch
git stash
git checkout correct-branch
git stash pop
```

## Best Practices

1. **Commit often** - Small, atomic commits
2. **Write good messages** - Use conventional commits
3. **Pull before push** - Avoid merge conflicts
4. **Don't commit secrets** - Use .gitignore
5. **Review your changes** - `git diff` before commit
