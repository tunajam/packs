---
name: web-stack
description: |
  Tunajam's Next.js web stack. Next.js 15 + Convex + Clerk + shadcn + PostHog.
  End-to-end type safety. TDD with network-layer E2E tests.
  Use for all web projects. Includes setup, patterns, and rules.
---

# Web Stack

The Tunajam web toolkit. One stack, every app. **Type-safe. Test-driven.**

---

## ⚠️ BEFORE YOU BUILD ANYTHING

### Required Documents (STOP if missing)

Before writing ANY code, these documents MUST exist:

| Document | Purpose | Location |
|----------|---------|----------|
| `prd.md` | Product requirements with user stories | Repo root |
| `[product].md` | Vision doc | Repo root |
| User Stories | Acceptance criteria for each flow | In PRD |

**If these don't exist, STOP and ask:**
> "I need the PRD with user stories before I can build. What are the core user flows and acceptance criteria?"

### Why Requirements First?

```
PRD User Stories → E2E Tests → Implementation
     ↓                ↓              ↓
  "User can X"    Test for X    Code for X
```

Tests are derived FROM requirements. Not invented. Not guessed.

### User Story Format (Required in PRD)

```markdown
### US-1.1: User Sign Up

**As a** new user
**I want to** create an account with my email
**So that** I can access the app's features

**Acceptance Criteria:**
- [ ] User can enter email and password
- [ ] Password must be 8+ characters
- [ ] User sees error for invalid email format
- [ ] User is redirected to dashboard after signup
- [ ] User record created in database
```

Each acceptance criterion becomes a test assertion.

---

## Stack Overview

| Layer | Choice | Why |
|-------|--------|-----|
| Runtime | Bun | Speed, modern tooling |
| Framework | Next.js 15 (App Router) | SSR, routing, React 19 |
| Backend | Convex | Real-time, end-to-end types |
| Auth | Clerk | Sessions, social, orgs |
| Payments | Stripe (via Convex) | Subscriptions, webhooks |
| Styling | Tailwind CSS | Utility-first |
| Components | shadcn/ui (Nova preset) | Consistent, accessible |
| Forms | React Hook Form + Zod | Type-safe validation |
| State | Zustand (when needed) | Beyond Convex reactivity |
| Analytics | PostHog | Events, funnels, replay |
| Error Tracking | Sentry | Crash reports |
| E2E Testing | Playwright | TDD, network-layer tests |
| Hosting | Vercel | Zero-config deploys |
| Email | Resend | Transactional |

---

## 🔒 End-to-End Type Safety

**The entire data flow is type-checked at compile time.**

```
┌─────────────┐    ┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│   Convex    │───▶│   Convex    │───▶│    Zod      │───▶│    React    │
│   Schema    │    │  Functions  │    │   Schema    │    │    Form     │
│  (source)   │    │  (typed)    │    │ (validate)  │    │  (UI)       │
└─────────────┘    └─────────────┘    └─────────────┘    └─────────────┘
     │                   │                   │                   │
     └───────────────────┴───────────────────┴───────────────────┘
                    All TypeScript. All compile-time checked.
```

### The Type Chain

**1. Convex Schema (Source of Truth)**
```ts
// convex/schema.ts
import { defineSchema, defineTable } from "convex/server";
import { v } from "convex/values";

export default defineSchema({
  users: defineTable({
    clerkId: v.string(),
    email: v.string(),
    name: v.string(),
    plan: v.union(v.literal("free"), v.literal("pro"), v.literal("team")),
    createdAt: v.number(),
  }).index("by_clerk_id", ["clerkId"]),

  projects: defineTable({
    ownerId: v.id("users"),
    name: v.string(),
    description: v.optional(v.string()),
    status: v.union(v.literal("active"), v.literal("archived")),
  }).index("by_owner", ["ownerId"]),
});
```

