package main

import (
	"fmt"
	"os"

	"github.com/fini/skill-cli/internal/completion"
	"github.com/fini/skill-cli/internal/convert"
	"github.com/fini/skill-cli/internal/create"
	"github.com/fini/skill-cli/internal/discover"
	"github.com/fini/skill-cli/internal/validate"
	"github.com/spf13/cobra"
)

func main() {
	root := &cobra.Command{
		Use:   "skill-cli",
		Short: "Manage AI agent skills across platforms",
	}
	root.CompletionOptions.DisableDefaultCmd = true

	root.AddCommand(
		completion.NewCommand(root),
		create.NewCommand(),
		validate.NewCommand(),
		discover.NewListCommand(),
		convert.NewCommand(),
	)

	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
