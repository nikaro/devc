package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish|powershell]",
	Short: "Generate completion script",
	Long: `To load completions:

Bash:

  $ source <(devc completion bash)

  # To load completions for each session, execute once:
  # Linux:
  $ devc completion bash > /etc/bash_completion.d/devc
  # macOS:
  $ devc completion bash > /usr/local/etc/bash_completion.d/devc

Zsh:

  # If shell completion is not already enabled in your environment,
  # you will need to enable it.  You can execute the following once:

  $ echo "autoload -U compinit; compinit" >> ~/.zshrc

  # To load completions for each session, execute once:
  $ devc completion zsh > "${fpath[1]}/_devc"

  # You will need to start a new shell for this setup to take effect.

fish:

  $ devc completion fish | source

  # To load completions for each session, execute once:
  $ devc completion fish > ~/.config/fish/completions/devc.fish

PowerShell:

  PS> devc completion powershell | Out-String | Invoke-Expression

  # To load completions for every new session, run:
  PS> devc completion powershell > devc.ps1
  # and source this file from your PowerShell profile.
`,
	DisableFlagsInUseLine: true,
	ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
	Args:                  cobra.ExactValidArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		switch args[0] {
		case "bash":
			cmd.Root().GenBashCompletion(os.Stdout)
		case "zsh":
			cmd.Root().GenZshCompletion(os.Stdout)
		case "fish":
			cmd.Root().GenFishCompletion(os.Stdout, true)
		case "powershell":
			cmd.Root().GenPowerShellCompletionWithDesc(os.Stdout)
		}
	},
}

func init() {
	rootCmd.AddCommand(completionCmd)
}
