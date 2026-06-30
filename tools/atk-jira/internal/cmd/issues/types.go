package issues

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/wohsj110/atlassian_cli/tools/atk-jira/internal/cache"
	"github.com/wohsj110/atlassian_cli/tools/atk-jira/internal/cmd/root"
	atkpresent "github.com/wohsj110/atlassian_cli/tools/atk-jira/internal/present"
	"github.com/wohsj110/atlassian_cli/tools/atk-jira/internal/resolve"
)

func newTypesCmd(opts *root.Options) *cobra.Command {
	var project string

	cmd := &cobra.Command{
		Use:   "types",
		Short: "List valid issue types for a project",
		Long:  "List all valid issue types that can be used when creating issues in a specific project.",
		Example: `  # List issue types for a project
  atk-jira issues types --project MYPROJ

  # Emit only type IDs
  atk-jira issues types --project MYPROJ --id`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runTypes(cmd.Context(), opts, project)
		},
	}

	cmd.Flags().StringVarP(&project, "project", "p", "", "Project key (required)")
	_ = cmd.MarkFlagRequired("project")

	return cmd
}

func runTypes(ctx context.Context, opts *root.Options, project string) error {
	client, err := opts.APIClient()
	if err != nil {
		return err
	}

	resolvedProject, err := resolve.New(client).Project(ctx, project)
	if err != nil {
		return err
	}
	projectKey := resolvedProject.Key

	issueTypes, err := cache.GetIssueTypesCacheFirst(ctx, client, projectKey)
	if err != nil {
		return err
	}

	if len(issueTypes) == 0 {
		return atkpresent.Emit(opts, atkpresent.IssuePresenter{}.PresentNoTypes(projectKey))
	}

	if opts.EmitIDOnly() {
		ids := make([]string, len(issueTypes))
		for i, t := range issueTypes {
			ids[i] = t.ID
		}
		return atkpresent.EmitIDs(opts, ids)
	}

	model := atkpresent.IssuePresenter{}.PresentTypes(issueTypes)
	return atkpresent.Emit(opts, model)
}
