package discover

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type InstalledSkill struct {
	Name     string
	Path     string
	Platform string
}

func DiscoverInstalledSkills() ([]InstalledSkill, error) {
	roots, err := discoverRoots()
	if err != nil {
		return nil, err
	}

	var skills []InstalledSkill
	for platform, dirs := range roots {
		for _, dir := range dirs {
			entries, err := os.ReadDir(dir)
			if err != nil {
				continue
			}
			for _, entry := range entries {
				if !entry.IsDir() {
					continue
				}
				skillPath := filepath.Join(dir, entry.Name())
				if _, err := os.Stat(filepath.Join(skillPath, "SKILL.md")); err != nil {
					continue
				}
				skills = append(skills, InstalledSkill{
					Name:     entry.Name(),
					Path:     skillPath,
					Platform: platform,
				})
			}
		}
	}

	sort.Slice(skills, func(i, j int) bool {
		if skills[i].Platform == skills[j].Platform {
			return skills[i].Name < skills[j].Name
		}
		return skills[i].Platform < skills[j].Platform
	})
	return skills, nil
}

func discoverRoots() (map[string][]string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("resolve home directory: %w", err)
	}

	roots := map[string][]string{}

	claudeRoot := os.Getenv("CLAUDE_SKILLS_DIR")
	if claudeRoot == "" {
		claudeRoot = filepath.Join(home, ".claude", "skills")
	}
	if dirExists(claudeRoot) {
		roots["claude"] = []string{claudeRoot}
	}

	openclawRoot := os.Getenv("OPENCLAW_SKILLS_DIR")
	if openclawRoot != "" {
		if dirExists(openclawRoot) {
			roots["openclaw"] = []string{openclawRoot}
		}
		return roots, nil
	}

	nvmRoot := filepath.Join(home, ".nvm")
	if !dirExists(nvmRoot) {
		return roots, nil
	}

	openclawRoots, err := findOpenClawSkillDirs(nvmRoot)
	if err != nil {
		return nil, err
	}
	if len(openclawRoots) > 0 {
		roots["openclaw"] = openclawRoots
	}
	return roots, nil
}

func findOpenClawSkillDirs(base string) ([]string, error) {
	var roots []string
	seen := map[string]bool{}

	baseParts := strings.Split(filepath.Clean(base), string(filepath.Separator))
	baseDepth := len(baseParts)

	err := filepath.WalkDir(base, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if !d.IsDir() {
			return nil
		}

		parts := strings.Split(filepath.Clean(path), string(filepath.Separator))
		if len(parts)-baseDepth > 8 {
			return filepath.SkipDir
		}

		if filepath.Base(path) == "skills" && filepath.Base(filepath.Dir(path)) == "openclaw" {
			if !seen[path] {
				seen[path] = true
				roots = append(roots, path)
			}
			return filepath.SkipDir
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("search openclaw skill directories: %w", err)
	}

	sort.Strings(roots)
	return roots, nil
}

func dirExists(path string) bool {
	if path == "" {
		return false
	}
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}
