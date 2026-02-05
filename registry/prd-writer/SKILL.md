---
name: prd-writer
description: |
  Collaborative PRD creation through guided conversation. Brainstorms user flows,
  writes testable user stories, defines scope. Use BEFORE any building.
  Proposes options at each step — human picks, combines, or overrides.
---

# PRD Writer

Create thorough Product Requirements Documents through structured collaboration.

**Philosophy:** I propose, you decide. Every question comes with options AND room for your own answer.

---

## When to Use

- Starting a new product/feature
- Before using `web-stack` or `react-native-stack` skills
- When user stories are missing or weak
- When scope is unclear

**Output:** Complete `prd.md` file with testable user stories

---

## The Process

```
┌──────────┐   ┌──────────┐   ┌──────────┐   ┌──────────┐   ┌──────────┐
│ Problem  │──▶│ Personas │──▶│  Flows   │──▶│ Stories  │──▶│  Scope   │
│          │   │          │   │          │   │          │   │          │
│ What/Who │   │ 1-3 max  │   │ All of   │   │ Testable │   │ P0/P1/   │
│ /Today?  │   │          │   │ them     │   │ criteria │   │ Out      │
└──────────┘   └──────────┘   └──────────┘   └──────────┘   └──────────┘
```

---

## Phase 1: Problem Definition

### Questions to Ask

**1.1 What's the core problem?**

Present options based on context (if idea was mentioned) or ask open:

> **What problem are we solving?**
>
> A) [Inferred option based on context]
> B) [Alternative framing]
> C) [Broader version]
> D) [Narrower version]
> E) [Your take]
>
> Pick, combine, or write your own.

**1.2 Who has this problem?**

> **Who experiences this pain most acutely?**
>
> A) Solo founders / indie hackers
> B) Small teams (2-10 people)
> C) SMB companies (10-100)
> D) Specific role (marketers, developers, ops)
> E) [Your take]

**1.3 How are they solving it today?**

> **Current solutions (what are we replacing)?**
>
> A) [Major competitor 1]
> B) [Major competitor 2]
> C) Manual process / spreadsheets
> D) Not solving it (living with pain)
> E) [Your take]

**1.4 Why do current solutions fail?**

> **What's wrong with existing options?**
>
> A) Too expensive
> B) Too complex / feature bloat
> C) Missing key capability
> D) Bad UX / friction
> E) Wrong pricing model
> F) [Your take]

### Phase 1 Output

Synthesize into problem statement:

```markdown
## Problem Statement

**The Problem:** [One sentence]

**Who Has It:** [Specific persona description]

**Current Solutions:** [What they use today]

**Why They Fail:** [The gap we're filling]
```

---

## Phase 2: Personas

### Propose Personas

Based on Phase 1, propose 2-3 personas:

> **Proposed Personas:**
>
> **A) [Name] the [Role]**
> - Context: [Their situation]
> - Goal: [What they want]
> - Pain: [What frustrates them]
> - Behavior: [How they work]
>
> **B) [Name] the [Role]**
> - Context: ...
> - Goal: ...
> - Pain: ...
> - Behavior: ...
>
> **C) [Name] the [Role]** (optional)
> - ...
>
> Keep all? Cut one? Add someone I missed?

### Persona Rules

- **Max 3 personas** — more = unfocused product
- **One primary** — who we optimize for
- **Be specific** — "Sarah the solo founder" not "users"

### Phase 2 Output

```markdown
## Personas

### Primary: [Name] the [Role]
- **Context:** [Situation]
- **Goal:** [What they want to achieve]
- **Pain:** [Current frustration]
- **Behavior:** [How they work today]

### Secondary: [Name] the [Role]
- ...
```

---

## Phase 3: User Flows

### Brainstorm ALL Flows

This is the deep work. Cover:

**Core Flows (Happy Path)**
- First-time user experience
- Primary value action (the thing they came for)
- Return user experience

**Edge Cases**
- What if X is empty?
- What if X fails?
- What if X is huge?
- What if user does Y before X?

**Error States**
- Network failure
- Validation errors
- Permission denied
- Rate limits

**Empty States**
- No data yet
- No results found
- Account just created

### Flow Brainstorm Format

> **I've mapped out these flows:**
>
> **Core Flows:**
> 1. New user signup → onboarding → first value
> 2. [Primary action] — the main thing
> 3. Return user → pick up where left off
>
> **Edge Cases:**
> 4. [Edge case 1]
> 5. [Edge case 2]
>
> **Error States:**
> 6. Network fails mid-action
> 7. Validation errors on form
>
> **Empty States:**
> 8. Dashboard with no data
> 9. Search with no results
>
> **What am I missing?** Any flows specific to this problem?

### Flow Detail

For each flow, walk through:

> **Flow: [Name]**
>
> 1. User starts at [screen/state]
> 2. User does [action]
> 3. System [response]
> 4. User sees [result]
> 5. User can now [next action]
>
> **Questions:**
> - What if [edge case]?
> - Should we [design decision]?

### Phase 3 Output

```markdown
## User Flows

### Flow 1: New User Signup
1. User lands on marketing page
2. Clicks "Get Started"
3. Enters email/password (or OAuth)
4. Sees onboarding wizard
5. Completes setup
6. Lands on dashboard with first prompt

### Flow 2: [Primary Action]
1. ...
2. ...

### Flow 3: Return User
1. ...

### Edge Cases
- **Empty state:** Show [helpful prompt]
- **Large data:** Paginate at [N] items
- **Offline:** Queue actions, sync on reconnect
```

---

## Phase 4: User Stories

### Convert Flows to Stories

