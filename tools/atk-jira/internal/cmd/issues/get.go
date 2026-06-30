package issues

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/wohsj110/atlassian_cli/shared/present"

	"github.com/wohsj110/atlassian_cli/tools/atk-jira/api"
	"github.com/wohsj110/atlassian_cli/tools/atk-jira/internal/cache"
	"github.com/wohsj110/atlassian_cli/tools/atk-jira/internal/cmd/root"
	atkpresent "github.com/wohsj110/atlassian_cli/tools/atk-jira/internal/present"
	"github.com/wohsj110/atlassian_cli/tools/atk-jira/internal/present/projection"
)

func newGetCmd(opts *root.Options) *cobra.Command {
	var noTruncate bool
	var fieldsFlag string
	var customFields bool

	cmd := &cobra.Command{
		Use:   "get <issue-key> [issue-key...]",
		Short: "Get issue details",
		Long:  "Retrieve and display details for a specific issue, or a summary table when multiple keys are given.",
		Example: `  atk-jira issues get PROJ-123
  atk-jira issues get PROJ-123 PROJ-456 PROJ-789
  atk-jira issues get PROJ-123 --fulltext
  atk-jira issues get PROJ-123 --id
  atk-jira issues get PROJ-123 --fields Status,Assignee
  atk-jira issues get PROJ-123 --fields "Issue Type"
  atk-jira issues get PROJ-123 --custom-fields`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 1 {
				if fieldsFlag != "" {
					return fmt.Errorf("--fields is only supported with a single issue key")
				}
				if customFields {
					return fmt.Errorf("--custom-fields is only supported with a single issue key")
				}
				return runGetMulti(cmd.Context(), opts, args)
			}
			return runGet(cmd.Context(), opts, args[0], noTruncate || opts.IsFullText(), fieldsFlag, customFields)
		},
	}

	cmd.Flags().BoolVar(&noTruncate, "no-truncate", false, "Show full description without truncation")
	_ = cmd.Flags().MarkDeprecated("no-truncate", "use --fulltext instead")
	cmd.Flags().StringVar(&fieldsFlag, "fields", "", "Comma-separated display fields (labels, Jira field IDs, or human names)")
	cmd.Flags().BoolVar(&customFields, "custom-fields", false, "Append custom fields section to output")

	return cmd
}

func runGet(ctx context.Context, opts *root.Options, issueKey string, noTruncate bool, fieldsFlag string, customFields bool) error {
	client, err := opts.APIClient()
	if err != nil {
		return err
	}

	if opts.EmitIDOnly() {
		issue, err := client.GetIssue(ctx, issueKey)
		if err != nil {
			return err
		}
		return atkpresent.EmitIDs(opts, []string{issue.Key})
	}

	selected, projected, err := projection.Resolve(
		ctx,
		atkpresent.IssueDetailSpec,
		opts.IsExtended(),
		fieldsFlag,
		fieldsFetcher(client),
		"issues get",
	)
	if err != nil {
		return err
	}

	issue, err := client.GetIssue(ctx, issueKey)
	if err != nil {
		return err
	}

	presenter := atkpresent.IssuePresenter{}

	if opts.IsExtended() {
		noTruncate = true
	}

	if projected {
		model := presenter.PresentDetailProjection(issue, client.IssueURL(issue.Key), noTruncate)
		atkpresent.AppendDynamicDetailFields(model, issue, projection.DynamicSpecs(selected))
		projection.ApplyToDetailInModel(model, selected)
		if customFields {
			appendCustomFields(ctx, client, issue, model)
		}
		return atkpresent.Emit(opts, model)
	}
	model := presenter.PresentDetail(issue, client.IssueURL(issue.Key), opts.IsExtended(), noTruncate)
	if customFields {
		appendCustomFields(ctx, client, issue, model)
	}
	return atkpresent.Emit(opts, model)
}

func runGetMulti(ctx context.Context, opts *root.Options, issueKeys []string) error {
	client, err := opts.APIClient()
	if err != nil {
		return err
	}

	if opts.EmitIDOnly() {
		ids := make([]string, 0, len(issueKeys))
		for _, key := range issueKeys {
			issue, err := client.GetIssue(ctx, key)
			if err != nil {
				return err
			}
			ids = append(ids, issue.Key)
		}
		return atkpresent.EmitIDs(opts, ids)
	}

	issues := make([]api.Issue, 0, len(issueKeys))
	for _, key := range issueKeys {
		issue, err := client.GetIssue(ctx, key)
		if err != nil {
			return err
		}
		issues = append(issues, *issue)
	}

	model := atkpresent.IssuePresenter{}.PresentList(issues, opts.IsExtended())
	return atkpresent.Emit(opts, model)
}

func appendCustomFields(ctx context.Context, client *api.Client, issue *api.Issue, model *present.OutputModel) {
	fields, err := cache.GetFieldsCacheFirst(ctx, client)
	if err != nil {
		return
	}
	entries := api.ExtractIssueFieldValues(issue, fields)
	atkpresent.AppendCustomFieldsSection(model, entries)
}
