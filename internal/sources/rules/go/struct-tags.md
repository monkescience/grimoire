---
type: rule
description: Use consistent struct tag formatting with proper casing and spacing
globs: ["*.go"]
---

## Good

```go
type User struct {
    ID        int64     `json:"id" db:"id"`
    FirstName string    `json:"first_name" db:"first_name"`
    Email     string    `json:"email,omitempty"`
    Internal  string    `json:"-"`
}
```

## Bad

```go
type User struct {
    ID        int64  `json:"id"  db:"id"`
    FirstName string `json:"firstName" db:"first_name"`
    Email     string `JSON:"email"`
}
```
