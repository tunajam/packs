# Convex

Patterns and best practices for Convex backend development.

## Core Concepts

### Schema Definition

```typescript
// convex/schema.ts
import { defineSchema, defineTable } from "convex/server";
import { v } from "convex/values";

export default defineSchema({
  users: defineTable({
    email: v.string(),
    name: v.string(),
    role: v.union(v.literal("admin"), v.literal("user")),
    createdAt: v.number(),
  })
    .index("by_email", ["email"])
    .index("by_role", ["role"]),

  posts: defineTable({
    authorId: v.id("users"),
    title: v.string(),
    content: v.string(),
    published: v.boolean(),
  })
    .index("by_author", ["authorId"])
    .index("by_published", ["published"]),
});
```

### Queries (Read)

```typescript
// convex/posts.ts
import { query } from "./_generated/server";
import { v } from "convex/values";

export const list = query({
  args: {
    limit: v.optional(v.number()),
  },
  handler: async (ctx, args) => {
    const limit = args.limit ?? 10;
    return await ctx.db
      .query("posts")
      .withIndex("by_published", (q) => q.eq("published", true))
      .order("desc")
      .take(limit);
  },
});

export const get = query({
  args: { id: v.id("posts") },
  handler: async (ctx, args) => {
    return await ctx.db.get(args.id);
  },
});
```

### Mutations (Write)

```typescript
import { mutation } from "./_generated/server";
import { v } from "convex/values";

export const create = mutation({
  args: {
    title: v.string(),
    content: v.string(),
  },
  handler: async (ctx, args) => {
    const identity = await ctx.auth.getUserIdentity();
    if (!identity) throw new Error("Not authenticated");

    const user = await ctx.db
      .query("users")
      .withIndex("by_email", (q) => q.eq("email", identity.email))
      .unique();

    return await ctx.db.insert("posts", {
      authorId: user._id,
      title: args.title,
      content: args.content,
      published: false,
    });
  },
});
```

### Actions (External APIs)

```typescript
import { action } from "./_generated/server";
import { v } from "convex/values";
import { api } from "./_generated/api";

export const sendEmail = action({
  args: { to: v.string(), subject: v.string(), body: v.string() },
  handler: async (ctx, args) => {
    // Call external API
    await fetch("https://api.sendgrid.com/v3/mail/send", {
      method: "POST",
      headers: { Authorization: `Bearer ${process.env.SENDGRID_KEY}` },
      body: JSON.stringify({ to: args.to, subject: args.subject, body: args.body }),
    });

    // Call a mutation to record the send
    await ctx.runMutation(api.emails.record, { to: args.to });
  },
});
```

## React Integration

### useQuery

```tsx
import { useQuery } from "convex/react";
import { api } from "../convex/_generated/api";

function PostList() {
  const posts = useQuery(api.posts.list, { limit: 10 });

  if (posts === undefined) return <Loading />;

  return (
    <ul>
      {posts.map((post) => (
        <li key={post._id}>{post.title}</li>
      ))}
    </ul>
  );
}
```

### useMutation

```tsx
import { useMutation } from "convex/react";
import { api } from "../convex/_generated/api";

function CreatePost() {
  const createPost = useMutation(api.posts.create);

  const handleSubmit = async (data: FormData) => {
    await createPost({
      title: data.get("title") as string,
      content: data.get("content") as string,
    });
  };

  return <form onSubmit={handleSubmit}>...</form>;
}
```

### Optimistic Updates

```tsx
const createPost = useMutation(api.posts.create).withOptimisticUpdate(
  (localStore, args) => {
    const existing = localStore.getQuery(api.posts.list, { limit: 10 });
    if (existing !== undefined) {
      localStore.setQuery(api.posts.list, { limit: 10 }, [
        { _id: "temp", ...args, published: false },
        ...existing,
      ]);
    }
  }
);
```

## Patterns

### Pagination

```typescript
export const paginated = query({
  args: {
    cursor: v.optional(v.string()),
    limit: v.number(),
  },
  handler: async (ctx, args) => {
    const results = await ctx.db
      .query("posts")
      .order("desc")
      .paginate({ cursor: args.cursor ?? null, numItems: args.limit });

    return {
      items: results.page,
      nextCursor: results.continueCursor,
      isDone: results.isDone,
    };
  },
});
```

### Relationships

```typescript
export const getWithAuthor = query({
  args: { id: v.id("posts") },
  handler: async (ctx, args) => {
    const post = await ctx.db.get(args.id);
    if (!post) return null;

    const author = await ctx.db.get(post.authorId);
    return { ...post, author };
  },
});
```

### Scheduled Functions

```typescript
// Schedule a function to run later
await ctx.scheduler.runAfter(60000, api.emails.sendReminder, { userId });

// Schedule at specific time
await ctx.scheduler.runAt(Date.now() + 86400000, api.tasks.cleanup, {});
```

## Best Practices

1. **Use indexes** for filtered queries
2. **Validate with Zod** in addition to Convex validators
3. **Keep mutations small** - single responsibility
4. **Use actions for external calls** - mutations are transactional
5. **Paginate large results**
6. **Add auth checks** at the start of handlers
