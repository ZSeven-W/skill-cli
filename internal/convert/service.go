package convert

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fini/skill-cli/internal/formats"
)

type Options struct {
	From   string
	To     string
	Input  string
	Output string
	Force  bool
}

func Convert(opts Options) error {
	if opts.From != "openclaw" && opts.From != "claude" {
		return fmt.Errorf("unsupported --from format %q", opts.From)
	}
	if opts.To != "openclaw" && opts.To != "claude" {
		return fmt.Errorf("unsupported --to format %q", opts.To)
	}
	if opts.From == opts.To {
		return fmt.Errorf("--from and --to must be different")
	}

	skillFile, err := resolveSkillFile(opts.Input)
	if err != nil {
		return err
	}

	content, err := os.ReadFile(skillFile)
	if err != nil {
		return fmt.Errorf("read source SKILL.md: %w", err)
	}

	skill, err := formats.ParseSkillMarkdown(content)
	if err != nil {
		return fmt.Errorf("parse source skill: %w", err)
	}

	if err := prepareOutput(opts.Output, opts.Force); err != nil {
		return err
	}
	if opts.To == "claude" {
		for _, dir := range []string{"scripts", "references", "assets"} {
			if err := os.MkdirAll(filepath.Join(opts.Output, dir), 0o755); err != nil {
				return fmt.Errorf("create %s directory: %w", dir, err)
			}
		}
	}

	rendered, err := formats.RenderSkillMarkdown(skill)
	if err != nil {
		return err
	}

	if err := os.WriteFile(filepath.Join(opts.Output, "SKILL.md"), rendered, 0o644); err != nil {
		return fmt.Errorf("write output SKILL.md: %w", err)
	}

	fmt.Fprintf(os.Stdout, "Converted %s -> %s at %s\n", opts.From, opts.To, opts.Output)
	return nil
}

func prepareOutput(path string, force bool) error {
	if info, err := os.Stat(path); err == nil {
		if !info.IsDir() {
			return fmt.Errorf("output path %s exists and is not a directory", path)
		}
		if !force {
			return fmt.Errorf("output path %s already exists (use --force to overwrite)", path)
		}
		if err := os.RemoveAll(path); err != nil {
			return fmt.Errorf("clear output path: %w", err)
		}
	}
	if err := os.MkdirAll(path, 0o755); err != nil {
		return fmt.Errorf("create output directory: %w", err)
	}
	return nil
}

func resolveSkillFile(path string) (string, error) {
	info, err := os.Stat(path)
	if err != nil {
		return "", fmt.Errorf("stat input path: %w", err)
	}
	if info.IsDir() {
		file := filepath.Join(path, "SKILL.md")
		if _, err := os.Stat(file); err != nil {
			return "", fmt.Errorf("input directory does not contain SKILL.md")
		}
		return file, nil
	}
	return path, nil
}
