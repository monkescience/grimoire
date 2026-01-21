package grimoire

import "testing"

func TestDeriveName(t *testing.T) {
	tests := []struct {
		path string
		typ  Type
		want string
	}{
		// Rules with subfolders
		{"rules/go/error-assignment.md", TypeRule, "go/error-assignment"},
		{"rules/java/null-handling.md", TypeRule, "java/null-handling"},
		{"rules/go/errors/assignment.md", TypeRule, "go/errors/assignment"},

		// Rules without subfolders
		{"rules/naming.md", TypeRule, "naming"},

		// Skills with subfolders
		{"skills/debugging/advanced.md", TypeSkill, "debugging/advanced"},

		// Skills without subfolders
		{"skills/refactor.md", TypeSkill, "refactor"},

		// Flat structure (no type prefix)
		{"my-rule.md", TypeRule, "my-rule"},
		{"some/path/my-rule.md", TypeRule, "my-rule"},

		// Alternative singular prefix
		{"rule/go/error.md", TypeRule, "go/error"},
		{"skill/refactor.md", TypeSkill, "refactor"},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			got := deriveName(tt.path, tt.typ)
			if got != tt.want {
				t.Errorf("deriveName(%q, %q) = %q, want %q", tt.path, tt.typ, got, tt.want)
			}
		})
	}
}