**2. Convex Functions (Auto-typed)**
```ts
// convex/projects.ts
import { mutation, query } from "./_generated/server";
import { v } from "convex/values";

// Args and returns are fully typed from schema
export const create = mutation({
  args: {
    name: v.string(),
    description: v.optional(v.string()),
  },
  handler: async (ctx, args) => {
    const identity = await ctx.auth.getUserIdentity();
    if (!identity) throw new Error("Unauthorized");

    const user = await ctx.db
      .query("users")
      .withIndex("by_clerk_id", (q) => q.eq("clerkId", identity.subject))
      .unique();

    if (!user) throw new Error("User not found");

    // TypeScript knows exactly what fields are required
    return ctx.db.insert("projects", {
      ownerId: user._id,
      name: args.name,
      description: args.description,
      status: "active",
    });
  },
});

export const list = query({
  args: {},
  handler: async (ctx) => {
    const identity = await ctx.auth.getUserIdentity();
    if (!identity) return [];

    const user = await ctx.db
      .query("users")
      .withIndex("by_clerk_id", (q) => q.eq("clerkId", identity.subject))
      .unique();

    if (!user) return [];

    // Return type is inferred: Doc<"projects">[]
    return ctx.db
      .query("projects")
      .withIndex("by_owner", (q) => q.eq("ownerId", user._id))
      .collect();
  },
});
```

**3. Zod Schema (Runtime Validation)**
```ts
// lib/schemas/project.ts
import { z } from "zod";

// Mirror Convex schema for client-side validation
export const createProjectSchema = z.object({
  name: z.string().min(1, "Name is required").max(100),
  description: z.string().max(500).optional(),
});

export type CreateProjectInput = z.infer<typeof createProjectSchema>;
```

**4. React Hook Form (Type-safe UI)**
```tsx
// components/create-project-form.tsx
"use client";

import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { useMutation } from "convex/react";
import { api } from "@/convex/_generated/api";
import { createProjectSchema, type CreateProjectInput } from "@/lib/schemas/project";

export function CreateProjectForm() {
  const createProject = useMutation(api.projects.create);

  const form = useForm<CreateProjectInput>({
    resolver: zodResolver(createProjectSchema),
    defaultValues: {
      name: "",
      description: "",
    },
  });

  const onSubmit = async (data: CreateProjectInput) => {
    // data is typed, createProject expects the same shape
    await createProject(data);
  };

  return (
    <form onSubmit={form.handleSubmit(onSubmit)}>
      {/* Form fields - all type-checked */}
    </form>
  );
}
```

### Type Safety Rules

1. **Never use `any`** — If you need escape hatch, use `unknown` and narrow
2. **Zod mirrors Convex** — Keep validation schemas in sync with DB schema
3. **Strict mode always** — `tsconfig.json` must have `"strict": true`
4. **No type assertions** — Avoid `as` casts; fix the types instead

---

## 🧪 Test-Driven Development (TDD)

**We write tests FIRST. Tests come FROM user stories.**

### Testing Philosophy

> Unit tests don't validate user behavior. E2E tests do.

1. **Read the user story** — Understand acceptance criteria
2. **Write the test first** — Each criterion = test assertion
3. **Test through the network** — Real Convex calls, real auth
4. **Watch it fail** — Then implement
5. **Watch it pass** — Then refactor

### Test File Naming

Tests reference user story IDs:
```
e2e/
├── flows/
│   ├── auth/
│   │   ├── US-1.1-signup.spec.ts
│   │   ├── US-1.2-signin.spec.ts
│   │   └── US-1.3-signout.spec.ts
│   ├── core/
│   │   ├── US-2.1-create-project.spec.ts
│   │   └── US-2.2-edit-project.spec.ts
│   └── payment/
│       ├── US-3.1-view-pricing.spec.ts
│       └── US-3.2-checkout.spec.ts
└── fixtures/
    ├── auth.ts
    └── database.ts
```

### Playwright Setup

```bash
# Install
bun add -D @playwright/test
bunx playwright install
```

```ts
// playwright.config.ts
import { defineConfig, devices } from "@playwright/test";

export default defineConfig({
  testDir: "./e2e",
  fullyParallel: true,
  forbidOnly: !!process.env.CI,
  retries: process.env.CI ? 2 : 0,
  workers: process.env.CI ? 1 : undefined,
  reporter: "html",
  use: {
    baseURL: "http://localhost:3000",
    trace: "on-first-retry",
  },
  projects: [
    { name: "chromium", use: { ...devices["Desktop Chrome"] } },
    { name: "Mobile Safari", use: { ...devices["iPhone 13"] } },
  ],
  webServer: {
    command: "bun run dev",
    url: "http://localhost:3000",
    reuseExistingServer: !process.env.CI,
  },
});
```

