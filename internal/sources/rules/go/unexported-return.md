---
type: rule
description: Do not return unexported types from exported functions
globs: ["*.go"]
---

## Good

```go
type User struct {
    ID   int64
    Name string
}

func NewUser(name string) *User
```

## Bad

```go
type user struct {
    ID   int64
    Name string
}

func NewUser(name string) *user
```
