# Debugging

A systematic approach to finding and fixing bugs.

## When to Use

- Something isn't working as expected
- Error messages appear
- Tests are failing
- Behavior is inconsistent

## The Process

### 1. Reproduce

Before anything else, reproduce the bug reliably.

- What exact steps trigger it?
- Does it happen every time?
- What environment? (browser, OS, node version)
- Can you write a failing test?

### 2. Isolate

Narrow down where the bug lives.

```
Is it in the frontend or backend?
  → Is it in this component or its parent?
    → Is it in this function or its dependencies?
      → Is it on this line?
```

Techniques:
- Binary search through code (comment out half)
- Add logging at key points
- Check if the bug exists in older commits (`git bisect`)
- Simplify inputs to minimal reproduction

### 3. Understand

Read the code. Actually read it.

- What should happen vs. what actually happens?
- What are the inputs at the failing point?
- What assumptions is the code making?
- Check types, null values, async timing

### 4. Fix

Make the smallest change that fixes the issue.

- Fix the root cause, not the symptom
- Don't add defensive code without understanding why
- Consider if the fix could break other things

### 5. Verify

- Does the original reproduction pass?
- Do existing tests still pass?
- Add a new test for this bug

## Debugging Toolkit

### Console/Logging
```javascript
console.log('checkpoint 1', { value, state });
console.trace(); // stack trace
console.table(arrayOfObjects); // formatted table
```

### Breakpoints
- `debugger;` statement
- Browser DevTools breakpoints
- VS Code breakpoints
- Conditional breakpoints for loops

### Git Bisect
```bash
git bisect start
git bisect bad          # current commit is broken
git bisect good v1.0.0  # this tag was working
# Git checks out middle commit
# Test, then: git bisect good/bad
# Repeat until culprit found
```

### Network
- Browser Network tab
- Check request/response payloads
- Look for CORS errors
- Check status codes

## Common Gotchas

| Symptom | Likely Cause |
|---------|--------------|
| Works locally, fails in prod | Environment variables, CORS, build differences |
| Works sometimes | Race condition, timing, caching |
| Null/undefined errors | Missing null checks, async not awaited |
| Wrong data displayed | Stale state, cache, wrong variable |
| Changes don't appear | Build not running, caching, wrong file |

## When Stuck

1. Take a break (seriously)
2. Explain the problem out loud (rubber duck)
3. Check recent changes (`git diff`, `git log`)
4. Search the error message exactly
5. Ask someone to pair
