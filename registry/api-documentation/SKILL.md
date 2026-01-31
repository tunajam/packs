# API Documentation

Write clear, useful API documentation.

## When to Use

- Documenting new endpoints
- Improving existing API docs
- Creating API reference guides
- Generating OpenAPI specs

## Documentation Structure

### Endpoint Documentation

```markdown
## Create User

Creates a new user account.

### Request

`POST /api/users`

**Headers**
| Header | Required | Description |
|--------|----------|-------------|
| Authorization | Yes | Bearer token |
| Content-Type | Yes | application/json |

**Body**
```json
{
  "email": "user@example.com",
  "name": "Alice Smith",
  "role": "member"
}
```

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| email | string | Yes | Valid email address |
| name | string | Yes | Full name (2-100 chars) |
| role | string | No | One of: admin, member. Default: member |

### Response

**Success (201 Created)**
```json
{
  "id": "usr_abc123",
  "email": "user@example.com",
  "name": "Alice Smith",
  "role": "member",
  "createdAt": "2024-01-15T10:30:00Z"
}
```

**Errors**

| Status | Code | Description |
|--------|------|-------------|
| 400 | VALIDATION_ERROR | Invalid request body |
| 401 | UNAUTHORIZED | Missing or invalid token |
| 409 | EMAIL_EXISTS | Email already registered |

### Example

```bash
curl -X POST https://api.example.com/api/users \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"email": "alice@example.com", "name": "Alice"}'
```
```

## What to Document

### For Each Endpoint

1. **Description** - What it does, when to use it
2. **URL and method** - `POST /api/users`
3. **Authentication** - Required auth type
4. **Parameters** - Path, query, header, body
5. **Request body** - Schema with types and validation
6. **Response** - Success and error responses
7. **Examples** - curl commands, request/response pairs

### For the API Overall

1. **Overview** - What the API does
2. **Authentication** - How to get and use tokens
3. **Rate limiting** - Limits and headers
4. **Versioning** - How versions work
5. **Error handling** - Standard error format
6. **Pagination** - How to paginate results
7. **Changelog** - Breaking changes

## Guidelines

### Be Specific

```markdown
// ❌ Vague
Returns user data.

// ✅ Specific
Returns the user's profile including name, email, and role.
Does not include sensitive fields like password hash.
```

### Show Real Examples

```markdown
// ❌ Abstract
Returns a user object.

// ✅ Concrete
```json
{
  "id": "usr_abc123",
  "name": "Alice Smith",
  "email": "alice@example.com"
}
```

### Document Errors

```markdown
| Status | When |
|--------|------|
| 400 | Invalid request body or parameters |
| 401 | Missing or expired token |
| 403 | User lacks permission |
| 404 | Resource not found |
| 429 | Rate limit exceeded |
| 500 | Server error |
```

### Include Edge Cases

- What happens with empty arrays?
- What about null vs missing fields?
- Behavior at pagination boundaries?
- Rate limit exceeded response?

## OpenAPI/Swagger

```yaml
openapi: 3.0.0
info:
  title: My API
  version: 1.0.0

paths:
  /users:
    post:
      summary: Create a user
      operationId: createUser
      tags:
        - Users
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/CreateUserRequest'
      responses:
        '201':
          description: User created
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/User'
        '400':
          $ref: '#/components/responses/ValidationError'

components:
  schemas:
    User:
      type: object
      properties:
        id:
          type: string
          example: usr_abc123
        email:
          type: string
          format: email
        name:
          type: string
```

## Checklist

- [ ] Every endpoint documented
- [ ] Real example requests and responses
- [ ] All parameters listed with types
- [ ] Error responses documented
- [ ] Authentication explained
- [ ] Rate limits documented
- [ ] Curl examples provided
