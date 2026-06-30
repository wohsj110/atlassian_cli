package attachment

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/wohsj110/atlassian_cli/shared/prompt"

	"github.com/wohsj110/atlassian_cli/tools/atk-cfl/internal/cmd/root"
	atkpresent "github.com/wohsj110/atlassian_cli/tools/atk-cfl/internal/present"
)

type deleteOptions struct {
	*root.Options
	force bool
}

func newDeleteCmd(rootOpts *root.Options) *cobra.Command {
	opts := &deleteOptions{Options: rootOpts}

	cmd := &cobra.Command{
		Use:   "delete <attachment-id>",
		Short: "Delete an attachment",
		Long:  `Delete an attachment by its ID.`,
		Example: `  # Delete an attachment
  atk-cfl attachment delete att123

  # Delete without confirmation
  atk-cfl attachment delete att123 --force`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDeleteAttachment(cmd.Context(), args[0], opts)
		},
	}

	cmd.Flags().BoolVarP(&opts.force, "force", "f", false, "Skip confirmation prompt")

	return cmd
}

func runDeleteAttachment(ctx context.Context, attachmentID string, opts *deleteOptions) error {
	// §3.4: short-circuit BEFORE any API call so --non-interactive without
	// --force returns ErrConfirmationRequired even if the attachment
	// lookup would have failed first (auth/not-found/network).
	if opts.NonInteractive && !opts.force {
		return prompt.ErrConfirmationRequired
	}

	client, err := opts.APIClient()
	if err != nil {
		return err
	}

	attachment, err := client.GetAttachment(ctx, attachmentID)
	if err != nil {
		return fmt.Errorf("getting attachment: %w", err)
	}

	if !opts.force && !opts.NonInteractive {
		_, _ = fmt.Fprintf(opts.Stderr, "About to delete attachment: %s (ID: %s)\n", attachment.Title, attachment.ID)
		_, _ = fmt.Fprint(opts.Stderr, "Are you sure? [y/N]: ")
	}
	confirmed, err := prompt.ConfirmOrFail(opts.force, opts.NonInteractive, opts.Stdin)
	if err != nil {
		return err
	}
	if !confirmed {
		return atkpresent.Emit(opts.Options, atkpresent.PresentDeletionCancelled())
	}

	if err := client.DeleteAttachment(ctx, attachmentID); err != nil {
		return fmt.Errorf("deleting attachment: %w", err)
	}

	return atkpresent.Emit(opts.Options, atkpresent.AttachmentPresenter{}.PresentDelete(attachment))
}
