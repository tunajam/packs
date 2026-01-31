# Security Review

Review code for security vulnerabilities.

## When to Use

- Code review for security-sensitive features
- Auditing authentication/authorization code
- Reviewing API endpoints
- Pre-deployment security checks

## Vulnerability Checklist

### 1. Injection Attacks

#### SQL Injection
```typescript
// ❌ Vulnerable
const query = `SELECT * FROM users WHERE id = '${userId}'`;

// ✅ Parameterized
const query = await db.query('SELECT * FROM users WHERE id = ?', [userId]);
```

#### Command Injection
```typescript
// ❌ Vulnerable
exec(`convert ${filename} output.png`);

// ✅ Safe
execFile('convert', [filename, 'output.png']);
```

#### XSS (Cross-Site Scripting)
```tsx
// ❌ Vulnerable
<div dangerouslySetInnerHTML={{ __html: userInput }} />

// ✅ Safe - use text content or sanitize
<div>{userInput}</div>

// ✅ If HTML needed, sanitize first
import DOMPurify from 'dompurify';
<div dangerouslySetInnerHTML={{ __html: DOMPurify.sanitize(html) }} />
```

### 2. Authentication

- [ ] Passwords hashed with bcrypt/argon2 (not MD5/SHA1)
- [ ] Rate limiting on login attempts
- [ ] Secure session handling
- [ ] Proper logout (invalidate session)
- [ ] Password reset tokens expire
- [ ] MFA available for sensitive operations

```typescript
// ❌ Bad
const hash = md5(password);

// ✅ Good
const hash = await bcrypt.hash(password, 12);
```

### 3. Authorization

```typescript
// ❌ Missing authorization check
app.get('/users/:id', async (req, res) => {
  const user = await db.users.findById(req.params.id);
  res.json(user);
});

// ✅ Verify permission
app.get('/users/:id', async (req, res) => {
  const user = await db.users.findById(req.params.id);
  if (user.id !== req.user.id && !req.user.isAdmin) {
    return res.status(403).json({ error: 'Forbidden' });
  }
  res.json(user);
});
```

- [ ] Every endpoint checks authorization
- [ ] Can't access other users' data via ID manipulation
- [ ] Admin functions protected
- [ ] API keys scoped appropriately

### 4. Data Exposure

```typescript
// ❌ Exposing sensitive fields
return res.json(user);

// ✅ Explicit selection
return res.json({
  id: user.id,
  name: user.name,
  email: user.email,
  // NOT: password, apiKey, internalNotes
});
```

- [ ] No passwords in responses
- [ ] No API keys/tokens logged
- [ ] Sensitive data encrypted at rest
- [ ] PII handled appropriately

### 5. Secrets Management

```typescript
// ❌ Hardcoded secrets
const API_KEY = 'sk_live_abc123';

// ✅ Environment variables
const API_KEY = process.env.API_KEY;
```

- [ ] No secrets in code
- [ ] No secrets in git history
- [ ] .env files in .gitignore
- [ ] Production secrets in secure vault

### 6. CORS & Headers

```typescript
// ❌ Too permissive
app.use(cors({ origin: '*' }));

// ✅ Specific origins
app.use(cors({ 
  origin: ['https://myapp.com'],
  credentials: true,
}));
```

Security headers:
```typescript
app.use(helmet()); // Sets secure headers
```

### 7. Input Validation

```typescript
// ❌ Trust user input
const limit = req.query.limit;
db.query(`LIMIT ${limit}`);

// ✅ Validate and constrain
const limit = Math.min(Math.max(parseInt(req.query.limit) || 10, 1), 100);
```

- [ ] All inputs validated
- [ ] Type coercion applied
- [ ] Size limits enforced
- [ ] File upload restrictions

### 8. Dependencies

```bash
# Check for known vulnerabilities
npm audit
npx snyk test
```

- [ ] Dependencies up to date
- [ ] No critical vulnerabilities
- [ ] Lock file committed

## Output Format

```markdown
## 🔴 Critical

### [Location]: [Issue]
**Risk:** [Impact if exploited]
**Fix:** [How to fix]
**Example:**
```code
```

## 🟡 Medium

### ...

## 🟢 Low / Informational

### ...
```

## Quick Reference

| Issue | Look For |
|-------|----------|
| SQL Injection | String concatenation in queries |
| XSS | dangerouslySetInnerHTML, innerHTML |
| Auth bypass | Missing auth middleware |
| IDOR | Direct object references without ownership check |
| Secrets | Hardcoded API keys, passwords |
| SSRF | Unvalidated URLs in requests |
