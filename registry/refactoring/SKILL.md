# Refactoring

Improve code structure without changing behavior.

## When to Use

- Code is hard to understand
- Adding features requires many changes
- Duplicated logic exists
- Tests are hard to write
- "Code smell" is detected

## The Golden Rule

**Have tests before refactoring.** If you can't verify behavior didn't change, don't refactor.

## Common Refactorings

### Extract Function
Turn a code block into a named function.

```typescript
// Before
function processOrder(order) {
  // validate
  if (!order.items.length) throw new Error('Empty order');
  if (!order.customer) throw new Error('No customer');
  // ... 50 more lines
}

// After
function processOrder(order) {
  validateOrder(order);
  // ... rest
}

function validateOrder(order) {
  if (!order.items.length) throw new Error('Empty order');
  if (!order.customer) throw new Error('No customer');
}
```

### Extract Variable
Give a name to a complex expression.

```typescript
// Before
if (user.age >= 18 && user.country === 'US' && !user.banned) {

// After
const isEligible = user.age >= 18 && user.country === 'US' && !user.banned;
if (isEligible) {
```

### Inline
Remove unnecessary indirection.

```typescript
// Before
function getFullName(user) {
  return formatName(user.first, user.last);
}
function formatName(first, last) {
  return `${first} ${last}`;
}

// After (if formatName is only used once)
function getFullName(user) {
  return `${user.first} ${user.last}`;
}
```

### Replace Conditional with Polymorphism
Use objects instead of switch statements.

```typescript
// Before
function calculatePay(employee) {
  switch (employee.type) {
    case 'hourly': return employee.hours * employee.rate;
    case 'salary': return employee.salary / 12;
    case 'commission': return employee.sales * 0.1;
  }
}

// After
const payCalculators = {
  hourly: (e) => e.hours * e.rate,
  salary: (e) => e.salary / 12,
  commission: (e) => e.sales * 0.1,
};

function calculatePay(employee) {
  return payCalculators[employee.type](employee);
}
```

### Replace Magic Number with Constant

```typescript
// Before
if (items.length > 100) {

// After
const MAX_CART_ITEMS = 100;
if (items.length > MAX_CART_ITEMS) {
```

## Code Smells

| Smell | Solution |
|-------|----------|
| Long function | Extract smaller functions |
| Long parameter list | Use options object |
| Duplicated code | Extract shared function |
| Feature envy | Move method to the class it uses |
| Data clumps | Create a class/type |
| Primitive obsession | Create domain types |
| Comments explaining what | Rename to be self-documenting |

## Refactoring Process

1. **Identify** the smell or issue
2. **Write tests** if they don't exist
3. **Make one small change**
4. **Run tests**
5. **Commit**
6. **Repeat**

## Tips

- Small steps beat big rewrites
- Commit after each successful refactor
- Don't refactor and add features in the same commit
- Use IDE refactoring tools (rename, extract, inline)
- When in doubt, leave it alone
