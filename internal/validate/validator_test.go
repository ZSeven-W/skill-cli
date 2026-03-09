package validate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestValidatePath_ValidSkillWithDirectories(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	for _, name := range []string{"scripts", "references", "assets"} {
		if err := os.Mkdir(filepath.Join(dir, name), 0o755); err != nil {
			t.Fatalf("mkdir %s: %v", name, err)
		}
	}

	content := `---
name: Good Skill
description: Detects duplicate entries and normalizes response formatting for agents.
version: 1.2.3
tags:
  - validation
  - quality
metadata:
  owner: platform
---
# Good Skill

## Overview
This skill ensures generated outputs meet quality expectations.

## Usage
Use this skill before publishing responses.

## Examples
- Run scripts/check.sh and review references/sample.md.
`

	if err := os.WriteFile(filepath.Join(dir, "SKILL.md"), []byte(content), 0o644); err != nil {
		t.Fatalf("write skill file: %v", err)
	}

	result, err := ValidatePath(dir, false)
	if err != nil {
		t.Fatalf("validate: %v", err)
	}
	if len(result.Errors) != 0 {
		t.Fatalf("expected no errors, got %v", result.Errors)
	}
	if len(result.Warnings) != 0 {
		t.Fatalf("expected no warnings, got %v", result.Warnings)
	}
	if result.Failed(false) {
		t.Fatalf("expected non-failed result")
	}
}

func TestValidatePath_UnknownFieldAndMissingStructure(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	content := `---
name: Skill
description: Useful
unknown: nope
---
# Skill
`
	if err := os.WriteFile(filepath.Join(dir, "SKILL.md"), []byte(content), 0o644); err != nil {
		t.Fatalf("write skill file: %v", err)
	}

	result, err := ValidatePath(dir, false)
	if err != nil {
		t.Fatalf("validate: %v", err)
	}
	if len(result.Errors) == 0 {
		t.Fatalf("expected schema error")
	}
	if !strings.Contains(result.Errors[0], "field unknown not found") {
		t.Fatalf("expected unknown field error, got %q", result.Errors[0])
	}
}

func TestValidatePath_StrictEscalatesBestPracticeWarnings(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	content := `---
name: Tiny
description: helper for things
---
# Tiny

## Overview
Quick notes.

## Usage
Do the task.
`
	if err := os.WriteFile(filepath.Join(dir, "SKILL.md"), []byte(content), 0o644); err != nil {
		t.Fatalf("write skill file: %v", err)
	}

	nonStrict, err := ValidatePath(dir, false)
	if err != nil {
		t.Fatalf("non-strict validate: %v", err)
	}
	if len(nonStrict.Errors) != 0 {
		t.Fatalf("expected non-strict to avoid best-practice errors, got %v", nonStrict.Errors)
	}
	if len(nonStrict.Warnings) == 0 {
		t.Fatalf("expected warnings in non-strict mode")
	}
	if nonStrict.Failed(false) {
		t.Fatalf("non-strict should not fail on warnings")
	}

	strictResult, err := ValidatePath(dir, true)
	if err != nil {
		t.Fatalf("strict validate: %v", err)
	}
	if len(strictResult.Errors) == 0 {
		t.Fatalf("expected strict mode to convert best-practice issues into errors")
	}
	if !strictResult.Failed(true) {
		t.Fatalf("strict mode should fail")
	}
}
