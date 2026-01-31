# Commit Message Generator

Generate a conventional commit message for staged or unstaged changes.

## When to Use

- Before any commit
- When you need a well-formatted conventional commit message
- After making code changes that need to be committed

## Instructions

1. **Gather context:**
   - Check staged changes: `git diff --cached`
   - If nothing staged, check unstaged changes: `git diff`
   - List changed files: `git status --short`

2. **Analyze what the code actually does** - read the diff carefully to understand the purpose, not just the files touched.

3. **Stage modified tracked files only:**
   - If nothing is staged, stage only modified tracked files (shown as ` M` in status)
   - Do NOT add untracked files (shown as `??`) - those require explicit user action
   - Confirm what was staged with `git status --short`

4. **Generate a conventional commit message:**

```
<type>(<scope>): <subject>

<body>
```

5. **Conventional commit types:**
   - `feat`: New feature
   - `fix`: Bug fix
   - `refactor`: Code change that neither fixes a bug nor adds a feature
   - `docs`: Documentation only
   - `chore`: Build, config, or tooling changes
   - `test`: Adding or updating tests
   - `style`: Formatting, whitespace (no code change)
   - `perf`: Performance improvement

6. **Guidelines:**
   - **Subject:** Imperative mood, lowercase, no period, max 50 chars
   - **Scope:** Optional, indicates area affected (e.g., `feat(auth):`, `fix(api):`)
   - **Body:** Explain the "why" not the "what". Omit if the subject is self-explanatory.
   - **Breaking changes:** Start body with `BREAKING CHANGE:` if applicable
   - **No filler:** Don't pad with unnecessary context

7. **Output:**
   - Write the message to `.git/COMMIT_MSG`
   - Show the commit message in a code block for review
   - On its own line at the end, output: `git commit -F .git/COMMIT_MSG`

## Example

**Input:** Staged changes add a login form component

**Output:**
```
feat(auth): add login form component

Implements email/password form with validation.
Connects to auth API endpoint.
```

`git commit -F .git/COMMIT_MSG`
