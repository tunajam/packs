# Error Handling

Best practices for handling errors in applications.

## When to Use

- Building new features
- Reviewing error handling in existing code
- Debugging failures
- Designing APIs

## Core Principles

### 1. Fail Fast

Catch problems early. Validate inputs at boundaries.

```typescript
function createUser(email: string, age: number) {
  // Validate early
  if (!email.includes('@')) throw new Error('Invalid email');
  if (age < 0) throw new Error('Age must be positive');
  
  // Now proceed with confidence
  return db.users.create({ email, age });
}
```

### 2. Fail Informatively

Include context in error messages.

```typescript
// ❌ Bad
throw new Error('Failed');

// ❌ Still bad
throw new Error('User not found');

// ✅ Good
throw new Error(`User not found: ${userId}`);

// ✅ Better - custom error class
class UserNotFoundError extends Error {
  constructor(public userId: string) {
    super(`User not found: ${userId}`);
    this.name = 'UserNotFoundError';
  }
}
```

### 3. Don't Swallow Errors

```typescript
// ❌ Silent failure
try {
  await riskyOperation();
} catch (error) {
  // nothing here
}

// ✅ At minimum, log it
try {
  await riskyOperation();
} catch (error) {
  console.error('riskyOperation failed:', error);
  throw error; // re-throw if you can't handle it
}
```

## Error Types

### Custom Error Classes

```typescript
class AppError extends Error {
  constructor(
    message: string,
    public code: string,
    public statusCode: number = 500
  ) {
    super(message);
    this.name = 'AppError';
  }
}

class ValidationError extends AppError {
  constructor(message: string) {
    super(message, 'VALIDATION_ERROR', 400);
  }
}

class NotFoundError extends AppError {
  constructor(resource: string, id: string) {
    super(`${resource} not found: ${id}`, 'NOT_FOUND', 404);
  }
}

class AuthenticationError extends AppError {
  constructor() {
    super('Authentication required', 'UNAUTHORIZED', 401);
  }
}
```

## API Error Responses

```typescript
// Consistent error response shape
interface ErrorResponse {
  error: {
    code: string;
    message: string;
    details?: unknown;
  };
}

// Express error handler
app.use((err: Error, req: Request, res: Response, next: NextFunction) => {
  console.error(err);
  
  if (err instanceof AppError) {
    return res.status(err.statusCode).json({
      error: {
        code: err.code,
        message: err.message,
      },
    });
  }
  
  // Unknown error - don't leak details
  res.status(500).json({
    error: {
      code: 'INTERNAL_ERROR',
      message: 'An unexpected error occurred',
    },
  });
});
```

## Async Error Handling

### Promises

```typescript
// ❌ Unhandled rejection
fetchData().then(data => process(data));

// ✅ Handle errors
fetchData()
  .then(data => process(data))
  .catch(error => {
    console.error('Failed to fetch:', error);
    showErrorToUser('Could not load data');
  });
```

### Async/Await

```typescript
async function loadUser(id: string) {
  try {
    const user = await fetchUser(id);
    return user;
  } catch (error) {
    if (error instanceof NotFoundError) {
      return null; // Handle expected case
    }
    throw error; // Re-throw unexpected errors
  }
}
```

### Result Types

```typescript
type Result<T, E = Error> =
  | { success: true; data: T }
  | { success: false; error: E };

async function safeLoadUser(id: string): Promise<Result<User>> {
  try {
    const user = await fetchUser(id);
    return { success: true, data: user };
  } catch (error) {
    return { success: false, error: error as Error };
  }
}

// Usage
const result = await safeLoadUser('123');
if (result.success) {
  console.log(result.data);
} else {
  console.error(result.error);
}
```

## React Error Boundaries

```tsx
class ErrorBoundary extends React.Component {
  state = { hasError: false, error: null };

  static getDerivedStateFromError(error: Error) {
    return { hasError: true, error };
  }

  componentDidCatch(error: Error, info: React.ErrorInfo) {
    // Log to error tracking service
    reportError(error, info);
  }

  render() {
    if (this.state.hasError) {
      return <ErrorFallback error={this.state.error} />;
    }
    return this.props.children;
  }
}

// Usage
<ErrorBoundary>
  <App />
</ErrorBoundary>
```

## Logging

```typescript
// Structured logging
function logError(error: Error, context: Record<string, unknown>) {
  console.error(JSON.stringify({
    timestamp: new Date().toISOString(),
    level: 'error',
    message: error.message,
    stack: error.stack,
    ...context,
  }));
}

// Usage
try {
  await processOrder(orderId);
} catch (error) {
  logError(error, { orderId, userId: currentUser.id });
  throw error;
}
```

## Checklist

- [ ] Custom error classes for domain errors
- [ ] Consistent error response format
- [ ] Proper async error handling
- [ ] Error boundaries in React
- [ ] Structured logging with context
- [ ] Don't expose internal details to users
- [ ] Monitor errors in production
