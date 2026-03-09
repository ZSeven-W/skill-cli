package validate

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fini/skill-cli/internal/formats"
)

type Result struct {
	Path   string
	Errors []string
}

func ValidatePath(path string) (Result, error) {
	target, err := resolveSkillFile(path)
	if err != nil {
		return Result{}, err
	}

	content, err := os.ReadFile(target)
	if err != nil {
		return Result{}, fmt.Errorf("read skill file: %w", err)
	}

	skill, err := formats.ParseSkillMarkdown(content)
	if err != nil {
		return Result{Path: target, Errors: []string{err.Error()}}, nil
	}

	var errs []string
	if strings.TrimSpace(skill.Metadata.Name) == "" {
		errs = append(errs, "metadata.name is required")
	}
	if strings.TrimSpace(skill.Metadata.Description) == "" {
		errs = append(errs, "metadata.description is required")
	}
	if strings.TrimSpace(skill.Body) == "" {
		errs = append(errs, "skill body content is required")
	}

	return Result{Path: target, Errors: errs}, nil
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
