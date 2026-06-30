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

func TestValidatePath_WarnsWhenSocialActionSkillLacksApprovalBoundary(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	content := `---
name: Twitter Reply Skill
description: Publishes Twitter replies for support teams after reading a user complaint.
---
# Twitter Reply Skill

## Overview
This skill writes and publishes replies for social media support queues.

## Usage
Use this skill when support teams need to reply to customer tweets.

## Examples
- Publish a reply to a tweet that asks for order help.
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

	found := false
	for _, warning := range result.Warnings {
		if strings.Contains(warning, "approval or confirmation boundaries") {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected approval boundary warning, got %v", result.Warnings)
	}
}

func TestValidatePath_AcceptsSocialActionSkillWithApprovalBoundary(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	content := `---
name: Social Draft Skill
description: Drafts Twitter replies for support teams and waits for explicit approval before publishing.
---
# Social Draft Skill

## Overview
This skill prepares social media replies while keeping publication user-approved.

## Usage
Use this skill to draft replies, then ask the user to confirm before posting.

## Examples
- Draft a reply and wait for approval before publishing.
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
	for _, warning := range result.Warnings {
		if strings.Contains(warning, "approval or confirmation boundaries") {
			t.Fatalf("did not expect approval boundary warning, got %v", result.Warnings)
		}
	}
}

func TestValidatePath_WarnsOnNegatedSocialApprovalBoundary(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	content := `---
name: Social Publish Skill
description: Publishes Twitter posts without approval from the user.
---
# Social Publish Skill

## Overview
This skill handles social media publishing without human approval.

## Examples
- Publish queued posts without approval.
`
	if err := os.WriteFile(filepath.Join(dir, "SKILL.md"), []byte(content), 0o644); err != nil {
		t.Fatalf("write skill file: %v", err)
	}

	result, err := ValidatePath(dir, false)
	if err != nil {
		t.Fatalf("validate: %v", err)
	}

	found := false
	for _, warning := range result.Warnings {
		if strings.Contains(warning, "approval or confirmation boundaries") {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected approval boundary warning, got %v", result.Warnings)
	}
}

func TestValidatePath_WarnsOnInflectedSocialActions(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	content := `---
name: Social Operations Skill
description: Handles Twitter publishing, scheduling, and deleting posts.
---
# Social Operations Skill

## Overview
This skill coordinates social media publishing, scheduling, and deleting posts.

## Examples
- Prepare publishing and deleting operations.
`
	if err := os.WriteFile(filepath.Join(dir, "SKILL.md"), []byte(content), 0o644); err != nil {
		t.Fatalf("write skill file: %v", err)
	}

	result, err := ValidatePath(dir, false)
	if err != nil {
		t.Fatalf("validate: %v", err)
	}

	found := false
	for _, warning := range result.Warnings {
		if strings.Contains(warning, "approval or confirmation boundaries") {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected approval boundary warning, got %v", result.Warnings)
	}
}
