---
type: rule
description: Use log/slog for logging instead of log or fmt.Print
globs: ["*.go"]
---

## Good

```go
slog.Info("server started", "port", port)
slog.Error("request failed", "error", err)
```

## Bad

```go
log.Printf("server started on port %d", port)
fmt.Println("request failed:", err)
```
