package main

import (
	"fmt"
	"os"

	"github.com/fini/skill-cli/internal/create"
	"github.com/spf13/cobra"
)

func main() {
	root := &cobra.Command{
		Use:   "skill-cli",
		Short: "Manage AI agent skills across platforms",
	}

	root.AddCommand(create.NewCommand())

	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
