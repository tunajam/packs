# TypeScript Patterns

Common TypeScript patterns and type-level programming techniques.

## Utility Types

### Built-in Utilities

```typescript
// Make all properties optional
type PartialUser = Partial<User>;

// Make all properties required
type RequiredUser = Required<User>;

// Make all properties readonly
type ReadonlyUser = Readonly<User>;

// Pick specific properties
type UserPreview = Pick<User, 'id' | 'name'>;

// Omit specific properties
type UserWithoutPassword = Omit<User, 'password'>;

// Extract return type of function
type Result = ReturnType<typeof fetchUser>;

// Extract parameters of function
type Params = Parameters<typeof fetchUser>;

// Extract instance type of class
type Instance = InstanceType<typeof UserClass>;

// Make properties nullable
type NullableUser = { [K in keyof User]: User[K] | null };
```

## Discriminated Unions

```typescript
type Result<T> =
  | { success: true; data: T }
  | { success: false; error: string };

function handleResult(result: Result<User>) {
  if (result.success) {
    console.log(result.data); // TS knows data exists
  } else {
    console.log(result.error); // TS knows error exists
  }
}
```

## Type Guards

### User-defined Type Guards

```typescript
function isUser(value: unknown): value is User {
  return (
    typeof value === 'object' &&
    value !== null &&
    'id' in value &&
    'email' in value
  );
}

function process(input: unknown) {
  if (isUser(input)) {
    console.log(input.email); // TS knows it's a User
  }
}
```

### Assertion Functions

```typescript
function assertUser(value: unknown): asserts value is User {
  if (!isUser(value)) {
    throw new Error('Not a user');
  }
}

function process(input: unknown) {
  assertUser(input);
  console.log(input.email); // TS knows it's a User after assertion
}
```

## Generics

### Constrained Generics

```typescript
function getProperty<T, K extends keyof T>(obj: T, key: K): T[K] {
  return obj[key];
}

const user = { name: 'Alice', age: 30 };
const name = getProperty(user, 'name'); // type: string
```

### Default Type Parameters

```typescript
interface Container<T = string> {
  value: T;
}

const stringContainer: Container = { value: 'hello' };
const numberContainer: Container<number> = { value: 42 };
```

### Generic Constraints with Interfaces

```typescript
interface HasId {
  id: string;
}

function findById<T extends HasId>(items: T[], id: string): T | undefined {
  return items.find(item => item.id === id);
}
```

## Mapped Types

```typescript
// Make all properties optional and nullable
type DeepPartial<T> = {
  [P in keyof T]?: T[P] extends object ? DeepPartial<T[P]> : T[P];
};

// Create readonly version
type Immutable<T> = {
  readonly [P in keyof T]: T[P] extends object ? Immutable<T[P]> : T[P];
};

// Transform property types
type Getters<T> = {
  [K in keyof T as `get${Capitalize<string & K>}`]: () => T[K];
};
```

## Template Literal Types

```typescript
type EventName = 'click' | 'focus' | 'blur';
type EventHandler = `on${Capitalize<EventName>}`;
// 'onClick' | 'onFocus' | 'onBlur'

type HTTPMethod = 'GET' | 'POST' | 'PUT' | 'DELETE';
type Endpoint = `/${string}`;
type Route = `${HTTPMethod} ${Endpoint}`;
// 'GET /users' | 'POST /users' | etc.
```

## Conditional Types

```typescript
type NonNullable<T> = T extends null | undefined ? never : T;

type Flatten<T> = T extends Array<infer U> ? U : T;
// Flatten<string[]> = string
// Flatten<number> = number

type UnwrapPromise<T> = T extends Promise<infer U> ? U : T;
// UnwrapPromise<Promise<string>> = string
```

## Function Overloads

```typescript
function createElement(tag: 'input'): HTMLInputElement;
function createElement(tag: 'div'): HTMLDivElement;
function createElement(tag: string): HTMLElement;
function createElement(tag: string): HTMLElement {
  return document.createElement(tag);
}

const input = createElement('input'); // HTMLInputElement
const div = createElement('div');     // HTMLDivElement
```

## Branded Types

```typescript
type UserId = string & { readonly brand: unique symbol };
type PostId = string & { readonly brand: unique symbol };

function createUserId(id: string): UserId {
  return id as UserId;
}

function getUser(id: UserId): User { ... }

const userId = createUserId('123');
const postId = '456' as PostId;

getUser(userId); // ✅ OK
getUser(postId); // ❌ Error: PostId not assignable to UserId
```

## Zod for Runtime Validation

```typescript
import { z } from 'zod';

const UserSchema = z.object({
  id: z.string().uuid(),
  email: z.string().email(),
  age: z.number().min(0).optional(),
});

type User = z.infer<typeof UserSchema>;

function parseUser(data: unknown): User {
  return UserSchema.parse(data);
}
```

## Best Practices

1. **Prefer `unknown` over `any`** for truly unknown values
2. **Use `const` assertions** for literal types: `as const`
3. **Avoid type assertions** when possible; use type guards
4. **Enable strict mode** in tsconfig
5. **Use discriminated unions** for state machines
6. **Extract complex types** into named types for readability
