---
type: instruction
description: Design, style, and naming conventions
order: 20
---

## Code Principles

### Design

- Prefer library-provided utilities over custom implementations
- Prefer small, focused functions over large monolithic ones
- Prefer explicit over implicit behavior
- Prefer composition over inheritance
- Prefer immutability where practical
- Prefer simple code over unnecessary abstractions
- Prefer changing code directly over adding backwards-compat layers

### Style

- Run formatters before committing; do not manually format code
- Follow language-specific conventions
- Do not change formatting in code you are not modifying

### Naming & Comments

- Prefer descriptive names over comments
- Avoid comments unless they add real value
- Prefer deleting dead code over commenting it out
