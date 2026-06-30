// Package refresh provides the `atk-jira refresh` command for updating the instance cache.
package refresh

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/wohsj110/atlassian_cli/shared/present"

	"github.com/wohsj110/atlassian_cli/tools/atk-jira/api"
	"github.com/wohsj110/atlassian_cli/tools/atk-jira/internal/cache"
	"github.com/wohsj110/atlassian_cli/tools/atk-jira/internal/cmd/root"
	"github.com/wohsj110/atlassian_cli/tools/atk-jira/internal/config"
	atkpresent "github.com/wohsj110/atlassian_cli/tools/atk-jira/internal/present"
)

// Register registers the refresh command.
func Register(parent *cobra.Command, opts *root.Options) {
	var statusOnly bool

	cmd := &cobra.Command{
		Use:   "refresh [resources...]",
		Short: "Refresh the atk-jira instance cache",
		Long: `Refresh the atk-jira instance cache — the local snapshot of fields, projects,
users, issue types, statuses, priorities, resolutions, boards, and link types.

With no arguments, refreshes every cacheable resource. With resource names,
refreshes only those (plus any declared dependencies, auto-bootstrapped in
dependency order). With --status, reports freshness without fetching.

Requires configuration (atk-jira init). Resources unavailable under the current
auth (e.g., boards under bearer auth) are silently skipped during a refresh
and reported as "unavailable" in --status.`,
		Example: `  # Refresh everything
  atk-jira refresh

  # Refresh a subset (auto-expands to include dependencies)
  atk-jira refresh statuses

  # Show freshness without fetching
  atk-jira refresh --status`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(cmd.Context(), opts, args, statusOnly)
		},
	}

	cmd.Flags().BoolVar(&statusOnly, "status", false, "Print cache freshness; no network calls")
	parent.AddCommand(cmd)
}

func run(ctx context.Context, opts *root.Options, names []string, statusOnly bool) error {
	if !config.IsConfigured() {
		return errors.New("configuration not found — run 'atk-jira init' first")
	}

	selected, err := cache.SelectWithDeps(names)
	if err != nil {
		return err
	}

	if statusOnly {
		return runStatus(opts, selected)
	}

	client, err := opts.APIClient()
	if err != nil {
		return err
	}
	return runRefresh(ctx, opts, client, selected)
}

// runStatus renders the freshness table.
func runStatus(opts *root.Options, selected []cache.Entry) error {
	// Client is only used to evaluate Entry.Available (e.g., the bearer-auth
	// gate on boards). If construction fails, the config was already validated
	// at run()'s entry so this is a programming error — surface it.
	client, err := opts.APIClient()
	if err != nil {
		return err
	}
	now := time.Now().UTC()

	rows := make([]atkpresent.StatusRow, 0, len(selected))
	for _, e := range selected {
		rows = append(rows, buildStatusRow(e, client, now))
	}

	model := atkpresent.RefreshPresenter{}.PresentStatus(rows, now)
	out := present.Render(model, opts.RenderStyle())
	_, _ = fmt.Fprint(opts.Stdout, out.Stdout)
	_, _ = fmt.Fprint(opts.Stderr, out.Stderr)
	return nil
}

func buildStatusRow(e cache.Entry, client *api.Client, now time.Time) atkpresent.StatusRow {
	if !e.IsAvailable(client) {
		return atkpresent.StatusRow{Resource: e.Name, TTL: e.TTL, Status: cache.StatusUnavailable}
	}

	env, err := cache.ReadResource[any](e.Name)
	if errors.Is(err, cache.ErrCacheMiss) {
		return atkpresent.StatusRow{Resource: e.Name, TTL: e.TTL, Status: cache.StatusUninitialized}
	}
	if err != nil {
		// Corrupt envelope on disk: report as uninitialized rather than a novel status;
		// a subsequent refresh will overwrite it.
		return atkpresent.StatusRow{Resource: e.Name, TTL: e.TTL, Status: cache.StatusUninitialized}
	}
	return atkpresent.StatusRow{
		Resource:  e.Name,
		TTL:       e.TTL,
		FetchedAt: env.FetchedAt,
		Status:    cache.Classify(env.FetchedAt, e.TTL, now),
	}
}

// runRefresh fetches each selected entry, collects per-resource results, hands
// them to the presenter, and returns ErrAlreadyReported if any resource failed.
// Unavailable entries (e.g., boards under bearer auth) are silently skipped.
func runRefresh(ctx context.Context, opts *root.Options, client *api.Client, selected []cache.Entry) error {
	results := make([]atkpresent.RefreshResult, 0, len(selected))
	failed := 0

	for _, e := range selected {
		if !e.IsAvailable(client) {
			continue
		}
		prev := previousCount(e.Name)
		count, err := e.Fetch(ctx, client)
		results = append(results, atkpresent.RefreshResult{
			Name:     e.Name,
			Count:    count,
			Previous: prev,
			Err:      err,
			At:       time.Now().UTC(),
		})
		if err != nil {
			failed++
		}
	}

	model := atkpresent.RefreshPresenter{}.PresentRefresh(results)
	out := present.Render(model, opts.RenderStyle())
	_, _ = fmt.Fprint(opts.Stdout, out.Stdout)
	_, _ = fmt.Fprint(opts.Stderr, out.Stderr)

	if failed > 0 {
		return fmt.Errorf("%w (%d resource(s))", root.ErrAlreadyReported, failed)
	}
	return nil
}

// previousCount reports the prior entry count for a resource, or -1 if unknown.
// Counts both slice-shaped envelopes (e.g. fields) and map-shaped ones
// (e.g. issuetypes = map[projectKey][]IssueType).
func previousCount(name string) int {
	env, err := cache.ReadResource[any](name)
	if err != nil {
		return -1
	}
	switch v := env.Data.(type) {
	case []any:
		return len(v)
	case map[string]any:
		total := 0
		for _, inner := range v {
			if arr, ok := inner.([]any); ok {
				total += len(arr)
			}
		}
		return total
	default:
		return -1
	}
}
