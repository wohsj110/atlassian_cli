package automation

import (
	"github.com/spf13/cobra"

	"github.com/wohsj110/atlassian_cli/tools/atk-jira/internal/cmd/root"
)

func newDisableCmd(opts *root.Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "disable <rule-id>",
		Short: "Disable an automation rule",
		Long:  "Disable an enabled automation rule. This is a safe operation that does not modify the rule definition.",
		Example: `  atk-jira automation disable 12345
  atk-jira auto disable 12345`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSetState(cmd.Context(), opts, args[0], false)
		},
	}

	return cmd
}
