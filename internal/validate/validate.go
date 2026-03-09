package validate

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	var outputFormat string
	var strict bool

	cmd := &cobra.Command{
		Use:   "validate [path-to-skill-or-skill-md]",
		Short: "Validate a skill definition",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			target := "."
			if len(args) == 1 {
				target = args[0]
			}

			if outputFormat != "text" && outputFormat != "json" {
				return fmt.Errorf("unsupported --format %q (use text or json)", outputFormat)
			}

			result, err := ValidatePath(target, strict)
			if err != nil {
				return err
			}

			if outputFormat == "json" {
				payload := struct {
					Path     string   `json:"path"`
					Valid    bool     `json:"valid"`
					Strict   bool     `json:"strict"`
					Errors   []string `json:"errors"`
					Warnings []string `json:"warnings"`
				}{
					Path:     result.Path,
					Valid:    !result.Failed(strict),
					Strict:   strict,
					Errors:   result.Errors,
					Warnings: result.Warnings,
				}
				enc := json.NewEncoder(os.Stdout)
				enc.SetIndent("", "  ")
				if err := enc.Encode(payload); err != nil {
					return fmt.Errorf("encode json output: %w", err)
				}
			} else {
				printTextResult(result, strict)
			}

			if result.Failed(strict) {
				return fmt.Errorf("validation failed")
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&outputFormat, "format", "text", "output format: text|json")
	cmd.Flags().BoolVar(&strict, "strict", false, "treat best-practice warnings as errors")

	return cmd
}

func printTextResult(result Result, strict bool) {
	if result.Failed(strict) {
		fmt.Fprintf(os.Stdout, "INVALID: %s\n", result.Path)
	} else {
		fmt.Fprintf(os.Stdout, "VALID: %s\n", result.Path)
	}

	if len(result.Errors) > 0 {
		fmt.Fprintln(os.Stdout, "Errors:")
		for _, item := range result.Errors {
			fmt.Fprintf(os.Stdout, "- %s\n", item)
		}
	}

	if len(result.Warnings) > 0 {
		fmt.Fprintln(os.Stdout, "Warnings:")
		for _, item := range result.Warnings {
			fmt.Fprintf(os.Stdout, "- %s\n", item)
		}
		if strict {
			fmt.Fprintln(os.Stdout, "Strict mode is enabled: warnings fail validation.")
		}
	}
}
