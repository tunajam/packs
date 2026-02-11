# API Design

Patterns for designing clean REST APIs. Consistent, predictable, developer-friendly.

---

## URL Structure

### Resources (nouns, not verbs)

```
✅ GET  /users
✅ POST /users
✅ GET  /users/123
✅ PUT  /users/123
✅ DELETE /users/123

❌ GET /getUsers
❌ POST /createUser
❌ GET /fetchUserById
```

### Nested Resources

```
GET /users/123/orders         # Orders for user 123
GET /orders?userId=123        # Alternative (flatter)

# Go max 2 levels deep
GET /users/123/orders/456     # OK
GET /users/123/orders/456/items/789  # Too deep → /order-items/789
```

### Query Parameters

```
# Filtering
GET /orders?status=pending&minTotal=100

# Pagination
GET /orders?page=2&limit=20
GET /orders?cursor=abc123&limit=20  # Cursor-based (preferred for infinite scroll)

# Sorting
GET /orders?sort=createdAt:desc

# Field selection
GET /users?fields=id,name,email
```

---

## HTTP Methods

| Method | Use For | Idempotent |
|--------|---------|------------|
| GET | Read data | Yes |
| POST | Create resource | No |
| PUT | Replace resource | Yes |
| PATCH | Update fields | Yes |
| DELETE | Remove resource | Yes |

### When to Use POST vs PUT vs PATCH

```
POST /users           # Create new user (server assigns ID)
PUT /users/123        # Replace user 123 entirely
PATCH /users/123      # Update specific fields of user 123
```

---

## Status Codes

### Success (2xx)

| Code | When |
|------|------|
| 200 | OK — Default success |
| 201 | Created — After POST creates resource |
| 204 | No Content — DELETE success, no body |

### Client Errors (4xx)

| Code | When |
|------|------|
| 400 | Bad Request — Invalid input |
| 401 | Unauthorized — No/invalid auth |
| 403 | Forbidden — Auth valid, permission denied |
| 404 | Not Found — Resource doesn't exist |
| 409 | Conflict — Duplicate, version mismatch |
| 422 | Unprocessable — Validation failed |
| 429 | Too Many Requests — Rate limited |

### Server Errors (5xx)

| Code | When |
|------|------|
| 500 | Internal Server Error — Unexpected failure |
| 502 | Bad Gateway — Upstream service failed |
| 503 | Service Unavailable — Temporarily down |

---

## Request/Response Bodies

### Consistent Envelope (optional)

```json
{
  "data": { ... },
  "meta": {
    "page": 1,
    "totalPages": 10,
    "totalCount": 195
  }
}
```

Or keep it flat (simpler):

```json
{
  "id": "123",
  "name": "John",
  "email": "john@example.com"
}
```

### Error Responses

```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Invalid email format",
    "details": [
      {
        "field": "email",
        "message": "Must be a valid email address"
      }
    ]
  }
}
```

---

## Naming Conventions

### JSON Properties

```json
// ✅ camelCase (JavaScript convention)
{
  "userId": "123",
  "firstName": "John",
  "createdAt": "2024-01-15T10:30:00Z"
}

// ❌ snake_case (unless targeting Python/Ruby clients)
{
  "user_id": "123",
  "first_name": "John"
}
```

### URLs

```
✅ /user-profiles         # kebab-case
✅ /userProfiles          # camelCase (also acceptable)
❌ /user_profiles         # snake_case (avoid)
❌ /UserProfiles          # PascalCase (avoid)
```

### IDs

```json
// ✅ String IDs (flexible, works with UUIDs)
{ "id": "abc123" }

// ⚠️ Numeric IDs (simpler but limited)
{ "id": 123 }
```

---

## Pagination

### Offset-Based (simple)

```
GET /items?page=2&limit=20

Response:
{
  "data": [...],
  "meta": {
    "page": 2,
    "limit": 20,
    "totalCount": 195,
    "totalPages": 10
  }
}
```

### Cursor-Based (scalable)

```
GET /items?cursor=eyJpZCI6MTAwfQ&limit=20

Response:
{
  "data": [...],
  "nextCursor": "eyJpZCI6MTIwfQ",
  "hasMore": true
}
```

**Use cursor-based for:** Real-time data, large datasets, infinite scroll.

---

## Versioning

### URL Path (recommended)

```
/v1/users
/v2/users
```

### Header (alternative)

```
Accept: application/vnd.api+json; version=2
```

### Tips

- Support old versions for 6-12 months
- Document breaking changes
- Use v1, v2 — not v1.0, v1.1

---

## Authentication

### Bearer Token (JWT)

```
Authorization: Bearer eyJhbGciOiJIUzI1NiIs...
```

### API Key

```
X-API-Key: sk_live_abc123
# Or in query (less secure, for convenience)
GET /data?apiKey=sk_live_abc123
```

---

## Rate Limiting

### Headers

```
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 45
X-RateLimit-Reset: 1640995200
```

### Response (429)

```json
{
  "error": {
    "code": "RATE_LIMITED",
    "message": "Too many requests",
    "retryAfter": 60
  }
}
```

---

## Common Patterns

### Bulk Operations

```
POST /users/bulk
{
  "users": [
    { "name": "John" },
    { "name": "Jane" }
  ]
}

Response:
{
  "created": 2,
  "errors": []
}
```

### Search

```
GET /search?q=john&type=user&limit=10
```

### Actions (RPC-style)

```
POST /orders/123/cancel
POST /users/123/verify-email
POST /payments/123/refund
```

### Soft Delete

```
DELETE /users/123          # Sets deletedAt, doesn't remove
GET /users?includeDeleted=true
```

---

## Checklist

Before shipping an API:

- [ ] URLs use nouns, not verbs
- [ ] Correct HTTP methods
- [ ] Consistent status codes
- [ ] Error format documented
- [ ] Pagination for lists
- [ ] Rate limiting in place
- [ ] Auth required on write endpoints
- [ ] Input validation
- [ ] OpenAPI spec generated

---

*A good API is one developers can guess.*
