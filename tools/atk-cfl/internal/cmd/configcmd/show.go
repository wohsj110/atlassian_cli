package configcmd

import (
	"github.com/spf13/cobra"

	"github.com/wohsj110/atlassian_cli/shared/credstore"
	"github.com/wohsj110/atlassian_cli/shared/keyring"

	"github.com/wohsj110/atlassian_cli/tools/atk-cfl/internal/cmd/root"
	"github.com/wohsj110/atlassian_cli/tools/atk-cfl/internal/config"
	atkpresent "github.com/wohsj110/atlassian_cli/tools/atk-cfl/internal/present"
)

func newShowCmd(opts *root.Options) *cobra.Command {
	return &cobra.Command{
		Use:   "show",
		Short: "Show current configuration",
		Long: `Display the current atk-cfl configuration.

The API token value is never displayed — only whether one is configured,
where it resolves from, and the OS keyring backend in use. Token/keyring
reporting is authoritative.

Note: the non-secret rows (URL, email, etc.) reflect environment
variables and the shared atlassian-agent-cli config file.`,
		Example: `  # Show current configuration
  atk-cfl config show`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runShow(cmd, opts)
		},
	}
}

func runShow(cmd *cobra.Command, opts *root.Options) error {
	configPath := config.DefaultConfigPath()
	if cmd != nil {
		if p, err := cmd.Root().PersistentFlags().GetString("config"); err == nil && p != "" {
			configPath = p
		}
	}

	// Load config file (if exists)
	fileCfg, fileErr := config.Load(configPath)
	if fileErr != nil {
		fileCfg = &config.Config{}
	}

	// Non-secret keyring description (non-migrating: show stays usable
	// during an unresolved §1.8 conflict). The token VALUE is never
	// shown — presence + source only (§1.12).
	kr, krErr := keyring.InspectForTool(credstore.ToolAtkCFL)

	proj := config.ProjectShow(configPath, fileCfg, fileErr, kr, krErr)
	return atkpresent.Emit(opts, atkpresent.ConfigShowPresenter{}.PresentDetail(proj))
}
