package validate

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/fini/skill-cli/internal/formats"
	"gopkg.in/yaml.v3"
)

type Result struct {
	Path     string   `json:"path"`
	Errors   []string `json:"errors"`
	Warnings []string `json:"warnings"`
}

type frontmatterSchema struct {
	Name        string                 `yaml:"name"`
	Description string                 `yaml:"description"`
	Version     string                 `yaml:"version,omitempty"`
	Metadata    map[string]interface{} `yaml:"metadata,omitempty"`
	Tags        []string               `yaml:"tags,omitempty"`
}

func (r Result) Failed(strict bool) bool {
	if len(r.Errors) > 0 {
		return true
	}
	return strict && len(r.Warnings) > 0
}

func ValidatePath(path string, strict bool) (Result, error) {
	target, err := resolveSkillFile(path)
	if err != nil {
		return Result{}, err
	}

	content, err := os.ReadFile(target)
	if err != nil {
		return Result{}, fmt.Errorf("read skill file: %w", err)
	}

	frontmatter, body, err := splitSkillMarkdown(content)
	if err != nil {
		return Result{Path: target, Errors: []string{err.Error()}}, nil
	}

	schema, err := parseFrontmatterSchema(frontmatter)
	if err != nil {
		return Result{Path: target, Errors: []string{err.Error()}}, nil
	}

	result := Result{Path: target}

	validateSchema(schema, &result)
	validateStructure(schema, body, &result)
	validateBestPractices(schema, body, target, strict, &result)

	return result, nil
}

func splitSkillMarkdown(content []byte) (string, string, error) {
	text := string(content)
	if !strings.HasPrefix(text, "---\n") {
		return "", "", formats.ErrMissingFrontmatter
	}

	rest := text[4:]
	idx := strings.Index(rest, "\n---\n")
	if idx == -1 {
		return "", "", formats.ErrMissingFrontmatter
	}

	return rest[:idx], rest[idx+5:], nil
}

func parseFrontmatterSchema(frontmatter string) (frontmatterSchema, error) {
	dec := yaml.NewDecoder(bytes.NewBufferString(frontmatter))
	dec.KnownFields(true)

	var schema frontmatterSchema
	if err := dec.Decode(&schema); err != nil {
		return frontmatterSchema{}, fmt.Errorf("invalid YAML frontmatter: %w", err)
	}

	return schema, nil
}

func validateSchema(meta frontmatterSchema, result *Result) {
	if strings.TrimSpace(meta.Name) == "" {
		result.Errors = append(result.Errors, "metadata.name is required")
	}
	if strings.TrimSpace(meta.Description) == "" {
		result.Errors = append(result.Errors, "metadata.description is required")
	}

	if meta.Version != "" {
		versionPattern := regexp.MustCompile(`^v?\d+\.\d+\.\d+(?:-[0-9A-Za-z.-]+)?(?:\+[0-9A-Za-z.-]+)?$`)
		if !versionPattern.MatchString(meta.Version) {
			result.Errors = append(result.Errors, "metadata.version must be a semver-like string (e.g. 1.2.3)")
		}
	}

	for i, tag := range meta.Tags {
		if strings.TrimSpace(tag) == "" {
			result.Errors = append(result.Errors, fmt.Sprintf("metadata.tags[%d] must not be empty", i))
		}
	}
}

func validateStructure(meta frontmatterSchema, body string, result *Result) {
	trimmed := strings.TrimSpace(body)
	if trimmed == "" {
		result.Errors = append(result.Errors, "skill body content is required")
		return
	}

	if !regexp.MustCompile(`(?m)^#\s+`).MatchString(body) {
		result.Errors = append(result.Errors, "SKILL.md must include a top-level heading (e.g. '# Skill Name')")
	}
	if !regexp.MustCompile(`(?im)^##\s+overview\b`).MatchString(body) {
		result.Errors = append(result.Errors, "SKILL.md should include a '## Overview' section")
	}
	if !regexp.MustCompile(`(?im)^##\s+usage\b`).MatchString(body) {
		result.Errors = append(result.Errors, "SKILL.md should include a '## Usage' section")
	}

	heading := regexp.MustCompile(`(?m)^#\s+(.+?)\s*$`).FindStringSubmatch(body)
	if len(heading) > 1 {
		title := strings.TrimSpace(heading[1])
		if title != "" && strings.TrimSpace(meta.Name) != "" && !strings.EqualFold(title, strings.TrimSpace(meta.Name)) {
			result.Warnings = append(result.Warnings, "top-level heading does not match metadata.name")
		}
	}
}

func validateBestPractices(meta frontmatterSchema, body, skillFile string, strict bool, result *Result) {
	desc := strings.TrimSpace(meta.Description)
	if desc != "" {
		minLen := 20
		maxLen := 180
		if strict {
			minLen = 30
			maxLen = 140
		}
		if len(desc) < minLen {
			addBestPracticeIssue(result, strict, fmt.Sprintf("description is too short (%d chars); aim for at least %d", len(desc), minLen))
		}
		if len(desc) > maxLen {
			addBestPracticeIssue(result, strict, fmt.Sprintf("description is too long (%d chars); keep it under %d", len(desc), maxLen))
		}

		vague := regexp.MustCompile(`(?i)\b(various|misc|general purpose|things|stuff|helper)\b`)
		if vague.MatchString(desc) {
			addBestPracticeIssue(result, strict, "description appears vague; be specific about scope and trigger conditions")
		}
	}

	hasExamples := regexp.MustCompile(`(?i)\bexample(s)?\b`).MatchString(body) || strings.Contains(body, "```")
	if !hasExamples {
		addBestPracticeIssue(result, strict, "no examples found; add an example section or fenced code blocks")
	}

	mentioned := map[string]bool{}
	for _, dir := range []string{"scripts", "references", "assets"} {
		pattern := regexp.MustCompile(fmt.Sprintf(`(?i)\b%s/`, regexp.QuoteMeta(dir)))
		if pattern.MatchString(body) {
			mentioned[dir] = true
		}
	}

	skillDir := filepath.Dir(skillFile)
	for dir := range mentioned {
		full := filepath.Join(skillDir, dir)
		if info, err := os.Stat(full); err != nil || !info.IsDir() {
			addBestPracticeIssue(result, strict, fmt.Sprintf("body references %s/ but %s does not exist", dir, full))
		}
	}
}

func addBestPracticeIssue(result *Result, strict bool, message string) {
	if strict {
		result.Errors = append(result.Errors, "strict: "+message)
		return
	}
	result.Warnings = append(result.Warnings, message)
}

func resolveSkillFile(path string) (string, error) {
	info, err := os.Stat(path)
	if err != nil {
		return "", fmt.Errorf("stat path: %w", err)
	}

	if info.IsDir() {
		candidate := filepath.Join(path, "SKILL.md")
		if _, err := os.Stat(candidate); err != nil {
			return "", fmt.Errorf("SKILL.md not found in %s", path)
		}
		return candidate, nil
	}
	return path, nil
}
