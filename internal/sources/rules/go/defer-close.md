---
type: rule
description: Defer Close() immediately after error check, not before
globs: ["*.go"]
---

## Good

```go
f, err := os.Open(path)
if err != nil {
    return nil, err
}
defer f.Close()
```

## Bad

```go
f, err := os.Open(path)
defer f.Close() // panic if f is nil
if err != nil {
    return nil, err
}
```
