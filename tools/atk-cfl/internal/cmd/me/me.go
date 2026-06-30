// Package me provides the me command for atk-cfl.
package me

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/wohsj110/atlassian_cli/tools/atk-cfl/internal/cmd/root"
	atkpresent "github.com/wohsj110/atlassian_cli/tools/atk-cfl/internal/present"
)

// Register adds the me command to the root command.
func Register(rootCmd *cobra.Command, opts *root.Options) {
	rootCmd.AddCommand(newMeCmd(opts))
}

func newMeCmd(opts *root.Options) *cobra.Command {
	var idOnly bool
	cmd := &cobra.Command{
		Use:   "me",
		Short: "Show the currently authenticated user",
		Long: `Show the user authenticated by the current atk-cfl configuration as a token-dense one-liner: accountId | displayName | email.

Missing fields render as "-" so the row is always exactly three pipe-delimited fields.`,
		Example: `  # Show current user
  atk-cfl me

  # Show only the account ID (for scripting)
  atk-cfl me --id`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return Run(cmd.Context(), opts, idOnly)
		},
	}
	cmd.Flags().BoolVar(&idOnly, "id", false, "Print only the account ID")
	return cmd
}

// Run fetches and renders the currently authenticated user.
func Run(ctx context.Context, opts *root.Options, idOnly bool) error {
	client, err := opts.APIClient()
	if err != nil {
		return fmt.Errorf("getting API client: %w", err)
	}
	user, err := client.GetCurrentUser(ctx)
	if err != nil {
		return fmt.Errorf("getting current user: %w", err)
	}

	presenter := atkpresent.UserPresenter{}
	if idOnly {
		return atkpresent.Emit(opts, presenter.PresentUserIDOnly(user))
	}
	return atkpresent.Emit(opts, presenter.PresentUserOneLiner(user))
}
