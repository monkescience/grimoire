---
type: skill
description: Create a well-structured git commit
tags: [git, commits, version-control]
arguments:
  - name: message
    description: Optional commit message to use or refine
    required: false
  - name: files
    description: Specific files to focus on for the commit
    required: false
---

# Commit

Create a git commit for the staged changes.

## Process

1. **Analyze Changes**
   - Run `git status` to see staged and unstaged files
   - Run `git diff --cached` to review staged changes
   - Understand the purpose and scope of the changes

2. **Craft Commit Message**
   - Use Conventional Commits format: `type(scope): description`
   - Subject line only, max 72 characters
   - Use imperative mood ("add" not "added")
   - Be specific about what changed and why

{{message}}

3. **Commit Types**
   - `feat`: New feature (triggers release)
   - `fix`: Bug fix (triggers release)
   - `perf`: Performance improvement (triggers release)
   - `docs`: Documentation only
   - `style`: Code style (formatting, semicolons)
   - `refactor`: Code change without feature/fix
   - `test`: Adding or fixing tests
   - `chore`: Maintenance tasks
   - `build`: Build system changes
   - `ci`: CI configuration changes

4. **Scope** (optional)
   - Use lowercase
   - Identify affected component/module
   - Keep it short and consistent

{{files}}

## Examples

Good:
- `feat(auth): add OAuth2 login support`
- `fix(api): handle null response in user endpoint`
- `refactor(utils): simplify date formatting logic`

Bad:
- `updated files` (vague)
- `Fix bug` (no scope, not descriptive)
- `feat: Added new feature for user authentication` (past tense, too long)

## Verification

After committing:
- Run `git log -1` to verify the commit
- Ensure the message accurately describes the change
