---
type: rule
description: Return interfaces from constructors to hide implementation details
tags: [go, interfaces, api]
globs: ["*.go"]
---

## Good

```go
type Runner interface {
    Run() error
}

func NewRunner() Runner {
    return &runner{}
}
```

## Bad

```go
func NewRunner() *runner {
    return &runner{}
}
```
