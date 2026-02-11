# Testing Strategy

When to write unit tests, integration tests, and E2E tests. Make testing decisions fast.

## The Testing Pyramid

```
       /\
      /E2E\      ← Few, slow, expensive
     /──────\
    /  Integ  \   ← Medium count, medium speed
   /────────────\
  /    Unit      \ ← Many, fast, cheap
 ──────────────────
```

**Rule of thumb:** If you're unsure, default to the lower level that can catch the bug.

---

## When to Unit Test

**Always unit test:**
- Pure functions with complex logic
- Data transformations
- Validation logic
- Calculations and algorithms
- Edge cases you've already been burned by

**Skip unit tests for:**
- Simple getters/setters
- Thin wrappers around libraries
- Code that will be caught by integration tests anyway

### Unit Test Checklist
```
Is it a pure function? → Unit test it
Does it have branching logic? → Unit test it
Has it caused bugs before? → Unit test it
Is it just glue code? → Skip
```

### Example: When to Unit Test

```typescript
// ✅ Unit test this - complex transformation
function formatPrice(cents: number, currency: string): string {
  const dollars = cents / 100;
  return new Intl.NumberFormat('en-US', {
    style: 'currency',
    currency,
  }).format(dollars);
}

// ❌ Skip unit test - thin wrapper
function getUser(id: string) {
  return db.users.findUnique({ where: { id } });
}
```

---

## When to Integration Test

**Always integration test:**
- API endpoints
- Database queries (complex ones)
- Multi-service workflows
- Authentication flows
- Webhook handlers

**Skip integration tests for:**
- Simple CRUD that E2E will cover
- Third-party API calls (mock them)

### Integration Test Checklist
```
Does it cross a network boundary? → Integration test
Does it touch a database? → Integration test (or E2E)
Does it involve auth? → Integration test
Is it a single function? → Unit test instead
```

### Example: When to Integration Test

```typescript
// ✅ Integration test - API endpoint with DB
describe('POST /api/users', () => {
  it('creates user and sends welcome email', async () => {
    const res = await request(app)
      .post('/api/users')
      .send({ email: 'test@example.com' });
    
    expect(res.status).toBe(201);
    expect(await db.users.count()).toBe(1);
    expect(mockEmailService.sent).toHaveLength(1);
  });
});
```

---

## When to E2E Test

**Always E2E test:**
- Critical user flows (signup, checkout, core features)
- Flows that have broken in production
- Anything that touches multiple services

**Skip E2E tests for:**
- Every possible variation (use integration tests)
- Edge cases (use unit tests)
- Features still in flux

### E2E Test Checklist
```
Is it the happy path? → E2E test
Would it break revenue if it failed? → E2E test
Is it a rare edge case? → Unit/integration instead
Is the feature changing weekly? → Wait until stable
```

### Example: When to E2E Test

```typescript
// ✅ E2E test - critical user flow
test('user can sign up and make first purchase', async ({ page }) => {
  await page.goto('/signup');
  await page.fill('[name=email]', 'new@user.com');
  await page.fill('[name=password]', 'password123');
  await page.click('button[type=submit]');
  
  await expect(page).toHaveURL('/dashboard');
  
  await page.click('[data-testid=upgrade-button]');
  await page.fill('[data-testid=card-number]', '4242424242424242');
  await page.click('[data-testid=pay-button]');
  
  await expect(page.locator('[data-testid=pro-badge]')).toBeVisible();
});
```

---

## Quick Decision Tree

```
Is this a bug fix?
├─ Yes → Write test at lowest level that catches it
└─ No → Continue

Does it involve external services?
├─ Yes → Integration or E2E test
└─ No → Continue

Is it pure logic?
├─ Yes → Unit test
└─ No → Continue

Is it a critical user flow?
├─ Yes → E2E test
└─ No → Integration test (or skip if covered)
```

---

## Testing New Features

### Minimal Viable Testing (MVP approach)

For new features, start with:
1. **One E2E test for the happy path**
2. **Integration tests for each API endpoint**
3. **Unit tests for complex logic only**

Expand coverage after the feature stabilizes.

### TDD When It Makes Sense

TDD works best for:
- Well-defined requirements
- Pure functions
- Refactoring existing code

TDD is overkill for:
- Exploratory coding
- UI prototypes
- Features you'll throw away

---

## Common Patterns

### Database Tests

```typescript
// Use transactions for isolation
beforeEach(async () => {
  await db.$executeRaw`BEGIN`;
});

afterEach(async () => {
  await db.$executeRaw`ROLLBACK`;
});
```

### Mocking External Services

```typescript
// Mock at the boundary
vi.mock('./emailService', () => ({
  sendEmail: vi.fn().mockResolvedValue({ id: 'mock-id' }),
}));
```

### Time-Dependent Tests

```typescript
// Use fake timers
beforeEach(() => {
  vi.useFakeTimers();
  vi.setSystemTime(new Date('2024-01-15'));
});
```

---

## Red Flags

🚩 **E2E testing everything** — You're paying the time cost for no added confidence

🚩 **No tests at all** — You're manually testing every change

🚩 **Tests that test the mock** — Your tests pass but the code is broken

🚩 **Flaky tests** — Fix them immediately or delete them

🚩 **Slow test suite** — Run unit tests locally, E2E in CI

---

## Framework Quick Reference

| Framework | Level | Language | Best For |
|-----------|-------|----------|----------|
| Vitest | Unit/Integration | TypeScript | Vite/React projects |
| Jest | Unit/Integration | TypeScript | Node.js/React |
| Playwright | E2E | TypeScript | Cross-browser, reliable |
| Cypress | E2E | TypeScript | Fast feedback, debugging |
| pytest | Unit/Integration | Python | Python projects |
| Go testing | Unit/Integration | Go | Go projects |

---

*Write tests that give you confidence to ship. Not more, not less.*
