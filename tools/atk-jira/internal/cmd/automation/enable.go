package automation

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/wohsj110/atlassian_cli/shared/present"

	"github.com/wohsj110/atlassian_cli/tools/atk-jira/internal/cmd/root"
	"github.com/wohsj110/atlassian_cli/tools/atk-jira/internal/mutation"
	atkpresent "github.com/wohsj110/atlassian_cli/tools/atk-jira/internal/present"
)

func newEnableCmd(opts *root.Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "enable <rule-id>",
		Short: "Enable an automation rule",
		Long:  "Enable a disabled automation rule. This is a safe operation that does not modify the rule definition.",
		Example: `  atk-jira automation enable 12345
  atk-jira auto enable 12345`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSetState(cmd.Context(), opts, args[0], true)
		},
	}

	return cmd
}

func runSetState(ctx context.Context, opts *root.Options, ruleID string, enabled bool) error {
	client, err := opts.APIClient()
	if err != nil {
		return err
	}

	current, err := client.GetAutomationRule(ctx, ruleID)
	if err != nil {
		return err
	}

	newState := "DISABLED"
	if enabled {
		newState = "ENABLED"
	}

	if current.State == newState {
		return atkpresent.Emit(opts, atkpresent.AutomationPresenter{}.PresentNoChange(current.Name, newState))
	}

	if err := client.SetAutomationRuleState(ctx, ruleID, enabled); err != nil {
		return err
	}

	if opts.EmitIDOnly() {
		return atkpresent.EmitIDs(opts, []string{ruleID})
	}

	savedName := current.Name
	savedState := current.State
	return mutation.WriteAndPresent(ctx, opts, mutation.Config{
		Write: func(_ context.Context) (string, error) {
			return ruleID, nil
		},
		Fetch: func(ctx context.Context, id string) (*present.OutputModel, error) {
			rule, err := client.GetAutomationRule(ctx, id)
			if err != nil {
				return nil, err
			}
			return atkpresent.AutomationPresenter{}.PresentDetail(rule, false), nil
		},
		IsFresh: func(model *present.OutputModel) bool {
			return mutation.DetailFieldEquals(model, "State", newState)
		},
		Fallback: func(_ string) *present.OutputModel {
			return atkpresent.AutomationPresenter{}.PresentStateChanged(savedName, savedState, newState)
		},
	})
}
