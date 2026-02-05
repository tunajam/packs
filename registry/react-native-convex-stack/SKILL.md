---
name: react-native-stack
description: |
  Tunajam's React Native app toolkit. Expo + Convex + RevenueCat + NativeWind.
  Use for all mobile app projects. Includes setup, patterns, and rules.
---

# React Native Stack

The Tunajam mobile app toolkit. One stack, every app. **Type-safe. Test-driven.**

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
- [ ] User is redirected to home after signup
- [ ] User receives welcome push notification
```

Each acceptance criterion becomes a test assertion.

---

## 🧪 Test-Driven Development (TDD)

**We write tests FIRST. Tests come FROM user stories.**

### Maestro Setup

```bash
# Install Maestro
curl -Ls "https://get.maestro.mobile.dev" | bash

# Verify
maestro --version
```

### Test Structure

```
e2e/
├── flows/
│   ├── auth/
│   │   ├── US-1.1-signup.yaml      # Maps to user story
│   │   ├── US-1.2-signin.yaml
│   │   └── US-1.3-signout.yaml
│   ├── core/
│   │   ├── US-2.1-create-item.yaml
│   │   └── US-2.2-edit-item.yaml
│   └── payment/
│       ├── US-3.1-view-plans.yaml
│       └── US-3.2-purchase.yaml
└── helpers/
    └── auth.yaml                    # Reusable auth flow
```

### Writing Tests FROM User Stories

**User Story:**
```markdown
### US-2.1: Create Item

**As a** logged-in user
**I want to** create a new item
**So that** I can track my things

**Acceptance Criteria:**
- [ ] User can tap "Add" button from home
- [ ] User can enter item name (required)
- [ ] User can enter optional description
- [ ] User sees new item in list after save
- [ ] Empty name shows validation error
```

**Maestro Test (derived from story):**
```yaml
# e2e/flows/core/US-2.1-create-item.yaml
appId: com.tunajam.myapp
name: "US-2.1: Create Item"
---
# Prerequisites: User is logged in
- runFlow: "../helpers/auth.yaml"

# AC: User can tap "Add" button from home
- tapOn: "Add"

# AC: User can enter item name (required)
- tapOn:
    id: "input-name"
- inputText: "My Test Item"

# AC: User can enter optional description
- tapOn:
    id: "input-description"
- inputText: "A test description"

# Save the item
- tapOn: "Save"

# AC: User sees new item in list after save
- assertVisible: "My Test Item"

---
# Separate test for validation
appId: com.tunajam.myapp
name: "US-2.1: Create Item - Validation"
---
- runFlow: "../helpers/auth.yaml"
- tapOn: "Add"

# AC: Empty name shows validation error
- tapOn: "Save"
- assertVisible: "Name is required"
```

### TDD Workflow

```bash
# 1. Read the user story
# 2. Write the Maestro test FIRST
maestro test e2e/flows/core/US-2.1-create-item.yaml

# 3. Watch it fail (screens don't exist yet)
# 4. Implement the feature
# 5. Run test again — should pass
# 6. Refactor if needed

# Run all tests
maestro test e2e/flows/

# Run specific flow
maestro test e2e/flows/auth/
```

### Test Requirements

**Before any PR is merged:**
- [ ] All existing Maestro tests pass
- [ ] New user stories have corresponding tests
- [ ] Test names reference user story IDs (US-X.X)

**Before v1 launch:**
- [ ] Every user story in PRD has a passing test
- [ ] Payment flow tested (RevenueCat sandbox)
- [ ] Auth flows tested (signup, signin, signout)
- [ ] Error states tested (validation, network failure)

---

## 🔒 End-to-End Type Safety

**Type-safe from database to UI.**

```
┌─────────────┐    ┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│   Convex    │───▶│   Convex    │───▶│    Zod      │───▶│  TanStack   │
│   Schema    │    │  Functions  │    │   Schema    │    │    Form     │
│  (source)   │    │  (typed)    │    │ (validate)  │    │   (UI)      │
└─────────────┘    └─────────────┘    └─────────────┘    └─────────────┘
                    All TypeScript. All compile-time checked.
