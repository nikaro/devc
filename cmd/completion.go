package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var completionOutput string

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
		if completionOutput != "" {
			file, err := os.Create(completionOutput)
			if err != nil {
				panic(err)
			}
			defer file.Close()
			out = file
		}
		rootCmd.GenBashCompletion(out)
	},
}

func init() {
	completionCmd.PersistentFlags().StringVarP(&completionOutput, "output", "o", "", "output file (default stdout)")

	rootCmd.AddCommand(completionCmd)
}
