package issues

import (
	"context"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/wohsj110/atlassian_cli/tools/atk-jira/api"
	"github.com/wohsj110/atlassian_cli/tools/atk-jira/internal/cmd/root"
	atkpresent "github.com/wohsj110/atlassian_cli/tools/atk-jira/internal/present"
	"github.com/wohsj110/atlassian_cli/tools/atk-jira/internal/present/projection"
)

func newHistoryCmd(opts *root.Options) *cobra.Command {
	var maxResults int
	var nextPageToken string
	var fieldsFlag string

	cmd := &cobra.Command{
		Use:   "history <issue-key>",
		Short: "List issue changelog history",
		Long:  "List Jira changelog history for an issue as compact changed-field rows.",
		Example: `  atk-jira issues history PROJ-123
  atk-jira issues history PROJ-123 --id
  atk-jira issues history PROJ-123 --extended
  atk-jira issues history PROJ-123 --fields CREATED,FIELD,TO
  atk-jira issues history PROJ-123 --max 1
  atk-jira issues history PROJ-123 --next-page-token 50`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runHistory(cmd.Context(), opts, args[0], maxResults, nextPageToken, fieldsFlag)
		},
	}

	cmd.Flags().IntVarP(&maxResults, "max", "m", 50, "Maximum number of history groups to return")
	cmd.Flags().StringVar(&nextPageToken, "next-page-token", "", "Token for next page of results")
	cmd.Flags().StringVar(&fieldsFlag, "fields", "", "Comma-separated display columns")

	return cmd
}

func runHistory(ctx context.Context, opts *root.Options, issueKey string, maxResults int, nextPageToken string, fieldsFlag string) error {
	client, err := opts.APIClient()
	if err != nil {
		return err
	}

	idOnly := opts.EmitIDOnly()
	var selected []projection.ColumnSpec
	var projected bool
	if !idOnly {
		selected, projected, err = projection.Resolve(
			ctx,
			atkpresent.IssueHistorySpec,
			opts.IsExtended(),
			fieldsFlag,
			noIssueFieldsFetcher,
			"issues history",
		)
		if err != nil {
			return err
		}
	}

	startAt, err := atkpresent.ParseStartAtToken(nextPageToken)
	if err != nil {
		return err
	}

	page, err := client.GetIssueChangelog(ctx, issueKey, api.IssueChangelogOptions{
		StartAt:    startAt,
		MaxResults: maxResults,
	})
	if err != nil {
		return err
	}

	hasMore, nextToken := computeHistoryPageCursor(page)

	if idOnly {
		ids := make([]string, len(page.Histories))
		for i, history := range page.Histories {
			ids[i] = history.ID
		}
		return atkpresent.EmitIDsWithPaginationToken(opts, ids, hasMore, nextToken)
	}

	rows := atkpresent.FlattenIssueHistory(page.Histories)
	if len(rows) == 0 {
		if hasMore {
			return atkpresent.Emit(opts, atkpresent.PaginationOnlyModel(nextToken))
		}
		return atkpresent.Emit(opts, atkpresent.IssueHistoryPresenter{}.PresentNoIssueHistory(issueKey))
	}

	fulltext := opts.IsFullText() || opts.IsExtended()
	model := atkpresent.IssueHistoryPresenter{}.PresentIssueHistoryWithPagination(rows, opts.IsExtended(), fulltext, hasMore, nextToken)
	if projected {
		projection.ApplyToTableInModel(model, selected)
	}
	return atkpresent.Emit(opts, model)
}

func computeHistoryPageCursor(page *api.IssueChangelogPage) (bool, string) {
	if page == nil {
		return false, ""
	}
	advance := len(page.Histories)
	if advance == 0 && page.MaxResults > 0 {
		// Avoid re-emitting the current offset if Jira returns an empty page mid-set.
		advance = page.MaxResults
	}
	nextStart := page.StartAt + advance
	if nextStart >= page.Total {
		return false, ""
	}
	return true, strconv.Itoa(nextStart)
}

func noIssueFieldsFetcher(context.Context) ([]api.Field, error) {
	return nil, nil
}
