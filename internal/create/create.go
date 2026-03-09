package create

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fini/skill-cli/internal/formats"
	"github.com/spf13/cobra"
)

type options struct {
	Name        string
	Description string
	Version     string
	Format      string
	OutputDir   string
	Force       bool
}

func NewCommand() *cobra.Command {
	opts := &options{}

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new skill from template",
		RunE: func(cmd *cobra.Command, args []string) error {
			if opts.Name == "" || opts.Description == "" {
				return fmt.Errorf("both --name and --description are required")
			}
			if opts.Format != "openclaw" && opts.Format != "claude" {
				return fmt.Errorf("unsupported format %q (use openclaw or claude)", opts.Format)
			}
			return createSkill(*opts)
		},
	}

	cmd.Flags().StringVar(&opts.Name, "name", "", "skill name")
	cmd.Flags().StringVar(&opts.Description, "description", "", "skill description")
	cmd.Flags().StringVar(&opts.Version, "version", "0.1.0", "skill version")
	cmd.Flags().StringVar(&opts.Format, "format", "claude", "skill format: claude|openclaw")
	cmd.Flags().StringVar(&opts.OutputDir, "output", ".", "directory where skill folder is created")
	cmd.Flags().BoolVar(&opts.Force, "force", false, "overwrite existing skill directory")

	return cmd
}

func createSkill(opts options) error {
	slug := sanitizeName(opts.Name)
	skillDir := filepath.Join(opts.OutputDir, slug)

	if _, err := os.Stat(skillDir); err == nil {
		if !opts.Force {
			return fmt.Errorf("directory %s already exists (use --force to overwrite)", skillDir)
		}
		if err := os.RemoveAll(skillDir); err != nil {
			return fmt.Errorf("remove existing directory: %w", err)
		}
	}

	if err := os.MkdirAll(skillDir, 0o755); err != nil {
		return fmt.Errorf("create skill directory: %w", err)
	}

	if opts.Format == "claude" {
		for _, dir := range []string{"scripts", "references", "assets"} {
			if err := os.MkdirAll(filepath.Join(skillDir, dir), 0o755); err != nil {
				return fmt.Errorf("create %s directory: %w", dir, err)
			}
		}
	}

	skill := formats.Skill{
		Metadata: formats.Metadata{
			Name:        opts.Name,
			Description: opts.Description,
			Version:     opts.Version,
		},
		Body: defaultBody(opts.Name),
	}

	rendered, err := formats.RenderSkillMarkdown(skill)
	if err != nil {
		return err
	}

	skillFile := filepath.Join(skillDir, "SKILL.md")
	if err := os.WriteFile(skillFile, rendered, 0o644); err != nil {
		return fmt.Errorf("write SKILL.md: %w", err)
	}

	fmt.Fprintf(os.Stdout, "Created %s skill at %s\n", opts.Format, skillDir)
	return nil
}

func sanitizeName(name string) string {
	n := strings.TrimSpace(strings.ToLower(name))
	n = strings.ReplaceAll(n, " ", "-")
	n = strings.ReplaceAll(n, "_", "-")
	return n
}

func defaultBody(name string) string {
	return fmt.Sprintf(`# %s

## Overview
Describe what this skill does and when the agent should use it.

## Usage
- Add clear usage guidance here.
- Add examples in references/ as needed.
`, name)
}
