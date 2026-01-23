---
type: rule
description: Use context-aware slog functions (DebugContext, InfoContext, etc.) when context is available
tags: [go, logging, slog, context]
globs: ["*.go"]
---

## Good

```go
func ProcessOrder(ctx context.Context, orderID string) error {
    slog.InfoContext(ctx, "processing order", slog.String("order_id", orderID))
    // ...
    slog.ErrorContext(ctx, "processing failed",
        slog.String("order_id", orderID),
        slog.Any("error", err))
    return err
}
```

## Bad

```go
func ProcessOrder(ctx context.Context, orderID string) error {
    slog.Info("processing order", slog.String("order_id", orderID))  // ctx available but not used
    // ...
    slog.Error("processing failed",
        slog.String("order_id", orderID),
        slog.Any("error", err))
    return err
}
```
