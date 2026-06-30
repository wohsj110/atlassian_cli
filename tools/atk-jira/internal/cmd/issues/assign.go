// Package issues provides CLI commands for managing Jira issues.
package issues

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/wohsj110/atlassian_cli/shared/present"

	"github.com/wohsj110/atlassian_cli/tools/atk-jira/internal/cmd/root"
	"github.com/wohsj110/atlassian_cli/tools/atk-jira/internal/mutation"
	atkpresent "github.com/wohsj110/atlassian_cli/tools/atk-jira/internal/present"
	"github.com/wohsj110/atlassian_cli/tools/atk-jira/internal/resolve"
)

func newAssignCmd(opts *root.Options) *cobra.Command {
	var unassign bool

	cmd := &cobra.Command{
		Use:   "assign <issue-key> [user]",
		Short: "Assign an issue to a user",
		Long:  `Assign an issue to a user, or unassign it. The <user> argument accepts an accountId, email, display name, or "me" — it is resolved via the instance cache.`,
		Example: `  # Assign by display name, email, "me", or raw accountId
  atk-jira issues assign PROJ-123 "Aaron Wong"
  atk-jira issues assign PROJ-123 aaron@example.com
  atk-jira issues assign PROJ-123 me
  atk-jira issues assign PROJ-123 5b10ac8d82e05b22cc7d4ef5

  # Unassign an issue
  atk-jira issues assign PROJ-123 --unassign`,
		Args: cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			accountID := ""
			if len(args) > 1 {
				accountID = args[1]
			}
			return runAssign(cmd.Context(), opts, args[0], accountID, unassign)
		},
	}

	cmd.Flags().BoolVar(&unassign, "unassign", false, "Remove current assignee")

	return cmd
}

func runAssign(ctx context.Context, opts *root.Options, issueKey, userInput string, unassign bool) error {
	client, err := opts.APIClient()
	if err != nil {
		return err
	}

	accountID := ""
	displayName := ""

	if !unassign && userInput != "" {
		resolvedUser, err := resolve.New(client).User(ctx, userInput)
		if err != nil {
			return err
		}
		accountID = resolvedUser.AccountID
		displayName = resolvedUser.DisplayName
		if displayName == "" {
			// Pass-through path: resolver returned synthetic api.User with only
			// AccountID populated. Fall back to the raw ID in the message.
			displayName = accountID
		}
	}

	if opts.EmitIDOnly() {
		if err := client.AssignIssue(ctx, issueKey, accountID); err != nil {
			return err
		}
		return atkpresent.EmitIDs(opts, []string{issueKey})
	}

	var isFresh func(*present.OutputModel) bool
	if displayName != "" {
		isFresh = func(m *present.OutputModel) bool {
			return mutation.ModelContainsField(m, "Assignee: ", displayName)
		}
	} else {
		isFresh = func(m *present.OutputModel) bool {
			return mutation.ModelContainsField(m, "Assignee: ", "-")
		}
	}

	return mutation.WriteAndPresent(ctx, opts, mutation.Config{
		Write: func(ctx context.Context) (string, error) {
			return issueKey, client.AssignIssue(ctx, issueKey, accountID)
		},
		Fetch: func(ctx context.Context, id string) (*present.OutputModel, error) {
			issue, err := client.GetIssue(ctx, id)
			if err != nil {
				return nil, err
			}
			return atkpresent.IssuePresenter{}.PresentDetail(
				issue, client.IssueURL(id), opts.IsExtended(), opts.IsFullText(),
			), nil
		},
		IsFresh: isFresh,
		Fallback: func(id string) *present.OutputModel {
			return atkpresent.IssuePresenter{}.PresentAssigned(id, displayName)
		},
	})
}