### E2E Test Structure

```
e2e/
├── fixtures/
│   ├── auth.ts           # Auth helpers (Clerk test users)
│   └── database.ts       # Test data setup/teardown
├── flows/
│   ├── auth.spec.ts      # Auth flow tests
│   ├── onboarding.spec.ts
│   ├── core-action.spec.ts
│   └── payment.spec.ts
├── pages/
│   └── *.ts              # Page object models
└── global-setup.ts       # One-time setup
```

### Writing Tests FROM User Stories

**User Story (from PRD):**
```markdown
### US-2.1: Create Project

**As a** logged-in user
**I want to** create a new project
**So that** I can organize my work

**Acceptance Criteria:**
- [ ] User can click "New Project" from dashboard
- [ ] User can enter project name (required)
- [ ] User can enter optional description
- [ ] User sees new project page after save
- [ ] Project appears in dashboard list
- [ ] Empty name shows validation error
```

**Playwright Test (derived from story):**
```ts
// e2e/flows/core/US-2.1-create-project.spec.ts
import { test, expect } from "@playwright/test";
import { authenticatedPage } from "../../fixtures/auth";

test.describe("US-2.1: Create Project", () => {
  test.beforeEach(async ({ page }) => {
    await authenticatedPage(page);
  });

  test("user can create a new project", async ({ page }) => {
    // AC: User can click "New Project" from dashboard
    await page.goto("/dashboard");
    await page.click('text="New Project"');

    // AC: User can enter project name (required)
    await page.fill('input[name="name"]', "My Test Project");

    // AC: User can enter optional description
    await page.fill('textarea[name="description"]', "A test project");
    await page.click('button[type="submit"]');

    // AC: User sees new project page after save
    await expect(page).toHaveURL(/\/projects\/[a-z0-9]+/);
    await expect(page.locator("h1")).toContainText("My Test Project");

    // AC: Project appears in dashboard list
    await page.goto("/dashboard");
    await expect(page.locator("text=My Test Project")).toBeVisible();
  });

  test("validation prevents empty project name", async ({ page }) => {
    await page.goto("/projects/new");

    // Submit without filling name
    await page.click('button[type="submit"]');

    // AC: Empty name shows validation error
    await expect(page.locator("text=Name is required")).toBeVisible();
  });
});
```

**Auth Flow Test (from US-1.x stories):**
```ts
// e2e/flows/auth/US-1.1-signup.spec.ts
import { test, expect } from "@playwright/test";

test.describe("US-1.1: User Sign Up", () => {
  test("new user can sign up and land on dashboard", async ({ page }) => {
    await page.goto("/");
    
    // AC: User can click Get Started
    await page.click('text="Get Started"');

    // AC: User can enter email and password
    await page.fill('input[name="emailAddress"]', `test+${Date.now()}@example.com`);
    await page.fill('input[name="password"]', "TestPassword123!");
    await page.click('button[type="submit"]');

    // AC: User is redirected to dashboard after signup
    await expect(page).toHaveURL("/dashboard");
    await expect(page.locator("h1")).toContainText("Dashboard");
  });
});
```

**Payment Flow Test**
```ts
// e2e/flows/payment.spec.ts
import { test, expect } from "@playwright/test";
import { authenticatedPage } from "../fixtures/auth";

test.describe("Subscription", () => {
  test("user can upgrade to pro plan", async ({ page }) => {
    await authenticatedPage(page);

    // Navigate to pricing
    await page.goto("/pricing");
    await page.click('text="Upgrade to Pro"');

    // Stripe checkout (test mode)
    await expect(page).toHaveURL(/checkout\.stripe\.com/);

    // Fill test card
    await page.fill('input[name="cardNumber"]', "4242424242424242");
    await page.fill('input[name="cardExpiry"]', "12/30");
    await page.fill('input[name="cardCvc"]', "123");
    await page.fill('input[name="billingName"]', "Test User");
    await page.click('button[type="submit"]');

    // Verify success redirect
    await expect(page).toHaveURL("/dashboard?upgraded=true");
    await expect(page.locator("text=Pro Plan")).toBeVisible();
  });
});
```

