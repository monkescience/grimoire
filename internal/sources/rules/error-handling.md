---
type: rule
description: Handle errors at boundaries, propagate with context, prefer returns over panics
globs: ["*.go", "*.ts", "*.js", "*.py", "*.rs", "*.java", "*.rb", "*.c", "*.cpp", "*.h", "*.cs"]
---

# Error Handling

- Handle errors at appropriate boundaries (API, user input, external calls)
- Propagate errors with context rather than swallowing them
- Prefer returning errors over panicking/throwing
- Fail fast on programmer errors; handle gracefully on user/external errors
- Do not add defensive checks for conditions the type system already prevents
