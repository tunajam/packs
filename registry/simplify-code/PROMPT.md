# Simplify Code

Simplify the given code while preserving its behavior.

## Instructions

1. **Identify complexity sources:**
   - Unnecessary nesting
   - Redundant code
   - Over-abstraction
   - Verbose patterns

2. **Apply simplifications:**
   - Use modern language features
   - Extract repeated logic
   - Flatten nested conditions
   - Use built-in methods

3. **Preserve:**
   - Original behavior
   - Edge case handling
   - Type safety

## Simplification Patterns

### Early Returns
```typescript
// Before
function process(user) {
  if (user) {
    if (user.isActive) {
      return doWork(user);
    } else {
      return null;
    }
  } else {
    return null;
  }
}

// After
function process(user) {
  if (!user || !user.isActive) return null;
  return doWork(user);
}
```

### Optional Chaining
```typescript
// Before
const city = user && user.address && user.address.city;

// After
const city = user?.address?.city;
```

### Nullish Coalescing
```typescript
// Before
const name = user.name !== null && user.name !== undefined ? user.name : 'Anonymous';

// After
const name = user.name ?? 'Anonymous';
```

### Array Methods
```typescript
// Before
const results = [];
for (let i = 0; i < items.length; i++) {
  if (items[i].active) {
    results.push(items[i].name);
  }
}

// After
const results = items.filter(x => x.active).map(x => x.name);
```

### Object Shorthand
```typescript
// Before
return { name: name, email: email, age: age };

// After
return { name, email, age };
```

### Destructuring
```typescript
// Before
const name = user.name;
const email = user.email;

// After
const { name, email } = user;
```

## Output Format

Provide:
1. The simplified code
2. Brief explanation of changes (if not obvious)

## Guidelines

- Don't over-engineer in the other direction
- Keep readable; clever ≠ simple
- Maintain type safety
- If behavior changes, call it out
- Don't remove useful comments