### Auth Fixture

```ts
// e2e/fixtures/auth.ts
import { Page } from "@playwright/test";

// Test user credentials (Clerk test mode)
const TEST_USER = {
  email: "e2e-test@tunajam.com",
  password: process.env.E2E_TEST_PASSWORD!,
};

export async function authenticatedPage(page: Page) {
  // Check if already authenticated
  await page.goto("/dashboard");
  if (!page.url().includes("sign-in")) {
    return; // Already logged in
  }

  // Sign in
  await page.fill('input[name="identifier"]', TEST_USER.email);
  await page.fill('input[name="password"]', TEST_USER.password);
  await page.click('button[type="submit"]');

  // Wait for dashboard
  await page.waitForURL("/dashboard");
}
```

### TDD Workflow

```bash
# 1. Write the test first
bun run test:e2e e2e/flows/new-feature.spec.ts

# 2. Watch it fail (red)
# 3. Implement the feature
# 4. Watch it pass (green)
# 5. Refactor if needed

# Run all E2E tests
bun run test:e2e

# Run with UI (debugging)
bun run test:e2e --ui

# Run specific test
bun run test:e2e --grep "user can create"
```

### package.json Scripts

```json
{
  "scripts": {
    "dev": "next dev --turbo",
    "build": "next build",
    "start": "next start",
    "lint": "next lint",
    "test:e2e": "playwright test",
    "test:e2e:ui": "playwright test --ui",
    "test:e2e:debug": "playwright test --debug"
  }
}
```

### Test Requirements

**Before any PR is merged:**
- [ ] All existing E2E tests pass
- [ ] New user flows have corresponding E2E tests
- [ ] Tests run through the real network (no mocks for Convex)

**Before v1 launch:**
- [ ] Every core user flow from PRD has an E2E test
- [ ] Payment flow tested in Stripe test mode
- [ ] Auth flows tested (signup, signin, signout)
- [ ] Error states tested (network failure, validation)

---

## Project Setup

### 1. Create Next.js App

```bash
bunx create-next-app@latest my-app --typescript --tailwind --eslint --app --src-dir=false --import-alias="@/*"
cd my-app
```

### 2. Initialize shadcn/ui

```bash
bunx shadcn@latest init
# Select: Nova style, Stone base, Amber accent, Outfit font
```

Or use the preset directly:
```bash
bunx shadcn@latest init --preset "https://ui.shadcn.com/init?base=base&style=nova&baseColor=stone&theme=amber&iconLibrary=phosphor&font=outfit&menuAccent=subtle&menuColor=default&radius=large&template=next"
```

### 3. Initialize Convex

```bash
bun add convex
bunx convex init
# Creates new project under tunajam team
```

### 4. Install Dependencies

```bash
# Auth
bun add @clerk/nextjs

# Forms
bun add react-hook-form @hookform/resolvers zod

# State (if needed beyond Convex)
bun add zustand

# Analytics
bun add posthog-js

# Error tracking
bun add @sentry/nextjs

# Payments
bun add stripe @stripe/stripe-js

# E2E Testing
bun add -D @playwright/test
bunx playwright install
```

### 5. Configure TypeScript (Strict)

```json
// tsconfig.json
{
  "compilerOptions": {
    "strict": true,
    "noUncheckedIndexedAccess": true,
    "noImplicitReturns": true,
    "noFallthroughCasesInSwitch": true,
    "forceConsistentCasingInFileNames": true,
    // ... rest of Next.js defaults
  }
}
```

### 6. Project Structure

