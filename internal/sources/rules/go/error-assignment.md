---
type: rule
description: Use plain assignment for error handling, not inline declaration in if statements
globs: ["*.go"]
---

## Good

```go
err := doSomething()
if err != nil {
    return err
}
```

## Bad

```go
if err := doSomething(); err != nil {
    return err
}
```
