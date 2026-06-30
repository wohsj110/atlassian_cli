package space

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/wohsj110/atlassian_cli/tools/atk-cfl/api"
	"github.com/wohsj110/atlassian_cli/tools/atk-cfl/internal/cmd/root"
	atkpresent "github.com/wohsj110/atlassian_cli/tools/atk-cfl/internal/present"
)

type listOptions struct {
	*root.Options
	limit     int
	spaceType string
	cursor    string
}

func newListCmd(rootOpts *root.Options) *cobra.Command {
	opts := &listOptions{Options: rootOpts}

	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List Confluence spaces",
		Long:    `List all Confluence spaces you have access to.`,
		Example: `  # List all spaces
  atk-cfl space list

  # List only global spaces
  atk-cfl space list --type global

  # Paginate through results
  atk-cfl space list --cursor "eyJpZCI6MTIzfQ=="`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runList(cmd.Context(), opts)
		},
	}

	cmd.Flags().IntVarP(&opts.limit, "limit", "l", 25, "Maximum number of spaces to return")
	cmd.Flags().StringVarP(&opts.spaceType, "type", "t", "", "Filter by space type (global, personal)")
	cmd.Flags().StringVar(&opts.cursor, "cursor", "", "Pagination cursor for next page")

	return cmd
}

func runList(ctx context.Context, opts *listOptions) error {
	if opts.limit < 0 {
		return fmt.Errorf("invalid limit: %d (must be >= 0)", opts.limit)
	}

	if opts.limit == 0 {
		return atkpresent.Emit(opts.Options, atkpresent.SpacePresenter{}.PresentEmpty())
	}

	client, err := opts.APIClient()
	if err != nil {
		return err
	}

	apiOpts := &api.ListSpacesOptions{
		Limit:  opts.limit,
		Type:   opts.spaceType,
		Cursor: opts.cursor,
	}

	result, err := client.ListSpaces(ctx, apiOpts)
	if err != nil {
		return fmt.Errorf("listing spaces: %w", err)
	}

	if len(result.Results) == 0 {
		return atkpresent.Emit(opts.Options, atkpresent.SpacePresenter{}.PresentEmpty())
	}
	return atkpresent.Emit(opts.Options, atkpresent.SpacePresenter{}.PresentList(result.Results, opts.Full, atkpresent.ExtractCursor(result.Links.Next), result.HasMore()))
}