Each flow becomes one or more user stories:

> **From Flow 2, I've written these stories:**
>
> ---
> ### US-2.1: [Action Name]
>
> **As a** [persona]
> **I want to** [action]
> **So that** [outcome]
>
> **Acceptance Criteria:**
> - [ ] [Testable criterion 1]
> - [ ] [Testable criterion 2]
> - [ ] [Testable criterion 3]
>
> ---
>
> **Check each criterion:** Can we write a test for it? Is it specific enough?

### Acceptance Criteria Rules

**Good criteria are:**
- Testable (yes/no, not subjective)
- Specific (exact behavior, not vague)
- Independent (don't require other stories)
- User-visible (not implementation details)

**Examples:**

❌ Bad: "Form should be user-friendly"
✅ Good: "Form shows inline error when email is invalid"

❌ Bad: "Should be fast"
✅ Good: "Page loads in under 2 seconds on 3G"

❌ Bad: "Handle errors gracefully"
✅ Good: "Network error shows retry button with message"

### Story ID Convention

```
US-[Flow].[Story]

US-1.1 = Flow 1, Story 1 (e.g., signup)
US-2.1 = Flow 2, Story 1 (e.g., create item)
US-2.2 = Flow 2, Story 2 (e.g., edit item)
```

Tests will reference these IDs.

### Phase 4 Output

```markdown
## User Stories

### US-1.1: User Sign Up

**As a** new user
**I want to** create an account with email
**So that** I can access the product

**Acceptance Criteria:**
- [ ] User can enter email and password
- [ ] Password requires 8+ characters
- [ ] Invalid email shows inline error
- [ ] Successful signup redirects to onboarding
- [ ] User record created in database

---

### US-1.2: User Sign In

**As a** returning user
**I want to** sign in to my account
**So that** I can access my data

**Acceptance Criteria:**
- [ ] User can enter email and password
- [ ] Wrong password shows error (no account enumeration)
- [ ] Successful login redirects to dashboard
- [ ] Session persists across browser refresh

---

### US-2.1: [Primary Action]
...
```

---

## Phase 5: Scope & Output

### Define Priorities

> **Let's prioritize. Here's my read:**
>
> **P0 — Must ship (MVP):**
> - US-1.1, US-1.2 (auth)
> - US-2.1 (core action)
> - US-2.2 (edit)
>
> **P1 — Next version:**
> - US-3.1 (payment)
> - US-4.1 (sharing)
>
> **Out of Scope (not now):**
> - Team/collaboration features
> - Mobile app
> - API access
>
> Agree? Move anything?

### Gray Areas

Capture unresolved questions:

> **Gray areas to resolve during build:**
> - [ ] Exact onboarding flow (wizard vs progressive?)
> - [ ] Free tier limits?
> - [ ] Data export format?

### Generate Final PRD

After all phases, output the complete document:

```markdown
# [Product Name] — PRD

*Generated: YYYY-MM-DD*

## Problem Statement

**The Problem:** [One sentence]

**Who Has It:** [Persona summary]

**Current Solutions:** [Competitors/alternatives]

**Why They Fail:** [Our opportunity]

---

## Personas

### Primary: [Name] the [Role]
...

### Secondary: [Name] the [Role]
...

---

## User Flows

### Flow 1: [Name]
...

### Flow 2: [Name]
...

---

## User Stories

### US-1.1: [Name]
...

### US-2.1: [Name]
...

---

## Scope

### P0 — Must Ship (MVP)
- US-1.1, US-1.2: Authentication
- US-2.1, US-2.2: Core functionality

### P1 — Next Version
- US-3.x: Payment
- US-4.x: Advanced features

### Out of Scope
- [Explicit list of what we're NOT doing]

---

## Gray Areas
- [ ] [Unresolved question 1]
- [ ] [Unresolved question 2]

---

## Appendix: Acceptance Criteria Summary

| Story | Criteria Count | Testable |
|-------|---------------|----------|
| US-1.1 | 5 | ✅ |
| US-1.2 | 4 | ✅ |
| US-2.1 | 6 | ✅ |
...
```

---

## Interaction Style

### Always Propose Options

Never ask open-ended without offering starting points:

```
❌ "What should the onboarding flow look like?"

✅ "Onboarding options:
   A) Single welcome screen → dashboard
   B) 3-step wizard (profile, preferences, first action)
   C) Progressive (learn as you use)
   D) Skip entirely (power users)
   E) [Your idea]"
```

### Push Back on Vague

If an answer is vague, ask for specifics:

```
❌ Accept: "It should be easy to use"

✅ Push: "Can you give me a specific scenario? Like: 
   'A user should be able to [do X] in under [N] clicks/seconds'"
```

### Capture Decisions

Note decisions as they're made:

```
📝 Decision: Going with wizard onboarding (option B)
📝 Decision: Free tier limited to 3 projects
📝 Decision: No team features in v1
```

### Phase Gates

Don't move forward until current phase is solid:

```
"Before we move to user stories, let me confirm the flows:
1. ✅ Signup/onboarding
2. ✅ Create project
3. ✅ Edit project
4. ❓ Delete project — should this exist in v1?
5. ✅ View dashboard

Resolve #4 and we'll continue."
```

---

## Quick Start

When user says "let's write a PRD" or "I have an idea":

1. **Ask for the idea** (if not provided)
2. **Start Phase 1** — Problem questions with options
3. **Work through each phase** — Don't skip
4. **Output final PRD** to `prd.md` in repo

---

*PRDs are the foundation. Tests come from stories. Stories come from flows. Flows come from problems. Start here.*
