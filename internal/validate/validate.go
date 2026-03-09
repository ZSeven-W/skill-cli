package validate

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validate [path-to-skill-or-skill-md]",
		Short: "Validate a skill definition",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			target := "."
			if len(args) == 1 {
				target = args[0]
			}

			result, err := ValidatePath(target)
			if err != nil {
				return err
			}

			if len(result.Errors) == 0 {
				fmt.Fprintf(os.Stdout, "VALID: %s\n", result.Path)
				return nil
			}

			fmt.Fprintf(os.Stdout, "INVALID: %s\n", result.Path)
			for _, item := range result.Errors {
				fmt.Fprintf(os.Stdout, "- %s\n", item)
			}
			return fmt.Errorf("validation failed")
		},
	}

	return cmd
}