```

### Convex Schema (Source of Truth)

```ts
// convex/schema.ts
import { defineSchema, defineTable } from "convex/server";
import { v } from "convex/values";

export default defineSchema({
  users: defineTable({
    clerkId: v.string(),
    email: v.string(),
    name: v.string(),
    plan: v.union(v.literal("free"), v.literal("pro")),
  }).index("by_clerk_id", ["clerkId"]),

  items: defineTable({
    ownerId: v.id("users"),
    name: v.string(),
    description: v.optional(v.string()),
  }).index("by_owner", ["ownerId"]),
});
```

### Zod Schema (Mirrors Convex)

```ts
// lib/schemas/item.ts
import { z } from "zod";

export const createItemSchema = z.object({
  name: z.string().min(1, "Name is required").max(100),
  description: z.string().max(500).optional(),
});

export type CreateItemInput = z.infer<typeof createItemSchema>;
```

### Type-Safe Form

```tsx
// components/create-item-form.tsx
import { useForm } from "@tanstack/react-form";
import { zodValidator } from "@tanstack/zod-form-adapter";
import { useMutation } from "convex/react";
import { api } from "@/convex/_generated/api";
import { createItemSchema } from "@/lib/schemas/item";

export function CreateItemForm() {
  const createItem = useMutation(api.items.create);

  const form = useForm({
    defaultValues: { name: "", description: "" },
    onSubmit: async ({ value }) => {
      await createItem(value); // Fully typed
    },
    validatorAdapter: zodValidator(),
  });

  // ... form JSX
}
```

---

## Stack Overview

| Layer | Choice | Why |
|-------|--------|-----|
| Framework | Expo (managed) | No Xcode wrestling, cloud builds |
| Backend | Convex | Full type safety, real-time |
| Auth | Clerk via Convex | Social logins, session management |
| Payments | RevenueCat | IAP/subs handled, analytics |
| Navigation | Expo Router | File-based, like Next.js |
| Styling | NativeWind | Tailwind for RN |
| Components | react-native-reusables | shadcn for RN, copy-paste |
| Local State | Zustand | Simple, tiny, works with Convex |
| Forms | TanStack Form + Zod | Type-safe validation |
| Crash Reporting | Sentry | Industry standard, free tier |
| OTA Updates | Expo Updates | Push fixes without App Store |
| Push | Expo Push | Built-in, free, simple |
| E2E Testing | Maestro | Tests from user stories |
| Analytics | PostHog | Same as web |

---

## Project Setup

### 1. Create Expo App

```bash
bunx create-expo-app@latest my-app --template expo-template-blank-typescript
cd my-app
```

### 2. Install Core Dependencies

```bash
# Navigation
bunx expo install expo-router expo-linking expo-constants

# Styling
bun add nativewind
bun add -D tailwindcss

# State & Forms
bun add zustand @tanstack/react-form @tanstack/zod-form-adapter zod

# Backend
bun add convex
bunx convex init

# Auth (Clerk + Convex)
bun add @clerk/clerk-expo

# Payments
bun add react-native-purchases

# Crash Reporting
bunx expo install @sentry/react-native

# Push
bunx expo install expo-notifications expo-device

# OTA
bunx expo install expo-updates

