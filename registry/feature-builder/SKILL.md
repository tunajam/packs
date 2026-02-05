---
name: feature-builder
description: |
  End-to-end feature development from MC task to merged PR.
  Reads task → finds PRD → extracts criteria → writes tests → implements → ships.
  Use for night shift, sub-agents, or any autonomous feature work.
---

# Feature Builder

Build features from MC tasks using TDD. Connects task management to PRD to code.

---

## When to Use

- Night shift autonomous work
- Sub-agent feature builds
- Any time you're picking up an MC task to implement

**Prerequisites:**
- MC task exists with story ID in title (e.g., "US-2.1: Dump an Idea")
- PRD exists in repo with matching story ID
- Stack skill loaded (web-stack or react-native-stack)

---

## The Flow

```
┌──────────┐    ┌──────────┐    ┌──────────┐    ┌──────────┐    ┌──────────┐
│ Pick Up  │───▶│  Find    │───▶│  Write   │───▶│  Build   │───▶│  Ship    │
│   Task   │    │   PRD    │    │  Tests   │    │ Feature  │    │   PR     │
│          │    │          │    │          │    │          │    │          │
│ MC list  │    │ Extract  │    │ From AC  │    │ TDD      │    │ Review   │
│          │    │ criteria │    │          │    │          │    │          │
└──────────┘    └──────────┘    └──────────┘    └──────────┘    └──────────┘
```

---

## Step 1: Pick Up Task

### Find the next task

```bash
# List high priority tasks
mc list -p <project> -s todo -P high

# Or in-progress tasks (resume work)
mc list -p <project> -s in-progress
```

### Start the task

```bash
mc start <task-id>
```

### Extract story ID from title

Task title format: `US-X.X: Description`

Examples:
- "US-2.1: Dump an Idea" → Story ID: `US-2.1`
- "US-3.1: Search Ideas" → Story ID: `US-3.1`

📝 **Note the story ID — you'll need it to find criteria.**

---

## Step 2: Find PRD & Extract Criteria

### Locate PRD

PRD lives in repo root:
```bash
# Check for PRD
ls -la prd.md

# Or find it
find . -name "prd.md" -type f
```

### Search for story

```bash
# Find the story section
grep -n "US-2.1" prd.md
```

### Extract acceptance criteria

Read the story section. Look for:

```markdown
### US-2.1: Dump an Idea

**As a** user with a thought
**I want to** capture it instantly
**So that** I don't lose it

**Acceptance Criteria:**
- [ ] Capture field visible at bottom of screen
- [ ] Keyboard appears when field is tapped
- [ ] User can type text and tap save (or enter)
- [ ] Idea appears in list above within 1 second
- [ ] Capture field clears after save
- [ ] Works offline (syncs later if cloud enabled)
```

📝 **Copy all acceptance criteria — each becomes a test assertion.**

---

## Step 3: Write Tests First (TDD)

### Create test file

**Web (Playwright):**
```bash
mkdir -p e2e/flows/core
touch e2e/flows/core/US-2.1-dump-idea.spec.ts
```

**Mobile (Maestro):**
```bash
mkdir -p e2e/flows/core
touch e2e/flows/core/US-2.1-dump-idea.yaml
```

### Write test from criteria

**Each acceptance criterion = test assertion**

Example (Playwright):
```typescript
// e2e/flows/core/US-2.1-dump-idea.spec.ts
import { test, expect } from "@playwright/test";

test.describe("US-2.1: Dump an Idea", () => {
  test("user can capture an idea instantly", async ({ page }) => {
    await page.goto("/");
    
    // AC: Capture field visible at bottom of screen
    const captureField = page.locator('[data-testid="capture-input"]');
    await expect(captureField).toBeVisible();
    
    // AC: Keyboard appears when field is tapped
    await captureField.click();
    await expect(captureField).toBeFocused();
    
    // AC: User can type text and tap save
    await captureField.fill("My test idea");
    await page.click('[data-testid="save-button"]');
    
    // AC: Idea appears in list above within 1 second
    await expect(page.locator("text=My test idea")).toBeVisible({ timeout: 1000 });
    
    // AC: Capture field clears after save
    await expect(captureField).toHaveValue("");
  });
});
```

Example (Maestro):
```yaml
# e2e/flows/core/US-2.1-dump-idea.yaml
appId: com.tunajam.dumpandsearch
name: "US-2.1: Dump an Idea"
---
# AC: Capture field visible at bottom of screen
- assertVisible:
    id: "capture-input"

# AC: User can type text and tap save
- tapOn:
    id: "capture-input"
- inputText: "My test idea"
- tapOn:
    id: "save-button"

# AC: Idea appears in list
- assertVisible: "My test idea"

# AC: Capture field clears after save
- assertVisible:
    id: "capture-input"
    text: ""
```

