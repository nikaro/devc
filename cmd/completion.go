package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var output string

var completionCmd = &cobra.Command{
	Use:   "completion",
	Short: "Generate completion script",
	Long: `To load completion run:

. <(devc completion)

To configure your bash shell to load completions for each session add to your bashrc:

# ~/.bashrc or ~/.profile
. <(devc completion)

If you have bash-completion installed:

devc completion --output /etc/bash_completion.d/devc
`,
	Run: func(cmd *cobra.Command, args []string) {
		out := os.Stdout
		if output != "" {
			f, err := os.Create(output)
			if err != nil {
				panic(err)
			}
			defer f.Close()
			out = f
		}
		rootCmd.GenBashCompletion(out)
	},
}

func init() {
	completionCmd.PersistentFlags().StringVarP(&output, "output", "o", "", "output file (default stdout)")

	rootCmd.AddCommand(completionCmd)
}
