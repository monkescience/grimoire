---
type: rule
name: naming
description: Guidelines for naming variables, functions, and types
tags: [naming, readability, style]
---

# Naming Conventions

## Principles

1. **Be descriptive** - Names should reveal intent. Prefer `userCount` over `n` or `uc`.

2. **Use consistent casing** - Follow language conventions:
   - Go: `camelCase` for private, `PascalCase` for exported
   - TypeScript/JavaScript: `camelCase` for variables/functions, `PascalCase` for classes/types
   - Python: `snake_case` for variables/functions, `PascalCase` for classes

3. **Avoid abbreviations** - Use full words unless the abbreviation is universally understood (e.g., `URL`, `HTTP`, `ID`).

4. **Scope-appropriate length** - Short names for small scopes, longer names for larger scopes.

5. **Boolean naming** - Use `is`, `has`, `can`, `should` prefixes for boolean variables.

## Anti-patterns

- Single-letter variables (except for loop indices)
- Hungarian notation (`strName`, `iCount`)
- Meaningless names (`data`, `info`, `temp`, `foo`)
