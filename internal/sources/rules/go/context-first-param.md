---
type: rule
description: Context should be the first parameter and named ctx
globs: ["*.go"]
---

## Good

```go
func ProcessOrder(ctx context.Context, orderID string) error
func (s *Service) FetchUser(ctx context.Context, id int64) (*User, error)
```

## Bad

```go
func ProcessOrder(orderID string, ctx context.Context) error
func (s *Service) FetchUser(id int64, c context.Context) (*User, error)
```
