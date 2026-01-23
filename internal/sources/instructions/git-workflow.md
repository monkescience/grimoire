---
type: instruction
description: Git commits and pull request conventions
order: 30
---

## Git & Pull Requests

- Use Conventional Commits: `type(scope): description`
- Keep commits focused and atomic (one logical change per commit)
- Subject line only, no body
- Do not add Co-Authored-By lines

### Commit Types

Trigger release: `feat`, `fix`, `perf`
No release: `docs`, `style`, `refactor`, `test`, `chore`, `build`, `ci`

### Breaking Changes

Add `!` after type (or scope) to indicate a breaking change that affects external consumers.

For services: API contract changes (endpoints, request/response formats, auth flows)
For libraries: Public API changes (exported functions, types, interfaces)

When uncertain if a change is breaking, ask before committing.

### History Rewriting

- Do not use `git commit --amend` on commits that have been pushed
- If a pushed commit needs fixing, create a new commit

### Pull Request Reviews

- Use Conventional Commits format for PR title
- Keep description concise
- Use Conventional Comments: `label: subject`
  - `praise:` - Highlight something done well
  - `nitpick:` - Minor style/preference, non-blocking
  - `suggestion:` - Propose an alternative approach
  - `issue:` - Something that must be addressed
  - `question:` - Seeking clarification or understanding
  - `thought:` - Share an idea without requiring action
