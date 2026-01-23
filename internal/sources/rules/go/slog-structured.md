---
type: rule
description: Use structured key-value pairs with type-safe constructors in slog
globs: ["*.go"]
---

Use type-safe attribute constructors for better performance and type checking:
- `slog.String(key, value)`
- `slog.Int(key, value)`, `slog.Int64`, `slog.Uint64`
- `slog.Float64(key, value)`
- `slog.Bool(key, value)`
- `slog.Time(key, value)`
- `slog.Duration(key, value)`
- `slog.Any(key, value)` for other types

## Good

```go
slog.Info("user logged in",
    slog.String("user_id", userID),
    slog.String("ip", remoteAddr),
    slog.Int("attempt", attemptCount))

slog.Error("request failed",
    slog.String("path", r.URL.Path),
    slog.Duration("latency", elapsed),
    slog.Any("error", err))
```

## Bad

```go
slog.Info(fmt.Sprintf("user %s logged in from %s", userID, remoteAddr))
slog.Info("user logged in", "user_id", userID, "ip", remoteAddr)  // untyped
slog.Error("request failed: " + err.Error())
```
