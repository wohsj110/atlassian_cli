package space

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/wohsj110/atlassian_cli/tools/atk-cfl/internal/cmd/root"
	atkpresent "github.com/wohsj110/atlassian_cli/tools/atk-cfl/internal/present"
)

type viewOptions struct {
	*root.Options
}

func newViewCmd(rootOpts *root.Options) *cobra.Command {
	opts := &viewOptions{Options: rootOpts}

	cmd := &cobra.Command{
		Use:     "view <space-key>",
		Aliases: []string{"get"},
		Short:   "View space details",
		Long:    `View details of a Confluence space by its key.`,
		Example: `  # View a space
  atk-cfl space view DEV`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runView(cmd.Context(), args[0], opts)
		},
	}

	return cmd
}

func runView(ctx context.Context, spaceKey string, opts *viewOptions) error {
	client, err := opts.APIClient()
	if err != nil {
		return err
	}

	space, err := client.GetSpaceByKey(ctx, spaceKey)
	if err != nil {
		return err
	}

	return atkpresent.Emit(opts.Options, atkpresent.SpacePresenter{}.PresentDetail(space, opts.Full))
}
