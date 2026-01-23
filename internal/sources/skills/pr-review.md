---
type: skill
trigger: When reviewing a pull request
tags: [git, pull-requests, review, github]
arguments:
  - name: pr
    description: PR number or URL to review
    required: false
  - name: focus
    description: Specific aspects to focus on (e.g., security, performance)
    required: false
---

# PR Review

Review a pull request systematically and provide actionable feedback.

{{pr}}

## Process

1. **Understand Context**
   - Read the PR title and description
   - Understand the purpose and scope of changes
   - Check linked issues or tickets

2. **Review Changes**
   - Use `gh pr diff` or review the diff
   - Consider the overall architecture impact
   - Check for breaking changes

3. **Evaluate Quality**

{{focus}}

### Code Quality
- Is the code readable and well-organized?
- Are names meaningful and consistent?
- Is complexity appropriate?

### Correctness
- Does the logic handle edge cases?
- Are error conditions handled?
- Are assumptions documented?

### Testing
- Are tests included for new functionality?
- Do tests cover edge cases?
- Are existing tests still passing?

### Security
- Any hardcoded secrets or credentials?
- Input validation present?
- SQL injection, XSS, or other vulnerabilities?

### Performance
- Any obvious performance issues?
- Unnecessary database queries?
- Memory leaks or resource management issues?

## Feedback Format

Use Conventional Comments for feedback:

- `praise:` Highlight something done well
- `nitpick:` Minor style/preference (non-blocking)
- `suggestion:` Propose an alternative approach
- `issue:` Must be addressed before merge
- `question:` Seeking clarification
- `thought:` Share an idea without requiring action

### Example Comments

```
praise: Clean separation of concerns here. The service layer is well-defined.

issue: This SQL query is vulnerable to injection. Use parameterized queries.

suggestion: Consider using a Map here instead of repeated array lookups.
The O(1) lookup would improve performance for larger datasets.

nitpick: Prefer `const` over `let` for variables that aren't reassigned.
```

## Summary

End with an overall assessment:
- **Approve**: Ready to merge
- **Request Changes**: Issues must be addressed
- **Comment**: Feedback provided, no blocking issues
