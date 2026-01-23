---
type: rule
description: Use explicit, scoped TODO comments for deferred work
tags: [todos, comments, task-tracking]
---

# TODO Comments

- When discovering missing or incomplete implementations, insert TODO comments directly in the relevant code
- TODOs must be explicit and scoped: `// TODO: implement retry logic with exponential backoff`
- Never leave vague TODOs like `// TODO: fix this` or `// TODO: implement`
- When deferring work, explain why and what's needed
- Prefer embedding intent in code over external documentation or plans
- Remove TODOs only when the work is complete
