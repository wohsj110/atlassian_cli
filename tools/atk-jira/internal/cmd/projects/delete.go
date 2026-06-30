package projects

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/wohsj110/atlassian_cli/shared/present"
	"github.com/wohsj110/atlassian_cli/shared/prompt"

	"github.com/wohsj110/atlassian_cli/tools/atk-jira/internal/cache"
	"github.com/wohsj110/atlassian_cli/tools/atk-jira/internal/cmd/root"
	atkpresent "github.com/wohsj110/atlassian_cli/tools/atk-jira/internal/present"
)

func newDeleteCmd(opts *root.Options) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "delete <project-key>",
		Short: "Delete a project",
		Long: `Soft-delete a Jira project (moves it to trash).

The project can be restored from trash using 'atk-jira projects restore'.`,
		Example: `  # Delete a project (will prompt for confirmation)
  atk-jira projects delete MYPROJ

  # Delete without confirmation
  atk-jira projects delete MYPROJ --force`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDelete(cmd.Context(), opts, args[0], force)
		},
	}

	cmd.Flags().BoolVar(&force, "force", false, "Skip confirmation prompt")

	return cmd
}

func runDelete(ctx context.Context, opts *root.Options, keyOrID string, force bool) error {
	if !force && !opts.NonInteractive {
		fmt.Fprintf(opts.Stderr, "This will delete project %s (moves to trash). It can be restored later.\n", keyOrID)
		fmt.Fprint(opts.Stderr, "Are you sure? [y/N]: ")
	}
	confirmed, err := prompt.ConfirmOrFail(force, opts.NonInteractive, opts.Stdin)
	if err != nil {
		return err
	}
	if !confirmed {
		model := atkpresent.ProjectPresenter{}.PresentDeleteCancelled()
		out := present.Render(model, opts.RenderStyle())
		_, _ = fmt.Fprint(opts.Stdout, out.Stdout)
		return nil
	}

	client, err := opts.APIClient()
	if err != nil {
		return err
	}

	if err := client.DeleteProject(ctx, keyOrID); err != nil {
		return err
	}

	_ = cache.Touch(cache.ProjectDependents()...)

	model := atkpresent.ProjectPresenter{}.PresentDeleted(keyOrID)
	out := present.Render(model, opts.RenderStyle())
	_, _ = fmt.Fprint(opts.Stdout, out.Stdout)
	return nil
}
