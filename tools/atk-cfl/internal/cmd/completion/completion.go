// Package completion provides shell completion generation commands.
package completion

import (
	"github.com/spf13/cobra"

	"github.com/wohsj110/atlassian_cli/tools/atk-cfl/internal/cmd/root"
)

// Register adds the completion command to the root command.
func Register(rootCmd *cobra.Command, _ *root.Options) {
	cmd := &cobra.Command{
		Use:   "completion [bash|zsh|fish|powershell]",
		Short: "Generate shell completion scripts",
		Long: `Generate shell completion scripts for atk-cfl.

To load completions:

Bash:
  $ source <(atk-cfl completion bash)
  # To load completions for each session, execute once:
  # Linux:
  $ atk-cfl completion bash > /etc/bash_completion.d/atk-cfl
  # macOS:
  $ atk-cfl completion bash > $(brew --prefix)/etc/bash_completion.d/atk-cfl

Zsh:
  # If shell completion is not already enabled in your environment,
  # you will need to enable it. You can execute the following once:
  $ echo "autoload -U compinit; compinit" >> ~/.zshrc
  # To load completions for each session, execute once:
  $ atk-cfl completion zsh > "${fpath[1]}/_atk-cfl"
  # You will need to start a new shell for this setup to take effect.

Fish:
  $ atk-cfl completion fish | source
  # To load completions for each session, execute once:
  $ atk-cfl completion fish > ~/.config/fish/completions/atk-cfl.fish

PowerShell:
  PS> atk-cfl completion powershell | Out-String | Invoke-Expression
  # To load completions for every new session, run:
  PS> atk-cfl completion powershell > atk-cfl.ps1
  # and source this file from your PowerShell profile.
`,
		DisableFlagsInUseLine: true,
		ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
		Args:                  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
		RunE: func(cmd *cobra.Command, args []string) error {
			out := cmd.OutOrStdout()
			switch args[0] {
			case "bash":
				return cmd.Root().GenBashCompletion(out)
			case "zsh":
				return cmd.Root().GenZshCompletion(out)
			case "fish":
				return cmd.Root().GenFishCompletion(out, true)
			case "powershell":
				return cmd.Root().GenPowerShellCompletionWithDesc(out)
			}
			return nil
		},
	}

	rootCmd.AddCommand(cmd)
}
