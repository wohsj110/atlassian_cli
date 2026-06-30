package automation

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/wohsj110/atlassian_cli/tools/atk-jira/internal/cmd/root"
	atkpresent "github.com/wohsj110/atlassian_cli/tools/atk-jira/internal/present"
)

func newGetCmd(opts *root.Options) *cobra.Command {
	var showComponents bool

	cmd := &cobra.Command{
		Use:   "get <rule-id>",
		Short: "Get automation rule details",
		Long: `Retrieve and display details for a specific automation rule.

Shows rule identifier, name, state, components summary, and description.
Use --show-components to see component type details.
Use --extended for additional fields (labels, tags, author, scope, timestamps).

For the exact JSON needed for editing, use 'atk-jira auto export' instead.`,
		Example: `  atk-jira automation get 12345
  atk-jira auto get 12345 --show-components
  atk-jira auto get 12345 --extended`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runGet(cmd.Context(), opts, args[0], showComponents)
		},
	}

	cmd.Flags().BoolVar(&showComponents, "show-components", false, "Show component type details")

	return cmd
}

func runGet(ctx context.Context, opts *root.Options, ruleID string, showComponents bool) error {
	client, err := opts.APIClient()
	if err != nil {
		return err
	}

	rule, err := client.GetAutomationRule(ctx, ruleID)
	if err != nil {
		return err
	}

	if opts.EmitIDOnly() {
		return atkpresent.EmitIDs(opts, []string{rule.Identifier()})
	}

	presenter := atkpresent.AutomationPresenter{}

	if opts.IsExtended() {
		authorName := ""
		if rule.AuthorAccountID != "" {
			user, err := client.GetUser(ctx, rule.AuthorAccountID, "")
			if err == nil && user.DisplayName != "" {
				authorName = user.DisplayName
			} else {
				authorName = rule.AuthorAccountID
			}
		}
		return atkpresent.Emit(opts, presenter.PresentGetDetailExtended(rule, showComponents, authorName))
	}

	return atkpresent.Emit(opts, presenter.PresentGetDetail(rule, showComponents))
}