```
my-app/
├── app/
│   ├── (auth)/                  # Auth-required routes
│   │   ├── dashboard/
│   │   │   └── page.tsx
│   │   └── layout.tsx           # Auth check wrapper
│   ├── (marketing)/             # Public pages
│   │   ├── page.tsx             # Landing
│   │   ├── pricing/
│   │   └── layout.tsx
│   ├── api/
│   │   └── webhooks/
│   │       ├── clerk/route.ts
│   │       └── stripe/route.ts
│   ├── sign-in/[[...sign-in]]/
│   ├── sign-up/[[...sign-up]]/
│   ├── layout.tsx               # Root layout + providers
│   └── globals.css
├── components/
│   ├── ui/                      # shadcn components
│   ├── providers.tsx            # ConvexClerkProvider, PostHog, etc.
│   └── [feature]/               # Feature components
├── convex/
│   ├── _generated/              # Auto-generated types
│   ├── schema.ts                # Database schema (source of truth)
│   ├── users.ts                 # User mutations/queries
│   ├── [feature].ts             # Feature-specific functions
│   └── http.ts                  # HTTP endpoints for webhooks
├── lib/
│   ├── schemas/                 # Zod schemas (mirror Convex)
│   ├── hooks/                   # Custom hooks
│   ├── stores/                  # Zustand stores
│   └── utils.ts                 # Helpers
├── e2e/
│   ├── fixtures/                # Auth, database helpers
│   ├── flows/                   # User flow tests
│   └── pages/                   # Page object models
├── public/
├── .env.local
├── .env.local.example
├── playwright.config.ts
├── convex.json
├── components.json
└── next.config.ts
```

---

## Key Patterns

### Providers Setup

```tsx
// components/providers.tsx
"use client";

import { ReactNode } from "react";
import { ConvexProviderWithClerk } from "convex/react-clerk";
import { ClerkProvider, useAuth } from "@clerk/nextjs";
import { ConvexReactClient } from "convex/react";
import posthog from "posthog-js";
import { PostHogProvider } from "posthog-js/react";

const convex = new ConvexReactClient(process.env.NEXT_PUBLIC_CONVEX_URL!);

if (typeof window !== "undefined") {
  posthog.init(process.env.NEXT_PUBLIC_POSTHOG_KEY!, {
    api_host: process.env.NEXT_PUBLIC_POSTHOG_HOST,
    capture_pageview: false, // We capture manually
  });
}

export function Providers({ children }: { children: ReactNode }) {
  return (
    <ClerkProvider>
      <ConvexProviderWithClerk client={convex} useAuth={useAuth}>
        <PostHogProvider client={posthog}>
          {children}
        </PostHogProvider>
      </ConvexProviderWithClerk>
    </ClerkProvider>
  );
}
```

### Auth Layout (Protected Routes)

```tsx
// app/(auth)/layout.tsx
import { auth } from "@clerk/nextjs/server";
import { redirect } from "next/navigation";

export default async function AuthLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const { userId } = await auth();

  if (!userId) {
    redirect("/sign-in");
  }

  return <>{children}</>;
}
```

### Clerk Webhook (User Sync)

```ts
// app/api/webhooks/clerk/route.ts
import { Webhook } from "svix";
import { headers } from "next/headers";
import { WebhookEvent } from "@clerk/nextjs/server";
import { ConvexHttpClient } from "convex/browser";
import { api } from "@/convex/_generated/api";

const convex = new ConvexHttpClient(process.env.NEXT_PUBLIC_CONVEX_URL!);

export async function POST(req: Request) {
  const WEBHOOK_SECRET = process.env.CLERK_WEBHOOK_SECRET!;

  const headerPayload = headers();
  const svix_id = headerPayload.get("svix-id");
  const svix_timestamp = headerPayload.get("svix-timestamp");
  const svix_signature = headerPayload.get("svix-signature");

  if (!svix_id || !svix_timestamp || !svix_signature) {
    return new Response("Missing svix headers", { status: 400 });
  }

  const payload = await req.json();
  const body = JSON.stringify(payload);

  const wh = new Webhook(WEBHOOK_SECRET);
  let evt: WebhookEvent;

  try {
    evt = wh.verify(body, {
      "svix-id": svix_id,
      "svix-timestamp": svix_timestamp,
      "svix-signature": svix_signature,
    }) as WebhookEvent;
  } catch (err) {
    return new Response("Invalid signature", { status: 400 });
  }

  if (evt.type === "user.created") {
    await convex.mutation(api.users.create, {
      clerkId: evt.data.id,
      email: evt.data.email_addresses[0]?.email_address ?? "",
      name: `${evt.data.first_name ?? ""} ${evt.data.last_name ?? ""}`.trim(),
    });
  }

  return new Response("OK", { status: 200 });
}
```