# Analytics
bun add posthog-react-native
```

### 3. Add react-native-reusables

```bash
bunx @react-native-reusables/cli@latest init
bunx @react-native-reusables/cli@latest add button input card
```

### 4. Configure Tailwind

```js
// tailwind.config.js
module.exports = {
  content: ["./app/**/*.{js,jsx,ts,tsx}", "./components/**/*.{js,jsx,ts,tsx}"],
  presets: [require("nativewind/preset")],
  theme: {
    extend: {},
  },
  plugins: [],
};
```

### 5. Configure TypeScript (Strict)

```json
// tsconfig.json
{
  "extends": "expo/tsconfig.base",
  "compilerOptions": {
    "strict": true,
    "noUncheckedIndexedAccess": true,
    "paths": {
      "@/*": ["./*"]
    }
  }
}
```

### 6. Project Structure

```
my-app/
├── app/                         # Expo Router pages
│   ├── (auth)/                  # Auth-required routes
│   │   ├── _layout.tsx
│   │   └── home.tsx
│   ├── (public)/                # Public routes
│   │   ├── _layout.tsx
│   │   └── login.tsx
│   ├── _layout.tsx              # Root layout
│   └── index.tsx                # Entry point
├── components/
│   ├── ui/                      # react-native-reusables
│   └── [feature]/               # Feature components
├── convex/
│   ├── _generated/              # Convex codegen
│   ├── schema.ts                # Database schema (source of truth)
│   └── [feature].ts             # Mutations/queries
├── lib/
│   ├── schemas/                 # Zod schemas (mirror Convex)
│   ├── stores/                  # Zustand stores
│   ├── hooks/                   # Custom hooks
│   └── utils.ts                 # Utilities
├── e2e/                         # Maestro tests
│   ├── flows/                   # Tests organized by user story
│   │   ├── auth/
│   │   ├── core/
│   │   └── payment/
│   └── helpers/                 # Reusable flows
├── docs/                        # Required documentation
│   ├── prd.md                   # Product requirements
│   └── [product].md             # Vision doc
├── app.json
├── eas.json
├── tailwind.config.js
└── tsconfig.json
```

---

## Key Patterns

### Convex + Clerk Auth

```tsx
// app/_layout.tsx
import { ConvexProviderWithClerk } from "convex/react-clerk";
import { ClerkProvider, useAuth } from "@clerk/clerk-expo";
import { ConvexReactClient } from "convex/react";

const convex = new ConvexReactClient(process.env.EXPO_PUBLIC_CONVEX_URL!);

export default function RootLayout() {
  return (
    <ClerkProvider publishableKey={process.env.EXPO_PUBLIC_CLERK_KEY!}>
      <ConvexProviderWithClerk client={convex} useAuth={useAuth}>
        <Stack />
      </ConvexProviderWithClerk>
    </ClerkProvider>
  );
}
```

### RevenueCat Setup

```tsx
// lib/purchases.ts
import Purchases from "react-native-purchases";

export async function initPurchases(userId?: string) {
  Purchases.configure({
    apiKey: process.env.EXPO_PUBLIC_REVENUECAT_KEY!,
    appUserID: userId,
  });
}

export async function purchasePackage(pkg: PurchasesPackage) {
  try {
    const { customerInfo } = await Purchases.purchasePackage(pkg);
    return customerInfo.entitlements.active["premium"]?.isActive;
  } catch (e) {
    if (e.userCancelled) return false;
    throw e;
  }
}
```

### Expo Push Registration

```tsx
// lib/notifications.ts
import * as Notifications from "expo-notifications";
import * as Device from "expo-device";

export async function registerForPush() {
  if (!Device.isDevice) return null;

  const { status } = await Notifications.requestPermissionsAsync();
  if (status !== "granted") return null;

  const token = await Notifications.getExpoPushTokenAsync({
    projectId: process.env.EXPO_PUBLIC_PROJECT_ID,
  });

  return token.data;
}
```

### PostHog Analytics

```tsx
// app/_layout.tsx
import { PostHogProvider } from "posthog-react-native";

export default function RootLayout() {
  return (
    <PostHogProvider
      apiKey={process.env.EXPO_PUBLIC_POSTHOG_KEY!}
      options={{ host: process.env.EXPO_PUBLIC_POSTHOG_HOST }}
    >
      <ClerkProvider>
        {/* ... */}
      </ClerkProvider>
    </PostHogProvider>
  );
}
```

---

## Environment Variables

```bash
# .env.local.example

# Clerk
EXPO_PUBLIC_CLERK_PUBLISHABLE_KEY=pk_test_...

# Convex
EXPO_PUBLIC_CONVEX_URL=https://....convex.cloud

