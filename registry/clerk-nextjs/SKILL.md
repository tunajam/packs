# Clerk + Next.js App Router Integration

**Use this skill when:** Setting up Clerk authentication in a Next.js App Router project.

## CRITICAL: Current Patterns Only

**DO:**
- Use `clerkMiddleware()` from `@clerk/nextjs/server` in `proxy.ts`
- Wrap app with `<ClerkProvider>` in `app/layout.tsx`
- Import from `@clerk/nextjs` (client) or `@clerk/nextjs/server` (server)
- Use `async/await` with `auth()` from `@clerk/nextjs/server`
- Use placeholder values only (`YOUR_PUBLISHABLE_KEY`, `YOUR_SECRET_KEY`)

**NEVER:**
- Use `authMiddleware()` (deprecated)
- Reference `_app.tsx` or `pages/` directory
- Import from deprecated APIs (`withAuth`, old `currentUser`)
- Write real API keys in code or tracked files

## Quick Setup

### 1. Install
```bash
npm install @clerk/nextjs@latest
```

### 2. Environment Variables (.env.local only)
```bash
NEXT_PUBLIC_CLERK_PUBLISHABLE_KEY=YOUR_PUBLISHABLE_KEY
CLERK_SECRET_KEY=YOUR_SECRET_KEY
```

### 3. Create proxy.ts
Location: `src/proxy.ts` (or root if no src/)

```typescript
import { clerkMiddleware } from "@clerk/nextjs/server";

export default clerkMiddleware();

export const config = {
  matcher: [
    "/((?!_next|[^?]*\\.(?:html?|css|js(?!on)|jpe?g|webp|png|gif|svg|ttf|woff2?|ico|csv|docx?|xlsx?|zip|webmanifest)).*)",
    "/(api|trpc)(.*)",
  ],
};
```

### 4. Wrap with ClerkProvider (app/layout.tsx)
```typescript
import {
  ClerkProvider,
  SignInButton,
  SignUpButton,
  SignedIn,
  SignedOut,
  UserButton,
} from "@clerk/nextjs";

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <ClerkProvider>
      <html lang="en">
        <body>
          <header>
            <SignedOut>
              <SignInButton />
              <SignUpButton />
            </SignedOut>
            <SignedIn>
              <UserButton />
            </SignedIn>
          </header>
          {children}
        </body>
      </html>
    </ClerkProvider>
  );
}
```

### 5. Protect Routes (optional)
```typescript
// proxy.ts - protect specific routes
import { clerkMiddleware, createRouteMatcher } from "@clerk/nextjs/server";

const isProtectedRoute = createRouteMatcher(["/dashboard(.*)", "/saved(.*)"]);

export default clerkMiddleware(async (auth, req) => {
  if (isProtectedRoute(req)) {
    await auth.protect();
  }
});
```

### 6. Access Auth in Server Components
```typescript
import { auth } from "@clerk/nextjs/server";

export default async function Page() {
  const { userId } = await auth();
  // ...
}
```

## Environment Setup for Vercel

| Vercel Env | Clerk Instance | When Used |
|------------|----------------|-----------|
| Production | Production app | main branch deploys |
| Preview | Development app | PR branch deploys |
| Development | Development app | local dev |

**Create TWO Clerk applications per project** - one for production, one for development.

## Verification Checklist

Before deploying, verify:
- [ ] `proxy.ts` uses `clerkMiddleware()` (not `authMiddleware`)
- [ ] `<ClerkProvider>` wraps app in `app/layout.tsx`
- [ ] All imports from `@clerk/nextjs` or `@clerk/nextjs/server`
- [ ] No `pages/` directory or `_app.tsx` references
- [ ] `.env.local.example` has placeholders only
- [ ] `.gitignore` excludes `.env*`
- [ ] Real keys only in `.env.local` (not tracked)

## Production DNS Setup

When deploying to production, Clerk requires DNS records for session management and verified emails.

### Steps
1. Create production instance in Clerk Dashboard
2. Go to **Domains** page: https://dashboard.clerk.com/~/domains
3. Add your production domain
4. Copy the DNS records shown and add them to your DNS provider

### Typical DNS Records Needed
| Type | Name | Value |
|------|------|-------|
| CNAME | `clerk` | `clerk.your-app.clerk.accounts.dev` |
| CNAME | `accounts` | (provided by Clerk) |
| TXT | (verification) | (provided by Clerk) |

### Cloudflare Users
If using Cloudflare, set DNS records to **"DNS only"** mode (gray cloud, not orange). Clerk's DNS check fails if proxied.

### After DNS Setup
1. Wait for DNS propagation (can take up to 48h, usually minutes)
2. Return to Clerk Dashboard
3. Click "Deploy certificates" when all checks pass

## Reference
- Docs: https://clerk.com/docs/quickstarts/nextjs
- Production: https://clerk.com/docs/guides/development/deployment/production
- Dashboard: https://dashboard.clerk.com
- Domains: https://dashboard.clerk.com/~/domains
