package discover

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func NewListCommand() *cobra.Command {
	var format string
	var platform string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List installed skills from supported platforms",
		RunE: func(cmd *cobra.Command, args []string) error {
			outputFormat := strings.ToLower(strings.TrimSpace(format))
			if outputFormat != "text" && outputFormat != "json" {
				return fmt.Errorf("unsupported format %q (expected text or json)", format)
			}

			platformFilter := strings.ToLower(strings.TrimSpace(platform))
			if platformFilter != "" && platformFilter != "claude" && platformFilter != "openclaw" {
				return fmt.Errorf("unsupported platform %q (expected claude or openclaw)", platform)
			}

			skills, err := DiscoverInstalledSkills()
			if err != nil {
				return err
			}

			if platformFilter != "" {
				filtered := make([]InstalledSkill, 0, len(skills))
				for _, skill := range skills {
					if skill.Platform == platformFilter {
						filtered = append(filtered, skill)
					}
				}
				skills = filtered
			}

			out := cmd.OutOrStdout()
			if outputFormat == "json" {
				encoder := json.NewEncoder(out)
				return encoder.Encode(skills)
			}

			if len(skills) == 0 {
				fmt.Fprintln(out, "No installed skills found")
				return nil
			}

			fmt.Fprintln(out, "PLATFORM\tNAME\tPATH")
			for _, skill := range skills {
				fmt.Fprintf(out, "%s\t%s\t%s\n", skill.Platform, skill.Name, skill.Path)
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&format, "format", "text", "Output format: text or json")
	cmd.Flags().StringVar(&platform, "platform", "", "Filter by platform (claude or openclaw)")
	return cmd
}
