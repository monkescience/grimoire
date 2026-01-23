---
type: skill
description: |
  Write and update tests using the given-when-then pattern for clear test structure.
  Use when writing unit tests, integration tests, or updating existing test suites.
  Covers test organization, assertions, mocking, and test-driven development (TDD)
  practices. Helps create readable, maintainable tests with clear setup, action,
  and verification phases.
---

# Testing

## Structure

Use the given-when-then pattern:

- **given**: setup and preconditions
- **when**: action being tested
- **then**: expected outcome

Use comments to mark each section when helpful.

## Example

```go
func TestUserService_Create(t *testing.T) {
    // given
    db := setupTestDB(t)
    svc := NewUserService(db)
    
    // when
    user, err := svc.Create("alice@example.com")
    
    // then
    require.NoError(t, err)
    assert.Equal(t, "alice@example.com", user.Email)
}
```