### Stripe Checkout

```ts
// convex/stripe.ts
"use node";

import { action } from "./_generated/server";
import { v } from "convex/values";
import Stripe from "stripe";

const stripe = new Stripe(process.env.STRIPE_SECRET_KEY!);

export const createCheckoutSession = action({
  args: {
    priceId: v.string(),
    successUrl: v.string(),
    cancelUrl: v.string(),
  },
  handler: async (ctx, args) => {
    const identity = await ctx.auth.getUserIdentity();
    if (!identity) throw new Error("Unauthorized");

    const session = await stripe.checkout.sessions.create({
      mode: "subscription",
      payment_method_types: ["card"],
      line_items: [{ price: args.priceId, quantity: 1 }],
      success_url: args.successUrl,
      cancel_url: args.cancelUrl,
      client_reference_id: identity.subject,
      customer_email: identity.email,
    });

    return session.url;
  },
});
```

### Stripe Webhook

```ts
// app/api/webhooks/stripe/route.ts
import { headers } from "next/headers";
import Stripe from "stripe";
import { ConvexHttpClient } from "convex/browser";
import { api } from "@/convex/_generated/api";

const stripe = new Stripe(process.env.STRIPE_SECRET_KEY!);
const convex = new ConvexHttpClient(process.env.NEXT_PUBLIC_CONVEX_URL!);

export async function POST(req: Request) {
  const body = await req.text();
  const signature = headers().get("stripe-signature")!;

  let event: Stripe.Event;

  try {
    event = stripe.webhooks.constructEvent(
      body,
      signature,
      process.env.STRIPE_WEBHOOK_SECRET!
    );
  } catch (err) {
    return new Response("Invalid signature", { status: 400 });
  }

  switch (event.type) {
    case "checkout.session.completed": {
      const session = event.data.object as Stripe.Checkout.Session;
      await convex.mutation(api.users.upgradePlan, {
        clerkId: session.client_reference_id!,
        plan: "pro",
        stripeCustomerId: session.customer as string,
        stripeSubscriptionId: session.subscription as string,
      });
      break;
    }
    case "customer.subscription.deleted": {
      const subscription = event.data.object as Stripe.Subscription;
      await convex.mutation(api.users.downgradePlan, {
        stripeCustomerId: subscription.customer as string,
      });
      break;
    }
  }

  return new Response("OK", { status: 200 });
}
```

---

## Environment Variables

```bash
# .env.local.example

# Clerk
NEXT_PUBLIC_CLERK_PUBLISHABLE_KEY=pk_test_...
CLERK_SECRET_KEY=sk_test_...
CLERK_WEBHOOK_SECRET=whsec_...

# Convex
NEXT_PUBLIC_CONVEX_URL=https://....convex.cloud
CONVEX_DEPLOY_KEY=prod:...

# Stripe
NEXT_PUBLIC_STRIPE_PUBLISHABLE_KEY=pk_test_...
STRIPE_SECRET_KEY=sk_test_...
STRIPE_WEBHOOK_SECRET=whsec_...

# PostHog
NEXT_PUBLIC_POSTHOG_KEY=phc_...
NEXT_PUBLIC_POSTHOG_HOST=https://us.i.posthog.com

# Sentry
SENTRY_DSN=https://...@sentry.io/...

# E2E Testing
E2E_TEST_PASSWORD=... # For test user in Clerk
```

---

## Commands

