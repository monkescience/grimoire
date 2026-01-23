---
type: rule
description: Single-method interfaces should use -er suffix, multi-method interfaces should use descriptive names
globs: ["*.go"]
---

## Good

```go
type Reader interface {
    Read(p []byte) (n int, err error)
}

type Validator interface {
    Validate() error
}

type UserRepository interface {
    FindUser(id int64) (*User, error)
    SaveUser(u *User) error
}
```

## Bad

```go
type IReader interface {
    Read(p []byte) (n int, err error)
}

type Validation interface {
    Validate() error
}

type UserFinderInterface interface {
    FindUser(id int64) (*User, error)
}
```