### Run test (should fail)

```bash
# Web
bun run test:e2e e2e/flows/core/US-2.1-dump-idea.spec.ts

# Mobile
maestro test e2e/flows/core/US-2.1-dump-idea.yaml
```

✅ Test fails = ready to implement

---

## Step 4: Build the Feature

### Create branch

```bash
git checkout -b feat/US-2.1-dump-idea
```

### Implement incrementally

Work through each acceptance criterion:

1. **Pick one criterion**
2. **Write minimal code to pass it**
3. **Run test**
4. **Commit when it passes**

```bash
# After each criterion passes
git add .
git commit -m "feat(capture): add capture field at bottom"

# Continue until all criteria pass
git commit -m "feat(capture): save idea to list"
git commit -m "feat(capture): clear field after save"
```

### Run full test suite

```bash
# Ensure nothing broke
bun run test:e2e  # or maestro test e2e/flows/
```

---

## Step 5: Ship PR

### Quality Check (REQUIRED)

Before creating PR, run:

```bash
bun ready
```

This runs typecheck + lint. **If it fails, fix the issues first.**

Do NOT create PR with failing types. This is enforced by:
1. This skill (you're reading it)
2. CI workflow (catches anything that slips through)

### Push branch

```bash
git push -u origin feat/US-2.1-dump-idea
```

### Create PR

Use pr-creator skill or manually:

```bash
gh pr create \
  --title "feat: US-2.1 Dump an Idea" \
  --body "## Summary
Implements US-2.1 from PRD.

## Acceptance Criteria
- [x] Capture field visible at bottom of screen
- [x] Keyboard appears when field is tapped
- [x] User can type text and tap save
- [x] Idea appears in list above within 1 second
- [x] Capture field clears after save
- [ ] Works offline (deferred to US-X.X)

## Tests
- Added e2e/flows/core/US-2.1-dump-idea.spec.ts

## Screenshots
[Add if UI change]
"
```

### Mark task done

```bash
mc done <task-id>
```

---

## Checklist Per Feature

```markdown
## Feature: US-X.X

### Setup
- [ ] Task started in MC (`mc start <id>`)
- [ ] Branch created (`feat/US-X.X-description`)
- [ ] PRD criteria extracted

### TDD
- [ ] Test file created
- [ ] All criteria converted to assertions
- [ ] Test fails initially (red)
- [ ] Feature implemented
- [ ] Test passes (green)
- [ ] Refactored if needed

### Ship
- [ ] `bun ready` passes (typecheck + lint)
- [ ] All E2E tests pass
- [ ] PR created with criteria checklist
- [ ] Fast Review passes
- [ ] Task marked done
```

---

## Night Shift Mode

When working autonomously overnight:

```
1. mc list -s todo -P high       # Find highest priority
2. mc start <task-id>            # Claim it
3. [Follow Steps 2-5]            # Build it
4. mc done <task-id>             # Complete it
5. REPEAT                        # Next task
```

### Rules for autonomous work

1. **One task at a time** — finish before starting next
2. **Stuck > 30 min?** — Pivot to next task, leave note
3. **Preview deploys only** — Don't merge to production
4. **Update MC as you go** — Status reflects reality
5. **Leave breadcrumbs** — Comments, commits, notes

---

## Sub-Agent Spawning

To delegate a feature to a sub-agent:

```
sessions_spawn with task:

"Build feature US-2.1 from dumpandsearch project.

1. Read ~/tunajam/dumpandsearch/prd.md
2. Find US-2.1 acceptance criteria
3. Follow feature-builder skill
4. Create PR when done
5. Report back with PR link"
```

---

## Troubleshooting

### Can't find story in PRD
- Check story ID matches exactly (US-2.1 not US-2.01)
- PRD might be in different location — `find . -name "prd.md"`
- Story might not be written yet — ask for requirements

### Test won't fail initially
- Likely testing something that already exists
- Check you're testing NEW functionality
- Might need to clear test data

### Feature works but test fails
- Test assertions might be too strict
- Check selectors/IDs match implementation
- Timing issues — add appropriate waits

### PR blocked by review
- Address review comments
- Re-run tests after changes
- Request re-review

---

## Quick Reference

```bash
# Task management
mc list -p <project> -s todo       # Find work
mc start <id>                       # Start task
mc done <id>                        # Complete task

# Git flow
git checkout -b feat/US-X.X-desc   # Branch
git commit -m "feat: description"   # Commit
gh pr create                        # PR

# Testing
bun run test:e2e                   # Web
maestro test e2e/flows/            # Mobile
```

---

*Requirements → Tests → Code → Ship. This is the way.*
