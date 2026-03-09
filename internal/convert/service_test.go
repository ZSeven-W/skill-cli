package convert

import (
	"os"
	"path/filepath"
	"testing"
)

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
	dir := filepath.Join(t.TempDir(), "safe-output-dir")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	if err := validateSafeDeletePath(dir); err != nil {
		t.Fatalf("expected safe path to pass validation, got %v", err)
	}
}

func TestPrepareOutputRejectsForceOnDangerousPath(t *testing.T) {
	if err := prepareOutput(".", true); err == nil {
		t.Fatal("expected dangerous path rejection")
	}
}