```bash
# Development
bun dev                          # Next.js dev (turbo)
bunx convex dev                  # Convex dev (separate terminal)

# Testing
bun run test:e2e                 # Run all E2E tests
bun run test:e2e --ui            # Playwright UI mode
bun run test:e2e --grep "auth"   # Run specific tests

# Build & Deploy
bun run build                    # Build Next.js
bunx convex deploy               # Deploy Convex to production

# Components
bunx shadcn@latest add button    # Add shadcn component

# Database
bunx convex dashboard            # Open Convex dashboard
bunx convex data                 # View data in terminal

# Utilities
bun run lint                     # ESLint
bun run lint --fix               # Auto-fix lint issues
```

---

## Checklist: New Web Project

### Prerequisites (MUST EXIST FIRST)
- [ ] PRD with user stories and acceptance criteria
- [ ] Vision doc (`[product].md`)
- [ ] Core user flows defined
- [ ] MC project created with tasks

### Setup
- [ ] `bunx create-next-app` with TypeScript
- [ ] shadcn init (Nova preset)
- [ ] Convex init (new project under tunajam)
- [ ] Clerk setup (create dev + prod instances)
- [ ] Playwright installed + configured
- [ ] TypeScript strict mode enabled

### Testing (BEFORE IMPLEMENTATION)
- [ ] E2E test structure created (`e2e/flows/`)
- [ ] Tests written for each user story
- [ ] Tests reference story IDs (US-X.X)
- [ ] Auth fixture for test user

### Integration
- [ ] Providers component with Clerk + Convex + PostHog
- [ ] Clerk webhook syncing users to Convex
- [ ] Protected route layout
- [ ] Stripe products + checkout + webhook (if paid)
- [ ] Sentry error boundary

### Deployment
- [ ] GitHub repo created (private, tunajam org)
- [ ] Fast Review workflow added
- [ ] Shipped changelog workflow added
- [ ] Vercel project linked
- [ ] Environment variables in Vercel (prod + preview)
- [ ] Convex production deployment

### Documentation
- [ ] `[product].md` — Vision doc
- [ ] `prd.md` — Requirements with user stories
- [ ] `tech-stack.md` — Architecture reference
- [ ] `CLAUDE.md` — Agent guide
- [ ] `STORY.md` — Founder story

---

## Rules

### 1. Requirements Before Code

No PRD = No code. Period.

If requirements don't exist:
1. Stop
2. Ask for user stories with acceptance criteria
3. Write tests from stories
4. Then implement

### 2. Type Safety is Non-Negotiable

- No `any` types
- Zod schemas mirror Convex schemas
- TypeScript strict mode always on
- Fix type errors, don't suppress them

### 3. TDD for Core Flows

- Write E2E test BEFORE implementing feature
- Tests derived from user story acceptance criteria
- Test names reference story IDs (US-X.X)
- Tests hit real Convex (no mocks)
- PR blocked if tests fail

### 4. Convex is the Backend

- No API routes for data operations (use Convex)
- API routes only for webhooks
- Server actions only for Convex calls
- Real-time by default

### 5. Component Library Hierarchy

1. shadcn/ui — First choice
2. Kibo UI — When shadcn doesn't have it
3. Build custom — Last resort

### 6. Analytics from Day One

PostHog on every page:
```tsx
import { usePostHog } from "posthog-js/react";

// Track key actions
posthog.capture("project_created", { plan: "pro" });
```

---

## Costs

| Service | Free Tier | Paid |
|---------|-----------|------|
| Vercel | Hobby | $20/mo Pro |
| Convex | 1M calls/mo | $25/mo |
| Clerk | 10K MAU | $25/mo |
| Stripe | — | 2.9% + 30¢ |
| PostHog | 1M events | $0.00045/event |
| Sentry | 5K errors | $26/mo |

**MVP cost: $0** — Everything free until traction.

---

## Quick Reference

### Docs
- [Next.js](https://nextjs.org/docs)
- [Convex](https://docs.convex.dev)
- [Clerk + Convex](https://docs.convex.dev/auth/clerk)
- [shadcn/ui](https://ui.shadcn.com)
- [Playwright](https://playwright.dev)
- [React Hook Form](https://react-hook-form.com)
- [Zod](https://zod.dev)
- [Stripe](https://stripe.com/docs)

---

*This is the Tunajam web stack. Requirements first. Tests from stories. Type-safe. Ship fast.*
