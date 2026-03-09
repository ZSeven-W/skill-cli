package discover

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func NewListCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List installed skills from supported platforms",
		RunE: func(cmd *cobra.Command, args []string) error {
			skills, err := DiscoverInstalledSkills()
			if err != nil {
				return err
			}

			if len(skills) == 0 {
				fmt.Fprintln(os.Stdout, "No installed skills found")
				return nil
			}

			fmt.Fprintln(os.Stdout, "PLATFORM\tNAME\tPATH")
			for _, skill := range skills {
				fmt.Fprintf(os.Stdout, "%s\t%s\t%s\n", skill.Platform, skill.Name, skill.Path)
			}
			return nil
		},
	}
}
