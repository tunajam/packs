# Cursor Rules

Write effective .cursorrules files for Cursor AI. Make Cursor work the way you want.

---

## What Are Cursor Rules?

A `.cursorrules` file in your project root tells Cursor how to assist you. It's like a system prompt that's always active for that project.

```
my-project/
├── .cursorrules    ← Cursor reads this automatically
├── src/
└── package.json
```

---

## Basic Structure

```markdown
# Project: My App

## Tech Stack
- Next.js 14 (App Router)
- TypeScript
- Tailwind CSS
- Prisma + PostgreSQL

## Code Style
- Use functional components
- Prefer named exports
- Use `async/await` over `.then()`

## File Conventions
- Components in `src/components/`
- API routes in `src/app/api/`
- Database models in `prisma/schema.prisma`
```

---

## What to Include

### 1. Tech Stack

Tell Cursor exactly what you're using:

```markdown
## Stack
- Framework: Next.js 14 (App Router, NOT Pages Router)
- Language: TypeScript (strict mode)
- Styling: Tailwind CSS + shadcn/ui
- Database: Prisma with PostgreSQL
- Auth: Clerk
- State: Zustand for client, React Query for server
```

### 2. File Structure

Where things live:

```markdown
## File Structure
- `src/app/` — Next.js routes and layouts
- `src/components/` — Reusable React components
- `src/lib/` — Utility functions and shared logic
- `src/hooks/` — Custom React hooks
- `src/types/` — TypeScript type definitions
```

### 3. Coding Conventions

Your preferences:

```markdown
## Code Style
- Prefer `function` declarations over arrow functions for components
- Use early returns to reduce nesting
- Prefer `const` over `let`
- Use template literals over string concatenation
- Destructure props in function signature
```

### 4. Common Patterns

How you do things:

```markdown
## Patterns

### API Routes
Always use Route Handlers (not API Routes):
- Export named functions: GET, POST, PUT, DELETE
- Return NextResponse.json()
- Handle errors with try/catch

### Components
- Use 'use client' only when needed
- Prefer Server Components by default
- Colocate styles with components
```

### 5. Don'ts

What to avoid:

```markdown
## Avoid
- ❌ Don't use `any` — use `unknown` or proper types
- ❌ Don't use `var` — use `const` or `let`
- ❌ Don't use default exports for components
- ❌ Don't put business logic in components
- ❌ Don't use `console.log` — use proper logging
```

---

## Examples by Project Type

### Next.js + TypeScript

```markdown
# My Next.js App

## Stack
- Next.js 14 with App Router
- TypeScript (strict)
- Tailwind CSS + shadcn/ui
- Prisma + PostgreSQL
- Clerk for auth

## Structure
- `src/app/` — Routes and layouts
- `src/components/` — UI components
- `src/lib/` — Utilities and helpers
- `src/server/` — Server-only code

## Conventions
- Use Server Components by default
- Mark client components with 'use client'
- Use Zod for validation
- Use React Query for data fetching

## Patterns
- Components: functional, named exports
- API Routes: use NextResponse
- Forms: use react-hook-form + zod

## Don't
- Don't use Pages Router patterns
- Don't use getServerSideProps
- Don't use `any` type
```

### React Native + Expo

```markdown
# My Mobile App

## Stack
- Expo SDK 50
- React Native
- TypeScript
- NativeWind (Tailwind for RN)
- Zustand for state

## Structure
- `app/` — Expo Router screens
- `components/` — Reusable components
- `hooks/` — Custom hooks
- `lib/` — Utilities

## Conventions
- Use Expo Router for navigation
- Style with NativeWind classNames
- Use expo-secure-store for secrets
- Handle platform differences with Platform.select()

## Don't
- Don't use inline styles
- Don't use React Navigation directly
- Don't store secrets in AsyncStorage
```

### Python + FastAPI

```markdown
# My FastAPI Service

## Stack
- Python 3.11
- FastAPI
- SQLAlchemy + PostgreSQL
- Pydantic for validation

## Structure
- `app/` — Main application
- `app/api/` — Route handlers
- `app/models/` — SQLAlchemy models
- `app/schemas/` — Pydantic schemas
- `app/services/` — Business logic

## Conventions
- Use async def for route handlers
- Use Pydantic models for request/response
- Use dependency injection
- Type hint everything

## Don't
- Don't use raw SQL queries
- Don't put business logic in routes
- Don't use global mutable state
```

---

## Pro Tips

### Be Specific

```markdown
# ❌ Too vague
Use TypeScript

# ✅ Specific
Use TypeScript with strict mode. Prefer interfaces over types for object shapes.
Define props interfaces inline for simple components, extract for reused ones.
```

### Include Examples

```markdown
## Component Pattern

Example:
\`\`\`tsx
interface ButtonProps {
  variant: 'primary' | 'secondary';
  children: React.ReactNode;
}

export function Button({ variant, children }: ButtonProps) {
  return (
    <button className={cn(baseStyles, variantStyles[variant])}>
      {children}
    </button>
  );
}
\`\`\`
```

### Reference External Docs

```markdown
## References
- shadcn/ui: https://ui.shadcn.com/docs
- Tailwind: https://tailwindcss.com/docs
- Next.js App Router: https://nextjs.org/docs/app
```

### Keep It Updated

Your `.cursorrules` should evolve with your project. Update it when:
- You add new patterns
- You deprecate old approaches
- Team conventions change

---

## Troubleshooting

**Cursor ignoring rules?**
- Make sure file is named exactly `.cursorrules`
- File must be in project root
- Try restarting Cursor

**Rules too long?**
- Focus on the most important conventions
- Link to external docs for details
- Split into sections for readability

**Conflicting with composer?**
- Cursor rules apply to regular chat
- For Composer, use inline instructions
- Rules provide context, not commands

---

## Starter Template

Copy and customize:

```markdown
# [Project Name]

## Stack
- 

## Structure
- `src/` — 

## Conventions
- 
- 

## Patterns
### [Pattern Name]
-

## Don't
- ❌ 
- ❌ 

## References
- 
```

---

*Good cursor rules = less corrections, faster coding.*
