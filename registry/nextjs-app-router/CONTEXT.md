# Next.js App Router

Patterns and conventions for the Next.js App Router (Next.js 13+).

## File Conventions

```
app/
├── layout.tsx          # Root layout (required)
├── page.tsx            # Home page (/)
├── loading.tsx         # Loading UI
├── error.tsx           # Error boundary
├── not-found.tsx       # 404 page
├── dashboard/
│   ├── page.tsx        # /dashboard
│   ├── layout.tsx      # Nested layout
│   └── [id]/
│       └── page.tsx    # /dashboard/:id
├── api/
│   └── users/
│       └── route.ts    # API route
└── (marketing)/        # Route group (no URL segment)
    ├── about/page.tsx  # /about
    └── blog/page.tsx   # /blog
```

## Server vs Client Components

### Server Components (default)

```tsx
// app/users/page.tsx
// No "use client" - runs on server only
async function UsersPage() {
  const users = await db.users.findMany(); // Direct DB access
  return <UserList users={users} />;
}
```

### Client Components

```tsx
// components/counter.tsx
'use client'; // Required for interactivity

import { useState } from 'react';

export function Counter() {
  const [count, setCount] = useState(0);
  return <button onClick={() => setCount(c => c + 1)}>{count}</button>;
}
```

### When to Use Each

| Server Component | Client Component |
|-----------------|------------------|
| Fetch data | useState, useEffect |
| Access backend directly | Event handlers |
| Sensitive data (keys, tokens) | Browser APIs |
| Large dependencies | Interactive UI |

## Data Fetching

### Server Components

```tsx
async function ProductPage({ params }: { params: { id: string } }) {
  const product = await fetch(`/api/products/${params.id}`, {
    next: { revalidate: 60 }, // ISR: revalidate every 60s
  }).then(r => r.json());
  
  return <ProductDetails product={product} />;
}
```

### Parallel Data Fetching

```tsx
async function Dashboard() {
  // Parallel fetching
  const [user, posts, analytics] = await Promise.all([
    getUser(),
    getPosts(),
    getAnalytics(),
  ]);
  
  return <DashboardView user={user} posts={posts} analytics={analytics} />;
}
```

### Server Actions

```tsx
// app/actions.ts
'use server';

export async function createPost(formData: FormData) {
  const title = formData.get('title');
  await db.posts.create({ data: { title } });
  revalidatePath('/posts');
}

// app/posts/new/page.tsx
import { createPost } from '@/app/actions';

function NewPost() {
  return (
    <form action={createPost}>
      <input name="title" />
      <button type="submit">Create</button>
    </form>
  );
}
```

## Layouts

### Root Layout

```tsx
// app/layout.tsx
export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="en">
      <body>
        <Header />
        <main>{children}</main>
        <Footer />
      </body>
    </html>
  );
}
```

### Nested Layouts

```tsx
// app/dashboard/layout.tsx
export default function DashboardLayout({ children }: { children: React.ReactNode }) {
  return (
    <div className="flex">
      <Sidebar />
      <div className="flex-1">{children}</div>
    </div>
  );
}
```

## Loading & Error States

### Loading UI

```tsx
// app/dashboard/loading.tsx
export default function Loading() {
  return <DashboardSkeleton />;
}
```

### Error Handling

```tsx
// app/dashboard/error.tsx
'use client';

export default function Error({
  error,
  reset,
}: {
  error: Error;
  reset: () => void;
}) {
  return (
    <div>
      <h2>Something went wrong!</h2>
      <button onClick={() => reset()}>Try again</button>
    </div>
  );
}
```

## API Routes

```tsx
// app/api/users/route.ts
import { NextRequest, NextResponse } from 'next/server';

export async function GET(request: NextRequest) {
  const users = await db.users.findMany();
  return NextResponse.json(users);
}

export async function POST(request: NextRequest) {
  const body = await request.json();
  const user = await db.users.create({ data: body });
  return NextResponse.json(user, { status: 201 });
}
```

### Dynamic Route Handlers

```tsx
// app/api/users/[id]/route.ts
export async function GET(
  request: NextRequest,
  { params }: { params: { id: string } }
) {
  const user = await db.users.findUnique({ where: { id: params.id } });
  if (!user) return NextResponse.json({ error: 'Not found' }, { status: 404 });
  return NextResponse.json(user);
}
```

## Metadata

```tsx
// Static metadata
export const metadata = {
  title: 'My App',
  description: 'Welcome to my app',
};

// Dynamic metadata
export async function generateMetadata({ params }: Props) {
  const product = await getProduct(params.id);
  return {
    title: product.name,
    openGraph: { images: [product.image] },
  };
}
```

## Route Handlers Best Practices

1. **Use Server Components** for data fetching when possible
2. **Cache aggressively** with `next: { revalidate: n }`
3. **Parallel fetch** with `Promise.all()`
4. **Streaming** with Suspense for progressive loading
5. **Server Actions** for mutations instead of API routes
6. **Keep Client Components small** - push interactivity to leaves
