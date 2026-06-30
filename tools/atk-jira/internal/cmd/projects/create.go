// Package projects provides CLI commands for managing Jira projects.
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
	"github.com/wohsj110/atlassian_cli/tools/atk-jira/internal/resolve"
)

func newCreateCmd(opts *root.Options) *cobra.Command {
	var key string
	var name string
	var projectType string
	var lead string
	var description string

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new project",
		Long: `Create a new Jira project.

The --lead flag accepts a cached user reference: accountId, email,
display name, or "me". Use 'atk-jira users search' for candidates or
'atk-jira me' to get your own accountId.`,
		Example: `  # Create a software project (--lead resolves via the users cache)
  atk-jira projects create --key MYPROJ --name "My Project" --lead me
  atk-jira projects create --key MYPROJ --name "My Project" --lead "Aaron Wong"
  atk-jira projects create --key MYPROJ --name "My Project" --lead aaron@example.com

  # Create a business project with description
  atk-jira projects create --key BIZ --name "Business" --type business --lead me --description "Business project"

  # Project types: software (default), service_desk, business
  atk-jira projects types`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runCreate(cmd.Context(), opts, key, name, projectType, lead, description)
		},
	}

	cmd.Flags().StringVarP(&key, "key", "k", "", "Project key (required)")
	cmd.Flags().StringVarP(&name, "name", "n", "", "Project name (required)")
	cmd.Flags().StringVarP(&projectType, "type", "t", "software", "Project type (software, service_desk, business)")
	cmd.Flags().StringVarP(&lead, "lead", "l", "", "Lead: accountId, email, display name, or \"me\" (required)")
	cmd.Flags().StringVarP(&description, "description", "d", "", "Project description")

	_ = cmd.MarkFlagRequired("key")
	_ = cmd.MarkFlagRequired("name")
	_ = cmd.MarkFlagRequired("lead")

	return cmd
}

func runCreate(ctx context.Context, opts *root.Options, key, name, projectType, lead, description string) error {
	client, err := opts.APIClient()
	if err != nil {
		return err
	}

	resolvedLead, err := resolve.New(client).User(ctx, lead)
	if err != nil {
		return err
	}

	req := &api.CreateProjectRequest{
		Key:            key,
		Name:           name,
		ProjectTypeKey: projectType,
		LeadAccountID:  resolvedLead.AccountID,
		Description:    description,
	}

	project, err := client.CreateProject(ctx, req)
	if err != nil {
		return err
	}

	_ = cache.Touch(cache.ProjectDependents()...)

	if opts.EmitIDOnly() {
		return atkpresent.EmitIDs(opts, []string{project.Key})
	}

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
		IsFresh: func(model *present.OutputModel) bool {
			return mutation.ModelContainsField(model, "", name)
		},
		Fallback: func(id string) *present.OutputModel {
			return atkpresent.ProjectPresenter{}.PresentCreated(id, name)
		},
	})
}
