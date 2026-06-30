// Package sprints provides CLI commands for managing Jira sprints.
package sprints

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/spf13/cobra"

	"github.com/wohsj110/atlassian_cli/tools/atk-jira/api"
	"github.com/wohsj110/atlassian_cli/tools/atk-jira/internal/cache"
	"github.com/wohsj110/atlassian_cli/tools/atk-jira/internal/cmd/root"
	"github.com/wohsj110/atlassian_cli/tools/atk-jira/internal/mutation"
	atkpresent "github.com/wohsj110/atlassian_cli/tools/atk-jira/internal/present"
	"github.com/wohsj110/atlassian_cli/tools/atk-jira/internal/present/projection"
	"github.com/wohsj110/atlassian_cli/tools/atk-jira/internal/resolve"
)

func noFieldFetch(_ context.Context) ([]api.Field, error) { return nil, nil }

func issueFieldsFetcher(client *api.Client) func(context.Context) ([]api.Field, error) {
	return func(ctx context.Context) ([]api.Field, error) {
		return cache.GetFieldsCacheFirst(ctx, client)
	}
}

// validateBoardRef rejects inputs that would parse as numeric but produce a
// synthetic Board{ID: n} with n <= 0, which the downstream Agile endpoints
// return confusing 404s for. Non-numeric names pass through unchanged —
// board-name resolution is handled by the resolver.
func validateBoardRef(board string) error {
	if board == "" {
		return fmt.Errorf("--board is required")
	}
	if n, err := strconv.Atoi(board); err == nil && n <= 0 {
		return fmt.Errorf("--board numeric ID must be positive (got %s)", board)
	}
	return nil
}

// Register registers the sprints commands
func Register(parent *cobra.Command, opts *root.Options) {
	cmd := &cobra.Command{
		Use:     "sprints",
		Aliases: []string{"sprint", "sp"},
		Short:   "Manage sprints",
		Long:    "Commands for viewing sprints and sprint issues.",
		// SupportsAgile checks AgileURL — the correct guard for Agile API commands.
		// Non-Agile scope-restricted commands (automation, dashboards) use IsBearerAuth() instead.
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			// Cobra does not chain PersistentPreRunE — this hook shadows
			// the root's, so we must invoke the backend-selection wiring
			// explicitly. Without this, --backend / keyring.backend silently
			// stop applying on the `sprints` command path.
			if err := root.WireBackendSelection(cmd); err != nil {
				return err
			}
			client, err := opts.APIClient()
			if err != nil {
				return err
			}
			if !client.SupportsAgile() {
				return api.ErrAgileUnavailable
			}
			return nil
		},
	}

	cmd.AddCommand(newListCmd(opts))
	cmd.AddCommand(newCurrentCmd(opts))
	cmd.AddCommand(newIssuesCmd(opts))
	cmd.AddCommand(newAddCmd(opts))
	cmd.AddCommand(newRemoveCmd(opts))

	parent.AddCommand(cmd)
}

func newListCmd(opts *root.Options) *cobra.Command {
	var board string
	var state string
	var maxResults int
	var nextPageToken string
	var fieldsFlag string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List sprints for a board",
		Long:  "List sprints for a specific board. --board accepts a board ID or name.",
		Example: `  # List all sprints
  atk-jira sprints list --board 123
  atk-jira sprints list --board "MON board"

  # List only active sprints
  atk-jira sprints list --board 123 --state active

  # Extended output with completion dates, board, goal
  atk-jira sprints list --board 123 --extended

  # Emit only sprint IDs
  atk-jira sprints list --board 123 --id`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if err := validateBoardRef(board); err != nil {
				return err
			}
			client, err := opts.APIClient()
			if err != nil {
				return err
			}
			resolvedBoard, err := resolve.New(client).Board(cmd.Context(), board)
			if err != nil {
				return err
			}
			return runList(cmd.Context(), opts, client, resolvedBoard.ID, state, maxResults, nextPageToken, fieldsFlag)
		},
	}

	cmd.Flags().StringVarP(&board, "board", "b", "", "Board ID or name (required)")
	cmd.Flags().StringVarP(&state, "state", "s", "", "Filter by state (active, closed, future)")
	cmd.Flags().IntVarP(&maxResults, "max", "m", 50, "Maximum number of results")
	cmd.Flags().StringVar(&nextPageToken, "next-page-token", "", "Decimal startAt for the next page")
	cmd.Flags().StringVar(&fieldsFlag, "fields", "", "Comma-separated display columns")

	return cmd
}

