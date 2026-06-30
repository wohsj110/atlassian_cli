package automation

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/wohsj110/atlassian_cli/shared/present"
	"github.com/wohsj110/atlassian_cli/shared/prompt"

	"github.com/wohsj110/atlassian_cli/tools/atk-jira/internal/cmd/root"
	atkpresent "github.com/wohsj110/atlassian_cli/tools/atk-jira/internal/present"
)

func newDeleteCmd(opts *root.Options) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "delete <rule-id>",
		Short: "Delete an automation rule",
		Long: `Delete an automation rule permanently. If the rule is currently ENABLED,
it will be automatically disabled before deletion.

This action cannot be undone.`,
		Example: `  # Delete a rule (will prompt for confirmation)
  atk-jira auto delete 019cd438-229b-75f4-a443-9a96e687b867

  # Delete without confirmation
  atk-jira auto delete 019cd438-229b-75f4-a443-9a96e687b867 --force`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDelete(cmd.Context(), opts, args[0], force)
		},
	}

	cmd.Flags().BoolVar(&force, "force", false, "Skip confirmation prompt")

	return cmd
}

func runDelete(ctx context.Context, opts *root.Options, ruleID string, force bool) error {
	// §3.4: short-circuit BEFORE any API call so --non-interactive without
	// --force returns ErrConfirmationRequired even if the API lookup
	// would have failed first (auth/not-found/network).
	if opts.NonInteractive && !force {
		return prompt.ErrConfirmationRequired
	}

	client, err := opts.APIClient()
	if err != nil {
		return err
	}

	current, err := client.GetAutomationRule(ctx, ruleID)
	if err != nil {
		return err
	}

	// Defense-in-depth: the early --non-interactive short-circuit above
	// would have returned by now, but pinning the gate on both `force`
	// and `NonInteractive` keeps the policy consistent with issues/page/
	// attachment delete so a future refactor that moves the short-circuit
	// can't leak warning text to stderr under --non-interactive.
	if !force && !opts.NonInteractive {
		fmt.Fprintf(opts.Stderr, "This will permanently delete rule %q (%s). This action cannot be undone.\n", current.Name, ruleID)
		fmt.Fprint(opts.Stderr, "Are you sure? [y/N]: ")
	}
	confirmed, err := prompt.ConfirmOrFail(force, opts.NonInteractive, opts.Stdin)
	if err != nil {
		return err
	}
	if !confirmed {
		model := atkpresent.AutomationPresenter{}.PresentDeleteCancelled()
		out := present.Render(model, opts.RenderStyle())
		fmt.Fprint(opts.Stdout, out.Stdout)
		return nil
	}

	// API rejects DELETE on ENABLED rules — disable first.
	wasEnabled := current.State == "ENABLED"
	if wasEnabled {
		if err := client.SetAutomationRuleState(ctx, ruleID, false); err != nil {
			return err
		}
	}

	if err := client.DeleteAutomationRule(ctx, ruleID); err != nil {
		if wasEnabled {
			return fmt.Errorf("rule was disabled but delete failed: %w — re-enable with: atk-jira auto enable %s", err, ruleID)
		}
		return err
	}

	model := atkpresent.AutomationPresenter{}.PresentDeleted(ruleID)
	out := present.Render(model, opts.RenderStyle())
	fmt.Fprint(opts.Stdout, out.Stdout)
	return nil
}
