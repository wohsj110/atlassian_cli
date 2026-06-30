package projects

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/wohsj110/atlassian_cli/shared/present"

	"github.com/wohsj110/atlassian_cli/tools/atk-jira/api"
	"github.com/wohsj110/atlassian_cli/tools/atk-jira/internal/cache"
	"github.com/wohsj110/atlassian_cli/tools/atk-jira/internal/cmd/root"
	"github.com/wohsj110/atlassian_cli/tools/atk-jira/internal/mutation"
	atkpresent "github.com/wohsj110/atlassian_cli/tools/atk-jira/internal/present"
)

func newRestoreCmd(opts *root.Options) *cobra.Command {
	return &cobra.Command{
		Use:     "restore <project-key>",
		Short:   "Restore a deleted project",
		Long:    "Restore a project from the trash.",
		Example: `  atk-jira projects restore MYPROJ`,
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runRestore(cmd.Context(), opts, args[0])
		},
	}
}

func runRestore(ctx context.Context, opts *root.Options, keyOrID string) error {
	client, err := opts.APIClient()
	if err != nil {
		return err
	}

	project, err := client.RestoreProject(ctx, keyOrID)
	if err != nil {
		return err
	}

	_ = cache.Touch(cache.ProjectDependents()...)

	if opts.EmitIDOnly() {
		return atkpresent.EmitIDs(opts, []string{project.Key})
	}

	restoredName := project.Name
	return mutation.WriteAndPresent(ctx, opts, mutation.Config{
		Write: func(_ context.Context) (string, error) {
			return project.Key, nil
		},
		Fetch: func(ctx context.Context, id string) (*present.OutputModel, error) {
			fetched, err := client.GetProject(ctx, id, api.ProjectGetExpand)
			if err != nil {
				return nil, err
			}
			return atkpresent.ProjectPresenter{}.PresentProjectDetail(fetched, opts.IsExtended()), nil
		},
		Fallback: func(id string) *present.OutputModel {
			return atkpresent.ProjectPresenter{}.PresentRestored(id, restoredName)
		},
	})
}