func runList(ctx context.Context, opts *root.Options, client *api.Client, boardID int, state string, maxResults int, nextPageToken, fieldsFlag string) error {
	idOnly := opts.EmitIDOnly()

	startAt, err := atkpresent.ParseStartAtToken(nextPageToken)
	if err != nil {
		return err
	}

	var selected []projection.ColumnSpec
	var projected bool
	if !idOnly {
		selected, projected, err = projection.Resolve(
			ctx,
			atkpresent.SprintListSpec,
			opts.IsExtended(),
			fieldsFlag,
			noFieldFetch,
			"sprints list",
		)
		if err != nil {
			return err
		}
	}

	allSprints, err := cache.GetSprintsCacheFirst(ctx, client, boardID, state, fetchAllSprints)
	if err != nil {
		return err
	}

	atkpresent.SortSprintsForDisplay(allSprints)

	// Client-side pagination window over the sorted slice.
	if maxResults <= 0 {
		maxResults = 50
	}
	total := len(allSprints)
	if startAt > total {
		startAt = total
	}
	end := startAt + maxResults
	if end > total {
		end = total
	}
	page := allSprints[startAt:end]
	hasMore := end < total
	nextToken := ""
	if hasMore {
		nextToken = strconv.Itoa(end)
	}

	if idOnly {
		ids := make([]string, len(page))
		for i, s := range page {
			ids[i] = strconv.Itoa(s.ID)
		}
		return atkpresent.EmitIDsWithPaginationToken(opts, ids, hasMore, nextToken)
	}

	if len(page) == 0 {
		return atkpresent.Emit(opts, atkpresent.SprintPresenter{}.PresentEmpty())
	}

	model := atkpresent.SprintPresenter{}.PresentListWithPagination(page, opts.IsExtended(), hasMore, nextToken)
	if projected {
		projection.ApplyToTableInModel(model, selected)
	}
	return atkpresent.Emit(opts, model)
}

const fetchPageSize = 50

func fetchAllSprints(ctx context.Context, client *api.Client, boardID int, state string) ([]api.Sprint, error) {
	var all []api.Sprint
	startAt := 0
	for {
		result, err := client.ListSprints(ctx, boardID, state, startAt, fetchPageSize)
		if err != nil {
			return nil, err
		}
		if !result.IsLast && len(result.Values) == 0 {
			return nil, fmt.Errorf("unexpected paginated response: IsLast=false with empty values (startAt=%d)", startAt)
		}
		all = append(all, result.Values...)
		if result.IsLast {
			break
		}
		startAt += len(result.Values)
	}
	return all, nil
}

func newCurrentCmd(opts *root.Options) *cobra.Command {
	var board string
	var fieldsFlag string

	cmd := &cobra.Command{
		Use:   "current",
		Short: "Show current sprint",
		Long:  "Show the current active sprint for a board. --board accepts a board ID or name.",
		Example: `  atk-jira sprints current --board 123
  atk-jira sprints current --board "MON board"
  atk-jira sprints current --board 123 --extended`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if err := validateBoardRef(board); err != nil {
				return err
			}
			client, err := opts.APIClient()
			if err != nil {
				return err
			}
			resolvedBoard, err := resolve.New(client).Board(cmd.Context(), board)
			if err != nil {
				return err
			}
			return runCurrent(cmd.Context(), opts, client, &resolvedBoard, fieldsFlag)
		},
	}

	cmd.Flags().StringVarP(&board, "board", "b", "", "Board ID or name (required)")
	cmd.Flags().StringVar(&fieldsFlag, "fields", "", "Comma-separated display fields")

	return cmd
}

