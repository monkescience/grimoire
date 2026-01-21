---
type: rule
description: Best practices for handling errors in code
tags: [errors, best-practices, reliability]
---

# Error Handling

## Principles

1. **Handle errors explicitly** - Never ignore errors. Either handle them, propagate them, or explicitly document why they can be ignored.

2. **Add context when propagating** - When returning an error up the call stack, wrap it with additional context about what operation failed.

3. **Fail fast on programmer errors** - Use assertions or panics for conditions that indicate bugs in the code.

4. **Graceful degradation for runtime errors** - Handle expected runtime errors (network failures, invalid input) gracefully with appropriate user feedback.

5. **Use typed errors when appropriate** - Create custom error types when callers need to distinguish between different error conditions.

## Examples

### Go
```go
// Add context when wrapping
if err := doSomething(); err != nil {
    return fmt.Errorf("failed to do something: %w", err)
}
```

### TypeScript
```typescript
// Use Result types for expected failures
type Result<T, E> = { ok: true; value: T } | { ok: false; error: E };
```
