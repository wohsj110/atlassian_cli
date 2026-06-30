package issues

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/wohsj110/atlassian_cli/shared/present"
	"github.com/wohsj110/atlassian_cli/shared/prompt"

	"github.com/wohsj110/atlassian_cli/tools/atk-jira/internal/cmd/root"
	atkpresent "github.com/wohsj110/atlassian_cli/tools/atk-jira/internal/present"
)

func newDeleteCmd(opts *root.Options) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "delete <issue-key> [<issue-key>...]",
		Short: "Delete one or more issues",
		Long:  "Permanently delete one or more Jira issues. This action cannot be undone.",
		Example: `  # Delete an issue (will prompt for confirmation)
  atk-jira issues delete PROJ-123

  # Delete multiple issues
  atk-jira issues delete PROJ-123 PROJ-124 PROJ-125

  # Delete without confirmation
  atk-jira issues delete PROJ-123 --force`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDelete(cmd.Context(), opts, args, force)
		},
	}

	cmd.Flags().BoolVar(&force, "force", false, "Skip confirmation prompt")

	return cmd
}

func runDelete(ctx context.Context, opts *root.Options, issueKeys []string, force bool) error {
	if !force && !opts.NonInteractive {
		var msg string
		if len(issueKeys) == 1 {
			msg = fmt.Sprintf("This will permanently delete issue %s. This action cannot be undone.", issueKeys[0])
		} else {
			msg = fmt.Sprintf("This will permanently delete %d issues: %s. This action cannot be undone.", len(issueKeys), strings.Join(issueKeys, ", "))
		}
		fmt.Fprintln(opts.Stderr, msg)
		fmt.Fprint(opts.Stderr, "Are you sure? [y/N]: ")
	}
	confirmed, err := prompt.ConfirmOrFail(force, opts.NonInteractive, opts.Stdin)
	if err != nil {
		return err
	}
	if !confirmed {
		model := atkpresent.IssuePresenter{}.PresentDeleteCancelled()
		out := present.Render(model, opts.RenderStyle())
		_, _ = fmt.Fprint(opts.Stdout, out.Stdout)
		return nil
	}

	client, err := opts.APIClient()
	if err != nil {
		return err
	}

	var failed int
	for _, key := range issueKeys {
		if ctx.Err() != nil {
			return ctx.Err()
		}

		if err := client.DeleteIssue(ctx, key); err != nil {
			fmt.Fprintf(opts.Stderr, "Failed to delete %s: %s\n", key, err)
			failed++
			continue
		}

		model := atkpresent.IssuePresenter{}.PresentDeleted(key)
		out := present.Render(model, opts.RenderStyle())
		_, _ = fmt.Fprint(opts.Stdout, out.Stdout)
		_, _ = fmt.Fprint(opts.Stderr, out.Stderr)
	}

	if failed > 0 {
		return fmt.Errorf("%w (%d of %d issue(s) failed)", root.ErrAlreadyReported, failed, len(issueKeys))
	}
	return nil
}