# RevenueCat
EXPO_PUBLIC_REVENUECAT_KEY=...

# PostHog
EXPO_PUBLIC_POSTHOG_KEY=phc_...
EXPO_PUBLIC_POSTHOG_HOST=https://us.i.posthog.com

# Expo
EXPO_PUBLIC_PROJECT_ID=...

# Sentry
SENTRY_DSN=https://...@sentry.io/...
```

---

## Commands

```bash
# Development
bunx expo start              # Start dev server
bunx convex dev              # Convex dev (separate terminal)

# Testing (TDD)
maestro test e2e/flows/      # Run all E2E tests
maestro test e2e/flows/auth/ # Run auth tests only
maestro studio               # Interactive test builder

# Build
eas build --platform ios
eas build --platform android
eas build --platform all

# Submit
eas submit --platform ios
eas submit --platform android

# OTA Update
eas update --branch production --message "fix: bug"

# Components
bunx @react-native-reusables/cli add button
```

---

## Rules

### 1. Requirements Before Code

No PRD = No code. Period.

If requirements don't exist:
1. Stop
2. Ask for user stories
3. Write tests from stories
4. Then implement

### 2. Tests From User Stories

- Every test file references a user story ID
- Test assertions map to acceptance criteria
- Tests run against real device/simulator
- No mocking Convex — test the real network

### 3. Type Safety is Non-Negotiable

- No `any` types
- Zod schemas mirror Convex schemas
- TypeScript strict mode always
- Fix type errors, don't suppress

### 4. RevenueCat for All Payments

Never use raw StoreKit/Play Billing. RevenueCat handles:
- Receipt validation
- Subscription status
- Analytics
- Cross-platform

### 5. Sentry from Day One

```bash
bunx @sentry/wizard@latest -i reactNative
```

Don't wait for bugs to add crash reporting.

---

## Checklist: New App Project

### Prerequisites (MUST EXIST FIRST)
- [ ] PRD with user stories and acceptance criteria
- [ ] Vision doc (`[product].md`)
- [ ] Core user flows defined

### Setup
- [ ] Create Expo app with TypeScript
- [ ] TypeScript strict mode enabled
- [ ] Convex init (new project under tunajam)
- [ ] Clerk setup (dev + prod instances)
- [ ] NativeWind + Tailwind config
- [ ] react-native-reusables init

### Testing (BEFORE IMPLEMENTATION)
- [ ] E2E test structure created (`e2e/flows/`)
- [ ] Tests written for each user story
- [ ] Tests reference story IDs (US-X.X)

### Integration
- [ ] RevenueCat setup (App Store Connect + Play Console)
- [ ] PostHog RN SDK
- [ ] Sentry wizard
- [ ] Expo Push config
- [ ] EAS Build config

### Deployment
- [ ] GitHub repo (private, tunajam org)
- [ ] Fast Review workflow
- [ ] MC project + tasks

---

## Costs

| Service | Free Tier | Paid |
|---------|-----------|------|
| Expo | 30 builds/mo | $99/mo |
| Convex | 1M calls/mo | $25/mo |
| RevenueCat | $2.5k MTR | 1% after |
| Sentry | 5K errors/mo | $26/mo |
| PostHog | 1M events | $0.00045/event |
| Apple Developer | — | $99/yr |

**MVP cost: $99/yr** (just Apple). Everything else free until traction.

---

## Quick Reference

### Docs
- [Expo](https://docs.expo.dev)
- [Convex](https://docs.convex.dev)
- [react-native-reusables](https://rnr-docs.vercel.app)
- [NativeWind](https://nativewind.dev)
- [TanStack Form](https://tanstack.com/form)
- [RevenueCat RN](https://docs.revenuecat.com/docs/reactnative)
- [Maestro](https://maestro.mobile.dev)
- [PostHog RN](https://posthog.com/docs/libraries/react-native)

---

*This is the Tunajam mobile stack. Requirements first. Tests from stories. Type-safe. Ship fast.*
