# Sources

Grimoire provides guidance through three types of sources: **instructions**, **skills**, and **rules**.

## Overview

| Type | Purpose | When Loaded |
|------|---------|-------------|
| **Instruction** | Base guidance injected into AI context | Always (server startup) |
| **Skill** | How to perform specific tasks | On-demand via `guidance()` |
| **Rule** | Project conventions and standards | On-demand via `guidance()` |

## Instructions

Instructions define general behavior and are always injected into the AI's context.

### Frontmatter

```yaml
---
type: instruction
description: <short phrase describing what this covers>
order: <number>
---
```

| Field | Required | Description |
|-------|----------|-------------|
| `type` | Yes | Must be `instruction` |
| `description` | Yes | Short phrase describing the instruction's purpose |
| `order` | Yes | Injection order (lower = earlier in context) |

### Body Format

- Start with `## <Topic Name>` heading
- Use bullet lists for actionable guidance
- Group related items under `###` subheadings

### Example

```markdown
---
type: instruction
description: Scope discipline and decision-making behavior
order: 10
---

## Behavior

- Ask before making significant changes
- Complete the requested task, nothing more
```

## Skills

Skills define HOW to perform specific tasks (debugging, code review, refactoring).
Following the [Agent Skills specification](https://agentskills.io), skills use a rich
description field that helps agents understand when to activate the skill.

### Frontmatter

```yaml
---
type: skill
description: |
  Detailed description of what this skill does and when to use it.
  Should be 1-4 sentences covering the skill's purpose, when to activate it,
  and what it helps accomplish. Keywords help with task matching.
---
```

| Field | Required | Description |
|-------|----------|-------------|
| `type` | Yes | Must be `skill` |
| `description` | Yes | Rich description (up to 1024 chars) explaining what the skill does AND when to use it |
| `arguments` | No | Parameters for templating with `{{argName}}` syntax |
| `agents` | No | Agent names this skill can delegate to |

### Description Guidelines

The description is the primary field for skill activation. Write it to:

1. **Explain what the skill does** - the capabilities it provides
2. **Describe when to use it** - task patterns that should trigger activation
3. **Include relevant keywords** - terms users might mention (e.g., "commit", "git", "version control")

### Body Format

- Start with `# <Skill Name>` heading (no "Skill" suffix)
- Use numbered steps or clear sections
- Include examples where helpful

### Example

```markdown
---
type: skill
description: |
  Write and update tests using the given-when-then pattern for clear test structure.
  Use when writing unit tests, integration tests, or updating existing test suites.
  Covers test organization, assertions, mocking, and test-driven development (TDD)
  practices.
---

# Testing

## Structure

Use the given-when-then pattern:

- **given**: setup and preconditions
- **when**: action being tested
- **then**: expected outcome
```

## Rules

Rules define project conventions and standards. They can be general or language-specific.

### Frontmatter

```yaml
---
type: rule
description: <concise statement of the rule>
globs: ["*.ext"]
---
```

| Field | Required | Description |
|-------|----------|-------------|
| `type` | Yes | Must be `rule` |
| `description` | Yes | Concise statement of what to do/avoid |
| `globs` | Recommended | File patterns this rule applies to |

### Body Format

Choose based on content type:

**Good/Bad format** - for code-centric rules with clear examples:

```markdown
## Good

\`\`\`go
// correct example
\`\`\`

## Bad

\`\`\`go
// incorrect example
\`\`\`
```

**Bullet list format** - for conceptual/behavioral rules:

```markdown
# Rule Topic

- Guideline 1
- Guideline 2
- Guideline 3
```

### Example

```markdown
---
type: rule
description: Context should be the first parameter and named ctx
globs: ["*.go"]
---

## Good

\`\`\`go
func ProcessOrder(ctx context.Context, orderID string) error
\`\`\`

## Bad

\`\`\`go
func ProcessOrder(orderID string, ctx context.Context) error
\`\`\`
```

## File Organization

```
internal/sources/
  instructions/     # Always-loaded base guidance
    behavior.md
    code-principles.md
    tooling.md
  skills/           # On-demand task guidance
    code-review.md
    debug.md
    git-workflow.md
    refactor.md
    testing.md
  rules/            # Project conventions
    error-handling.md
    todos.md
    go/             # Language-specific rules
      context-first-param.md
      defer-close.md
      ...
```

## How Rules are Applied

The AI sees all rules listed by name, globs, and description in the `guidance` tool. The server instructions tell the AI:

> "Apply rules based on their description. Load only if you need examples."

This means:

1. **The description IS the rule** - it should be self-contained and actionable
2. **Good/Bad examples are optional** - loaded only when the AI needs clarification
3. **Globs signal relevance** - help the AI know which file types a rule applies to

## Naming Conventions

| Type | Name Format | Examples |
|------|-------------|----------|
| Instructions | `kebab-case` | `code-principles`, `behavior` |
| Skills | `kebab-case` | `code-review`, `git-workflow` |
| Rules | `kebab-case` or `dir/kebab-case` | `todos`, `go/context-first-param` |