func runCurrent(ctx context.Context, opts *root.Options, client *api.Client, board *api.Board, fieldsFlag string) error {
	var selected []projection.ColumnSpec
	var projected bool
	if !opts.EmitIDOnly() {
		var err error
		selected, projected, err = projection.Resolve(
			ctx,
			atkpresent.SprintDetailSpec,
			opts.IsExtended(),
			fieldsFlag,
			noFieldFetch,
			"sprints current",
		)
		if err != nil {
			return err
		}
	}

	sprint, err := client.GetCurrentSprint(ctx, board.ID)
	if err != nil {
		return err
	}

	if opts.EmitIDOnly() {
		return atkpresent.EmitIDs(opts, []string{strconv.Itoa(sprint.ID)})
	}

	// Enrich synthetic board (no name) for table output paths.
	if board.Name == "" {
		if enriched, err := client.GetBoard(ctx, board.ID); err == nil && enriched.ID == board.ID && enriched.Name != "" {
			board = enriched
		}
	}

	presenter := atkpresent.SprintPresenter{}
	if projected {
		model := presenter.PresentDetailProjection(sprint, board)
		projection.ApplyToDetailInModel(model, selected)
		return atkpresent.Emit(opts, model)
	}

	model := presenter.PresentDetail(sprint, board, opts.IsExtended())
	return atkpresent.Emit(opts, model)
}

func newIssuesCmd(opts *root.Options) *cobra.Command {
	var maxResults int
	var nextPageToken string
	var fieldsFlag string

	cmd := &cobra.Command{
		Use:   "issues <sprint>",
		Short: "List issues in a sprint",
		Long:  "List all issues in a specific sprint. Accepts a sprint ID or name (resolved via cache).",
		Example: `  atk-jira sprints issues 456
  atk-jira sprints issues "MON Sprint 70"
  atk-jira sprints issues 456 --fields KEY,STATUS,customfield_10005`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.APIClient()
			if err != nil {
				return err
			}
			resolvedSprint, err := resolve.New(client).Sprint(cmd.Context(), args[0], 0)
			if err != nil {
				return err
			}
			return runIssues(cmd.Context(), opts, resolvedSprint.ID, maxResults, nextPageToken, fieldsFlag)
		},
	}

	cmd.Flags().IntVarP(&maxResults, "max", "m", 50, "Maximum number of results")
	cmd.Flags().StringVar(&nextPageToken, "next-page-token", "", "Decimal startAt for the next page")
	cmd.Flags().StringVar(&fieldsFlag, "fields", "", "Comma-separated display columns (headers, Jira field IDs, or human names)")

	return cmd
}

func runIssues(ctx context.Context, opts *root.Options, sprintID int, maxResults int, nextPageToken string, fieldsFlag string) error {
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
			atkpresent.IssueListSpec,
			opts.IsExtended(),
			fieldsFlag,
			issueFieldsFetcher(client),
			"sprints issues",
		)
		if err != nil {
			return err
		}
	}

	startAt, err := atkpresent.ParseStartAtToken(nextPageToken)
	if err != nil {
		return err
	}

	var fetchFields []string
	if projected {
		fetchFields = projection.DeriveFetchFields(selected)
	}

	result, err := client.GetSprintIssues(ctx, sprintID, startAt, maxResults, fetchFields)
	if err != nil {
		return err
	}

	var hasMore bool
	if result.Total < 0 {
		hasMore = len(result.Issues) == maxResults
	} else {
		hasMore = result.StartAt+len(result.Issues) < result.Total
	}
	nextToken := ""
	if hasMore {
		nextToken = strconv.Itoa(startAt + len(result.Issues))
	}

	if idOnly {
		ids := make([]string, len(result.Issues))
		for i, issue := range result.Issues {
			ids[i] = issue.Key
		}
		return atkpresent.EmitIDsWithPaginationToken(opts, ids, hasMore, nextToken)
	}

	if len(result.Issues) == 0 {
		return atkpresent.Emit(opts, atkpresent.SprintPresenter{}.PresentNoIssues())
	}

	model := atkpresent.IssuePresenter{}.PresentListWithPagination(result.Issues, opts.IsExtended(), hasMore, nextToken)
	if projected {
		atkpresent.AppendDynamicTableColumns(model, result.Issues, projection.DynamicSpecs(selected))
		projection.ApplyToTableInModel(model, selected)
	}
	return atkpresent.Emit(opts, model)
}

