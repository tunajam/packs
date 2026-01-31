# Code Review

Review code for bugs, security issues, and improvements.

## When to Use

- Before merging a PR
- When reviewing someone else's code
- For self-review before committing
- Auditing existing code for issues

## Review Checklist

### 1. Logic & Correctness
- Does the code do what it claims?
- Are edge cases handled?
- Are there off-by-one errors?
- Is error handling appropriate?
- Are there race conditions in async code?

### 2. Security
- SQL injection (parameterized queries?)
- XSS vulnerabilities (sanitized output?)
- Authentication/authorization gaps
- Sensitive data exposure
- Hardcoded secrets
- Input validation

### 3. Performance
- N+1 query problems
- Unnecessary re-renders (React)
- Missing memoization where needed
- Inefficient algorithms
- Missing indexes on queries

### 4. Maintainability
- Clear naming
- Appropriate abstractions
- DRY without over-abstraction
- Comments where non-obvious
- Consistent style

### 5. Testing
- Are new features tested?
- Do tests cover edge cases?
- Are tests readable?

## Output Format

Organize feedback by severity:

```markdown
## 🔴 Critical
Must fix before merge.

## 🟡 Suggestions
Improvements to consider.

## 💬 Nitpicks
Style/preference (optional).

## ✅ What's Good
Positive feedback.
```

## Guidelines

- Be specific: line numbers, code examples
- Explain why, not just what
- Offer alternatives when critiquing
- Acknowledge good patterns
- Ask questions for unclear code (don't assume bugs)

## Example

**Input:** React component with API call

```typescript
function UserProfile({ userId }) {
  const [user, setUser] = useState(null);
  
  useEffect(() => {
    fetch(`/api/users/${userId}`)
      .then(r => r.json())
      .then(setUser);
  }, []);
  
  return <div>{user.name}</div>;
}
```

**Output:**

## 🔴 Critical

**Line 10:** Null reference error. `user` is null on first render.
```typescript
return <div>{user?.name ?? 'Loading...'}</div>;
```

**Line 5-7:** Missing error handling. Network failures will fail silently.

## 🟡 Suggestions

**Line 6:** Missing dependency. `userId` should be in the dependency array:
```typescript
}, [userId]);
```

This also prevents stale data when userId changes.
