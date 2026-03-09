package convert

import "github.com/spf13/cobra"

func NewCommand() *cobra.Command {
	opts := Options{}

	cmd := &cobra.Command{
		Use:   "convert",
		Short: "Convert skills between supported formats",
		RunE: func(cmd *cobra.Command, args []string) error {
			return Convert(opts)
		},
	}

	cmd.Flags().StringVar(&opts.From, "from", "", "source format: openclaw|claude")
	cmd.Flags().StringVar(&opts.To, "to", "", "target format: openclaw|claude")
	cmd.Flags().StringVar(&opts.Input, "input", "", "input skill path (directory or SKILL.md)")
	cmd.Flags().StringVar(&opts.Output, "output", "", "output skill directory")
	cmd.Flags().BoolVar(&opts.Force, "force", false, "overwrite output directory if it exists")

	_ = cmd.MarkFlagRequired("from")
	_ = cmd.MarkFlagRequired("to")
	_ = cmd.MarkFlagRequired("input")
	_ = cmd.MarkFlagRequired("output")

	return cmd
}
