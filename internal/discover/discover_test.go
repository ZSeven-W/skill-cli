package discover

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDiscoverInstalledSkillsFromEnvRoots(t *testing.T) {
	home := t.TempDir()
	claudeRoot := filepath.Join(home, "claude-skills")
	openclawRoot := filepath.Join(home, "openclaw-skills")
	mustMkdirAll(t, claudeRoot)
	mustMkdirAll(t, openclawRoot)

	mustCreateSkillDir(t, claudeRoot, "beta")
	mustCreateSkillDir(t, claudeRoot, "alpha")
	mustCreateSkillDir(t, openclawRoot, "omega")
	mustMkdirAll(t, filepath.Join(claudeRoot, "not-a-skill"))

	t.Setenv("HOME", home)
	t.Setenv("CLAUDE_SKILLS_DIR", claudeRoot)
	t.Setenv("OPENCLAW_SKILLS_DIR", openclawRoot)

	got, err := DiscoverInstalledSkills()
	if err != nil {
		t.Fatalf("discover: %v", err)
	}

	if len(got) != 3 {
		t.Fatalf("expected 3 skills, got %d", len(got))
	}

	want := []InstalledSkill{
		{Name: "alpha", Path: filepath.Join(claudeRoot, "alpha"), Platform: "claude"},
		{Name: "beta", Path: filepath.Join(claudeRoot, "beta"), Platform: "claude"},
		{Name: "omega", Path: filepath.Join(openclawRoot, "omega"), Platform: "openclaw"},
	}

	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("skill[%d] mismatch: got %+v want %+v", i, got[i], want[i])
		}
	}
}

func TestFindOpenClawSkillDirs(t *testing.T) {
	t.Parallel()

	base := t.TempDir()
	mustMkdirAll(t, filepath.Join(base, "versions", "node", "lib", "openclaw", "skills"))
	mustMkdirAll(t, filepath.Join(base, "other", "openclaw", "skills"))
	mustMkdirAll(t, filepath.Join(base, "openclaw", "not-skills"))

	got, err := findOpenClawSkillDirs(base)
	if err != nil {
		t.Fatalf("find openclaw dirs: %v", err)
	}

	if len(got) != 2 {
		t.Fatalf("expected 2 directories, got %d", len(got))
	}
	for _, root := range got {
		if filepath.Base(root) != "skills" || filepath.Base(filepath.Dir(root)) != "openclaw" {
			t.Fatalf("unexpected root: %s", root)
		}
	}
}

func TestListCommandJSONAndPlatformFilter(t *testing.T) {
	home := t.TempDir()
	claudeRoot := filepath.Join(home, "claude")
	openclawRoot := filepath.Join(home, "openclaw")
	mustMkdirAll(t, claudeRoot)
	mustMkdirAll(t, openclawRoot)
	mustCreateSkillDir(t, claudeRoot, "alpha")
	mustCreateSkillDir(t, openclawRoot, "beta")

	t.Setenv("HOME", home)
	t.Setenv("CLAUDE_SKILLS_DIR", claudeRoot)
	t.Setenv("OPENCLAW_SKILLS_DIR", openclawRoot)

	cmd := NewListCommand()
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&out)
	cmd.SetArgs([]string{"--format", "json", "--platform", "claude"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("execute: %v", err)
	}

	var got []InstalledSkill
	if err := json.Unmarshal(out.Bytes(), &got); err != nil {
		t.Fatalf("decode json output: %v; output=%q", err, out.String())
	}

	if len(got) != 1 {
		t.Fatalf("expected 1 skill, got %d", len(got))
	}
	if got[0].Platform != "claude" || got[0].Name != "alpha" {
		t.Fatalf("unexpected skill: %+v", got[0])
	}
}

func TestListCommandRejectsUnsupportedFormat(t *testing.T) {
	t.Parallel()

	cmd := NewListCommand()
	cmd.SilenceUsage = true
	cmd.SetArgs([]string{"--format", "yaml"})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error for unsupported format")
	}
	if !strings.Contains(err.Error(), "unsupported format") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func mustCreateSkillDir(t *testing.T, root, name string) {
	t.Helper()
	skillDir := filepath.Join(root, name)
	mustMkdirAll(t, skillDir)
	if err := os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte("# Skill\n"), 0o644); err != nil {
		t.Fatalf("write SKILL.md: %v", err)
	}
}

func mustMkdirAll(t *testing.T, path string) {
	t.Helper()
	if err := os.MkdirAll(path, 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", path, err)
	}
}
