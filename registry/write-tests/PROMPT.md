# Write Tests

Generate comprehensive tests for the given code.

## Instructions

1. **Identify testable units** - Functions, methods, classes

2. **Cover these cases:**
   - Happy path (normal usage)
   - Edge cases (empty, null, max values)
   - Error conditions
   - Boundary conditions

3. **Use the appropriate testing framework** (Jest/Vitest for JS/TS)

4. **Follow AAA pattern:** Arrange, Act, Assert

## Test Categories

### Happy Path
Normal, expected usage.

### Edge Cases
- Empty inputs ([], '', null, undefined)
- Single item
- Maximum values
- Boundary conditions

### Error Cases
- Invalid inputs
- Network failures (for async)
- Missing dependencies

## Output Format

```typescript
describe('[FunctionName]', () => {
  describe('happy path', () => {
    it('should [expected behavior]', () => {
      // Arrange
      // Act
      // Assert
    });
  });

  describe('edge cases', () => {
    it('should handle empty input', () => {});
    it('should handle null', () => {});
  });

  describe('error cases', () => {
    it('should throw for invalid input', () => {});
  });
});
```

## Guidelines

- One assertion per test (when reasonable)
- Descriptive test names that explain the behavior
- Don't test implementation details
- Mock external dependencies
- Include setup/teardown if needed

## Example

**Input:**
```typescript
function divide(a: number, b: number): number {
  if (b === 0) throw new Error('Division by zero');
  return a / b;
}
```

**Output:**
```typescript
describe('divide', () => {
  describe('happy path', () => {
    it('divides two positive numbers', () => {
      expect(divide(10, 2)).toBe(5);
    });

    it('divides negative numbers', () => {
      expect(divide(-10, 2)).toBe(-5);
    });

    it('returns decimal results', () => {
      expect(divide(5, 2)).toBe(2.5);
    });
  });

  describe('edge cases', () => {
    it('divides by 1 returns the number', () => {
      expect(divide(42, 1)).toBe(42);
    });

    it('dividing 0 returns 0', () => {
      expect(divide(0, 5)).toBe(0);
    });
  });

  describe('error cases', () => {
    it('throws when dividing by zero', () => {
      expect(() => divide(10, 0)).toThrow('Division by zero');
    });
  });
});
```