func newAddCmd(opts *root.Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add <sprint> <issue-key>...",
		Short: "Move issues to a sprint",
		Long:  "Move one or more issues to a specific sprint. <sprint> accepts a sprint ID or name.",
		Example: `  # Move a single issue by sprint ID
  atk-jira sprints add 123 PROJ-456

  # Move by sprint name (resolved via cache)
  atk-jira sprints add "MON Sprint 70" PROJ-456

  # Move multiple issues
  atk-jira sprints add 123 PROJ-456 PROJ-789 PROJ-101`,
		Args: cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.APIClient()
			if err != nil {
				return err
			}
			resolvedSprint, err := resolve.New(client).Sprint(cmd.Context(), args[0], 0)
			if err != nil {
				return err
			}
			return runAdd(cmd.Context(), opts, client, resolvedSprint.ID, args[1:])
		},
	}

	return cmd
}

func runAdd(ctx context.Context, opts *root.Options, client *api.Client, sprintID int, issueKeys []string) error {
	if err := client.MoveIssuesToSprint(ctx, sprintID, issueKeys); err != nil {
		return err
	}

	if opts.EmitIDOnly() {
		return atkpresent.EmitIDs(opts, issueKeys)
	}

	// Verify membership via GetSprintIssues and present matched issues.
	keySet := make(map[string]bool, len(issueKeys))
	for _, k := range issueKeys {
		keySet[k] = true
	}

	var matched []api.Issue
	for i, delay := range mutation.BackoffSchedule {
		if i > 0 && delay > 0 {
			select {
			case <-ctx.Done():
				goto fallback
			case <-time.After(delay):
			}
		}

		matched = nil
		found := make(map[string]bool)
		startAt := 0
		for {
			result, err := client.GetSprintIssues(ctx, sprintID, startAt, 50, nil)
			if err != nil {
				break
			}
			for _, issue := range result.Issues {
				if keySet[issue.Key] {
					matched = append(matched, issue)
					found[issue.Key] = true
				}
			}
			if len(found) == len(issueKeys) {
				break
			}
			if len(result.Issues) == 0 || startAt+len(result.Issues) >= result.Total {
				break
			}
			startAt += len(result.Issues)
		}
		if len(found) == len(issueKeys) {
			break
		}
	}

	if len(matched) == len(issueKeys) {
		return atkpresent.Emit(opts, atkpresent.IssuePresenter{}.PresentList(matched, opts.IsExtended()))
	}

fallback:
	_ = atkpresent.Emit(opts, atkpresent.SprintPresenter{}.PresentPostStateUnavailable())
	return atkpresent.Emit(opts, atkpresent.SprintPresenter{}.PresentMoved(issueKeys, sprintID))
}

func newRemoveCmd(opts *root.Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove <issue-key>...",
		Short: "Move issues to the backlog",
		Long:  "Move one or more issues from their current sprint to the backlog.",
		Example: `  # Move a single issue to backlog
  atk-jira sprints remove PROJ-456

  # Move multiple issues
  atk-jira sprints remove PROJ-456 PROJ-789 PROJ-101`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.APIClient()
			if err != nil {
				return err
			}
			return runRemove(cmd.Context(), opts, client, args)
		},
	}

	return cmd
}

func runRemove(ctx context.Context, opts *root.Options, client *api.Client, issueKeys []string) error {
	if err := client.MoveIssuesToBacklog(ctx, issueKeys); err != nil {
		return err
	}

	if opts.EmitIDOnly() {
		return atkpresent.EmitIDs(opts, issueKeys)
	}

	var matched []api.Issue
	for i, delay := range mutation.BackoffSchedule {
		if i > 0 && delay > 0 {
			select {
			case <-ctx.Done():
				goto fallback
			case <-time.After(delay):
			}
		}

		matched = nil
		for _, key := range issueKeys {
			issue, err := client.GetIssue(ctx, key)
			if err != nil {
				matched = nil
				break
			}
			matched = append(matched, *issue)
		}
		if len(matched) == len(issueKeys) {
			break
		}
	}

	if len(matched) == len(issueKeys) {
		return atkpresent.Emit(opts, atkpresent.IssuePresenter{}.PresentList(matched, opts.IsExtended()))
	}

fallback:
	_ = atkpresent.Emit(opts, atkpresent.SprintPresenter{}.PresentPostStateUnavailable())
	return atkpresent.Emit(opts, atkpresent.SprintPresenter{}.PresentMovedToBacklog(issueKeys))
}
