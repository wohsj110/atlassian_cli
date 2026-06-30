package issues

import (
	"context"
	"strings"

	"github.com/spf13/cobra"

	"github.com/wohsj110/atlassian_cli/tools/atk-jira/api"
	"github.com/wohsj110/atlassian_cli/tools/atk-jira/internal/cache"
	"github.com/wohsj110/atlassian_cli/tools/atk-jira/internal/cmd/root"
	atkpresent "github.com/wohsj110/atlassian_cli/tools/atk-jira/internal/present"
)

func newFieldsCmd(opts *root.Options) *cobra.Command {
	var customOnly bool

	cmd := &cobra.Command{
		Use:   "fields [issue-key]",
		Short: "List available fields",
		Long:  "List fields and their metadata. If an issue key is provided, shows all fields with their current values for that issue.",
		Example: `  # List all fields
  atk-jira issues fields

  # List only custom fields
  atk-jira issues fields --custom-fields

  # Show field values for a specific issue
  atk-jira issues fields PROJ-123

  # Show only custom field values for an issue
  atk-jira issues fields PROJ-123 --custom-fields

  # Extended output with searchable/navigable/orderable/clause names
  atk-jira issues fields --extended

  # Emit only field IDs
  atk-jira issues fields --id`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			issueKey := ""
			if len(args) > 0 {
				issueKey = args[0]
			}
			return runFields(cmd.Context(), opts, issueKey, customOnly)
		},
	}

	cmd.Flags().BoolVar(&customOnly, "custom-fields", false, "Show only custom fields")
	cmd.Flags().BoolVar(&customOnly, "custom", false, "Show only custom fields")
	_ = cmd.Flags().MarkHidden("custom")

	return cmd
}

func runFields(ctx context.Context, opts *root.Options, issueKey string, customOnly bool) error {
	client, err := opts.APIClient()
	if err != nil {
		return err
	}

	if issueKey != "" {
		return runIssueFields(ctx, opts, client, issueKey, customOnly)
	}
	return runGlobalFields(ctx, opts, client, customOnly)
}

func runGlobalFields(ctx context.Context, opts *root.Options, client *api.Client, customOnly bool) error {
	fields, err := cache.GetFieldsCacheFirst(ctx, client)
	if err != nil {
		return err
	}
	if customOnly {
		var custom []api.Field
		for _, f := range fields {
			if f.Custom {
				custom = append(custom, f)
			}
		}
		fields = custom
	}

	if opts.EmitIDOnly() {
		ids := make([]string, len(fields))
		for i, f := range fields {
			ids[i] = f.ID
		}
		return atkpresent.EmitIDs(opts, ids)
	}

	if len(fields) == 0 {
		return atkpresent.Emit(opts, atkpresent.FieldPresenter{}.PresentEmpty())
	}

	model := atkpresent.FieldPresenter{}.PresentList(fields, opts.IsExtended())
	return atkpresent.Emit(opts, model)
}

func runIssueFields(ctx context.Context, opts *root.Options, client *api.Client, issueKey string, customOnly bool) error {
	issue, err := client.GetIssue(ctx, issueKey)
	if err != nil {
		return err
	}

	fields, err := cache.GetFieldsCacheFirst(ctx, client)
	if err != nil {
		return err
	}

	entries := api.ExtractIssueFieldValues(issue, fields)

	if customOnly {
		var filtered []api.IssueFieldEntry
		for _, e := range entries {
			if strings.HasPrefix(e.ID, "customfield_") {
				filtered = append(filtered, e)
			}
		}
		entries = filtered
	}

	if len(entries) == 0 {
		return atkpresent.Emit(opts, atkpresent.FieldPresenter{}.PresentEmpty())
	}

	if opts.EmitIDOnly() {
		ids := make([]string, len(entries))
		for i, e := range entries {
			ids[i] = e.ID
		}
		return atkpresent.EmitIDs(opts, ids)
	}

	model := atkpresent.FieldPresenter{}.PresentIssueFields(entries)
	return atkpresent.Emit(opts, model)
}
