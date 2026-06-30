// Package users provides CLI commands for searching Jira users.
package users

import (
	"context"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/wohsj110/atlassian_cli/tools/atk-jira/api"
	"github.com/wohsj110/atlassian_cli/tools/atk-jira/internal/cache"
	"github.com/wohsj110/atlassian_cli/tools/atk-jira/internal/cmd/root"
	atkpresent "github.com/wohsj110/atlassian_cli/tools/atk-jira/internal/present"
	"github.com/wohsj110/atlassian_cli/tools/atk-jira/internal/present/projection"
)

// noFieldFetch is the projection.Resolve fetcher for user commands. Users are
// not Jira issue fields, so there is no metadata to fetch; returning nil
// routes any deferred tokens cleanly to UnknownFieldError rather than into a
// real /rest/api/3/field call that would surface unrelated issue-field names.
// Package-local rather than shared because it is trivial and typed for the
// caller's context — consolidating it would obscure the intent.
func noFieldFetch(_ context.Context) ([]api.Field, error) { return nil, nil }

// Register registers the users commands
func Register(parent *cobra.Command, opts *root.Options) {
	cmd := &cobra.Command{
		Use:     "users",
		Aliases: []string{"user", "u"},
		Short:   "Search and lookup users",
		Long:    "Commands for searching and looking up Jira users.",
	}

	cmd.AddCommand(newGetCmd(opts))
	cmd.AddCommand(newSearchCmd(opts))

	parent.AddCommand(cmd)
}

func newGetCmd(opts *root.Options) *cobra.Command {
	var fieldsFlag string

	cmd := &cobra.Command{
		Use:   "get <account-id>",
		Short: "Get user details by account ID",
		Long:  "Retrieve and display details for a specific user by their Jira account ID.",
		Example: `  # Get user details (pipe one-liner)
  atk-jira users get 61292e4c4f29230069621c5f

  # Include timezone, locale, and group/application-role counts
  atk-jira users get 61292e4c4f29230069621c5f --extended

  # Just the account ID
  atk-jira users get 61292e4c4f29230069621c5f --id

  # Restrict output to selected fields
  atk-jira users get 61292e4c4f29230069621c5f --fields NAME,EMAIL`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runGet(cmd.Context(), opts, args[0], fieldsFlag)
		},
	}

	cmd.Flags().StringVar(&fieldsFlag, "fields", "", "Comma-separated display fields (UserDetailSpec headers)")

	return cmd
}

func runGet(ctx context.Context, opts *root.Options, accountID, fieldsFlag string) error {
	client, err := opts.APIClient()
	if err != nil {
		return err
	}

	// --id wins: collapse output before we do any projection work. The
	// account ID is its own canonical identifier (no ID→key lookup step like
	// projects), so we can round-trip the caller's input without fetching.
	// Saves an API call and an expand=groups,applicationRoles query that
	// would be thrown away immediately.
	if opts.EmitIDOnly() {
		return atkpresent.EmitIDs(opts, []string{accountID})
	}

	selected, projected, err := projection.Resolve(
		ctx,
		atkpresent.UserDetailSpec,
		opts.IsExtended(),
		fieldsFlag,
		noFieldFetch,
		"users get",
	)
	if err != nil {
		return err
	}

	expand := ""
	if opts.IsExtended() {
		expand = api.UserExtendedExpand
	}
	user, err := cache.GetUserCacheFirst(ctx, client, accountID, expand)
	if err != nil {
		return err
	}

	presenter := atkpresent.UserPresenter{}
	if projected {
		model := presenter.PresentUserDetailProjection(user)
		projection.ApplyToDetailInModel(model, selected)
		return atkpresent.Emit(opts, model)
	}

	var model = presenter.PresentUserOneLiner(user)
	if opts.IsExtended() {
		model = presenter.PresentUserExtended(user)
	}
	return atkpresent.Emit(opts, model)
}

func newSearchCmd(opts *root.Options) *cobra.Command {
	var maxResults int
	var nextPageToken string
	var fieldsFlag string

	cmd := &cobra.Command{
		Use:   "search <query>",
		Short: "Search for users",
		Long: `Search for users by name, email, or username.

The search is case-insensitive and matches against display name, email
address, and other user attributes. Pagination uses /user/search's offset
model — the continuation token is the next startAt as a decimal string. The
endpoint does not return an authoritative terminator, so hasMore is inferred
from the page being full (len(results) == --max).`,
		Example: `  # Search for users named "john"
  atk-jira users search john

  # Include timezone and locale columns
  atk-jira users search john --extended

  # Emit just account IDs, one per line
  atk-jira users search john --id

  # Project output to selected columns
  atk-jira users search john --fields ACCOUNT_ID,NAME

  # Fetch the second page
  atk-jira users search john --max 10 --next-page-token 10`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSearch(cmd.Context(), opts, args[0], maxResults, nextPageToken, fieldsFlag)
		},
	}

	cmd.Flags().IntVarP(&maxResults, "max", "m", 50, "Maximum number of results")
	cmd.Flags().StringVar(&nextPageToken, "next-page-token", "", "Decimal startAt for the next page")
	cmd.Flags().StringVar(&fieldsFlag, "fields", "", "Comma-separated display columns (UserListSpec headers)")

	return cmd
}

func runSearch(ctx context.Context, opts *root.Options, query string, maxResults int, nextPageToken, fieldsFlag string) error {
	client, err := opts.APIClient()
	if err != nil {
		return err
	}

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
			atkpresent.UserListSpec,
			opts.IsExtended(),
			fieldsFlag,
			noFieldFetch,
			"users search",
		)
		if err != nil {
			return err
		}
	}

	rawUsers, err := client.SearchUsers(ctx, query, startAt, maxResults)
	if err != nil {
		return err
	}
	users := filterSearchUsers(rawUsers)

	// /user/search has no native isLast; the heuristic is that a full page
	// implies more pages may exist. Over-reporting in the last window is the
	// documented tradeoff for a command whose endpoint lacks an authoritative
	// terminator. When maxResults <= 0 (no cap), hasMore stays false.
	hasMore := maxResults > 0 && len(rawUsers) == maxResults
	nextToken := ""
	if hasMore {
		nextToken = strconv.Itoa(startAt + len(rawUsers))
	}

	if idOnly {
		ids := make([]string, len(users))
		for i, u := range users {
			ids[i] = u.AccountID
		}
		return atkpresent.EmitIDsWithPaginationToken(opts, ids, hasMore, nextToken)
	}

	if len(users) == 0 {
		model := atkpresent.UserPresenter{}.PresentEmpty(query)
		model.Sections = atkpresent.AppendPaginationHintWithToken(model.Sections, hasMore, nextToken)
		return atkpresent.Emit(opts, model)
	}

	model := atkpresent.UserPresenter{}.PresentUserListWithPagination(users, opts.IsExtended(), hasMore, nextToken)
	if projected {
		projection.ApplyToTableInModel(model, selected)
	}
	return atkpresent.Emit(opts, model)
}

func filterSearchUsers(users []api.User) []api.User {
	filtered := make([]api.User, 0, len(users))
	for _, user := range users {
		if user.Active && user.AccountType == "atlassian" {
			filtered = append(filtered, user)
		}
	}
	return filtered
}
