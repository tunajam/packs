# Unit Testing

Write effective unit tests that catch bugs and document behavior.

## When to Use

- Adding new functions or modules
- Fixing bugs (write test first)
- Refactoring code (ensure tests pass)
- Documenting expected behavior

## Test Structure

Use AAA pattern: **Arrange, Act, Assert**

```typescript
test('calculates total with tax', () => {
  // Arrange
  const items = [{ price: 100 }, { price: 50 }];
  const taxRate = 0.1;
  
  // Act
  const result = calculateTotal(items, taxRate);
  
  // Assert
  expect(result).toBe(165);
});
```

## Naming Tests

Describe what behavior you're testing:

```typescript
describe('calculateTotal', () => {
  it('sums item prices', () => {});
  it('applies tax rate to subtotal', () => {});
  it('returns 0 for empty cart', () => {});
  it('throws for negative prices', () => {});
});
```

## What to Test

### Do Test
- Core business logic
- Edge cases (empty, null, max values)
- Error conditions
- Integration points (mocked)

### Don't Test
- Library code (trust your dependencies)
- Simple getters/setters
- Implementation details (test behavior, not how)

## Edge Cases Checklist

- Empty input ([], '', null, undefined)
- Single item
- Maximum values
- Negative numbers
- Invalid types
- Concurrent access (if relevant)

## Mocking

Mock external dependencies, not the unit under test:

```typescript
// Mock the database, test the service logic
jest.mock('./database');
import { db } from './database';
import { UserService } from './userService';

test('creates user with hashed password', async () => {
  db.save.mockResolvedValue({ id: 1 });
  
  const result = await UserService.create('user@test.com', 'pass123');
  
  expect(db.save).toHaveBeenCalledWith(
    expect.objectContaining({
      email: 'user@test.com',
      password: expect.not.stringContaining('pass123'), // hashed
    })
  );
});
```

## Testing Async Code

```typescript
// Promises
test('fetches user data', async () => {
  const user = await fetchUser(1);
  expect(user.name).toBe('Alice');
});

// Error cases
test('throws on invalid user', async () => {
  await expect(fetchUser(-1)).rejects.toThrow('Invalid ID');
});
```

## Test File Organization

```
src/
  utils/
    formatDate.ts
    formatDate.test.ts    # co-located
  
# or

src/
  utils/formatDate.ts
__tests__/
  utils/formatDate.test.ts
```

## Quick Reference

```typescript
// Matchers
expect(value).toBe(exact);
expect(value).toEqual(deepEqual);
expect(value).toBeTruthy();
expect(value).toContain(item);
expect(fn).toThrow('message');

// Setup/Teardown
beforeEach(() => {});
afterEach(() => {});
beforeAll(() => {});
afterAll(() => {});

// Async
await expect(promise).resolves.toBe(value);
await expect(promise).rejects.toThrow();

// Mocks
const mockFn = jest.fn();
mockFn.mockReturnValue(value);
mockFn.mockResolvedValue(value);
expect(mockFn).toHaveBeenCalledWith(args);
```

## Guidelines

1. **One assertion per test** (when reasonable)
2. **Tests should be independent** (no shared state)
3. **Fast tests** (mock slow dependencies)
4. **Readable tests** (clear names, simple setup)
5. **Test behavior, not implementation**
