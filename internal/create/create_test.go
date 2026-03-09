package create

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSanitizeNameRejectsNonAlnum(t *testing.T) {
	if got := sanitizeName("!!!"); got != "" {
		t.Fatalf("expected empty slug, got %q", got)
	}
}

func TestCreateSkillRejectsEmptySanitizedSlug(t *testing.T) {
	tempDir := t.TempDir()

	err := createSkill(options{
		Name:        "!!!",
		Description: "desc",
		Version:     "0.1.0",
		Format:      "claude",
		OutputDir:   tempDir,
	})
	if err == nil {
		t.Fatal("expected error for empty sanitized slug")
	}
}

func TestValidateSafeDeletePathBlocksDangerousTargets(t *testing.T) {
	cases := []string{".", "/", "~", "/home"}
	for _, tc := range cases {
		tc := tc
		t.Run(tc, func(t *testing.T) {
			if err := validateSafeDeletePath(tc); err == nil {
				t.Fatalf("expected %q to be rejected", tc)
			}
		})
	}
}

func TestValidateSafeDeletePathAllowsRegularDirectory(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "safe-skill-dir")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	if err := validateSafeDeletePath(dir); err != nil {
		t.Fatalf("expected safe path to pass validation, got %v", err)
	}
}
