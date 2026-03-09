package formats

import (
	"errors"
	"strings"
	"testing"
)

func TestParseSkillMarkdown_Success(t *testing.T) {
	content := []byte(`---
name: Test Skill
description: test description
version: 1.0.0
---
# Overview

Hello world.
`)

	skill, err := ParseSkillMarkdown(content)
	if err != nil {
		t.Fatalf("ParseSkillMarkdown returned error: %v", err)
	}

	if skill.Metadata.Name != "Test Skill" {
		t.Fatalf("unexpected name: %q", skill.Metadata.Name)
	}
	if skill.Metadata.Description != "test description" {
		t.Fatalf("unexpected description: %q", skill.Metadata.Description)
	}
	if skill.Metadata.Version != "1.0.0" {
		t.Fatalf("unexpected version: %q", skill.Metadata.Version)
	}
	if !strings.Contains(skill.Body, "# Overview") {
		t.Fatalf("expected body to contain heading, got: %q", skill.Body)
	}
}

func TestParseSkillMarkdown_MissingFrontmatter(t *testing.T) {
	_, err := ParseSkillMarkdown([]byte("# Overview\nNo frontmatter"))
	if !errors.Is(err, ErrMissingFrontmatter) {
		t.Fatalf("expected ErrMissingFrontmatter, got: %v", err)
	}
}

func TestParseSkillMarkdown_InvalidYAML(t *testing.T) {
	content := []byte(`---
name: [not-closed
---
body
`)

	_, err := ParseSkillMarkdown(content)
	if err == nil {
		t.Fatal("expected parse error for invalid YAML")
	}
	if !strings.Contains(err.Error(), "invalid YAML frontmatter") {
		t.Fatalf("expected invalid YAML error wrapper, got: %v", err)
	}
}

func TestRenderSkillMarkdown(t *testing.T) {
	skill := Skill{
		Metadata: Metadata{
			Name:        "Rendered Skill",
			Description: "render description",
			Version:     "2.0.0",
		},
		Body: "\n# Usage\n\nDo the thing.\n",
	}

	out, err := RenderSkillMarkdown(skill)
	if err != nil {
		t.Fatalf("RenderSkillMarkdown returned error: %v", err)
	}

	result := string(out)
	if !strings.HasPrefix(result, "---\n") {
		t.Fatalf("expected YAML frontmatter start, got: %q", result)
	}
	if !strings.Contains(result, "name: Rendered Skill\n") {
		t.Fatalf("expected name in frontmatter, got: %q", result)
	}
	if !strings.Contains(result, "description: render description\n") {
		t.Fatalf("expected description in frontmatter, got: %q", result)
	}
	if !strings.Contains(result, "# Usage\n\nDo the thing.\n") {
		t.Fatalf("expected normalized body in output, got: %q", result)
	}
	if !strings.HasSuffix(result, "\n") {
		t.Fatalf("expected trailing newline in output, got: %q", result)
	}
}
