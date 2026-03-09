package main

import (
	"fmt"
	"os"

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

	root.AddCommand(
		create.NewCommand(),
		validate.NewCommand(),
		discover.NewListCommand(),
	)

	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
