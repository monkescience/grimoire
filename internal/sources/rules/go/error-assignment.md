---
type: rule
description: Use plain assignment for error handling, not inline
tags: [go, errors, style]
globs: ["*.go"]
---

# Error Assignment

Use plain assignment for error handling:

```go
// Good
err := doSomething()
if err != nil {
    return err
}

// Bad
if err := doSomething(); err != nil {
    return err
}
```
