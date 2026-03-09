package formats

import (
	"bytes"
	"errors"
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

var ErrMissingFrontmatter = errors.New("missing YAML frontmatter")

type Metadata struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Version     string `yaml:"version,omitempty"`
}

type Skill struct {
	Metadata Metadata
	Body     string
}

func ParseSkillMarkdown(content []byte) (Skill, error) {
	text := string(content)
	if !strings.HasPrefix(text, "---\n") {
		return Skill{}, ErrMissingFrontmatter
	}

	rest := text[4:]
	idx := strings.Index(rest, "\n---\n")
	if idx == -1 {
		return Skill{}, ErrMissingFrontmatter
	}

	frontmatter := rest[:idx]
	body := rest[idx+5:]

	var meta Metadata
	dec := yaml.NewDecoder(bytes.NewBufferString(frontmatter))
	dec.KnownFields(false)
	if err := dec.Decode(&meta); err != nil {
		return Skill{}, fmt.Errorf("invalid YAML frontmatter: %w", err)
	}

	return Skill{Metadata: meta, Body: body}, nil
}

func RenderSkillMarkdown(skill Skill) ([]byte, error) {
	data, err := yaml.Marshal(skill.Metadata)
	if err != nil {
		return nil, fmt.Errorf("marshal metadata: %w", err)
	}

	var b strings.Builder
	b.WriteString("---\n")
	b.Write(data)
	b.WriteString("---\n\n")
	b.WriteString(strings.TrimSpace(skill.Body))
	b.WriteString("\n")

	return []byte(b.String()), nil
}
