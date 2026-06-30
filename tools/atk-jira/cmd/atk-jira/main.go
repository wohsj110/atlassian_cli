// Package main is the entry point for the atk-jira CLI.
//
// Distribution is fully automated: merges to main with feat:/fix: prefixes
// trigger auto-release, which runs GoReleaser (Homebrew + binary artifacts)
// and dispatches the chocolatey and winget publish workflows.
package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/wohsj110/atlassian_cli/shared/exitcode"
	"github.com/wohsj110/atlassian_cli/shared/keyring"

	"github.com/wohsj110/atlassian_cli/tools/atk-jira/internal/cmd/attachments"
	"github.com/wohsj110/atlassian_cli/tools/atk-jira/internal/cmd/automation"
	"github.com/wohsj110/atlassian_cli/tools/atk-jira/internal/cmd/boards"
	"github.com/wohsj110/atlassian_cli/tools/atk-jira/internal/cmd/comments"
	"github.com/wohsj110/atlassian_cli/tools/atk-jira/internal/cmd/completion"
	"github.com/wohsj110/atlassian_cli/tools/atk-jira/internal/cmd/configcmd"
	"github.com/wohsj110/atlassian_cli/tools/atk-jira/internal/cmd/dashboards"
	"github.com/wohsj110/atlassian_cli/tools/atk-jira/internal/cmd/fields"
	"github.com/wohsj110/atlassian_cli/tools/atk-jira/internal/cmd/initcmd"
	"github.com/wohsj110/atlassian_cli/tools/atk-jira/internal/cmd/issues"
	"github.com/wohsj110/atlassian_cli/tools/atk-jira/internal/cmd/links"
	"github.com/wohsj110/atlassian_cli/tools/atk-jira/internal/cmd/me"
	"github.com/wohsj110/atlassian_cli/tools/atk-jira/internal/cmd/projects"
	"github.com/wohsj110/atlassian_cli/tools/atk-jira/internal/cmd/refresh"
	"github.com/wohsj110/atlassian_cli/tools/atk-jira/internal/cmd/remotelinks"
	"github.com/wohsj110/atlassian_cli/tools/atk-jira/internal/cmd/root"
	"github.com/wohsj110/atlassian_cli/tools/atk-jira/internal/cmd/setcredential"
	"github.com/wohsj110/atlassian_cli/tools/atk-jira/internal/cmd/sprints"
	"github.com/wohsj110/atlassian_cli/tools/atk-jira/internal/cmd/transitions"
	"github.com/wohsj110/atlassian_cli/tools/atk-jira/internal/cmd/users"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	err := run(ctx)
	// Emit the one-time §1.8 migration notice (if migration ran this
	// invocation) before exiting — flushed here, not in a defer, so it
	// still prints when a command error triggers os.Exit.
	keyring.FlushMigrationNotice(os.Stderr)
	if err != nil {
		// set-credential --json may have already emitted its envelope on
		// stdout; in that case stderr stays empty per §1.5.2.
		if !errors.Is(err, root.ErrAlreadyReported) && !errors.Is(err, keyring.ErrSetCredentialEnvelopeEmitted) {
			fmt.Fprintln(os.Stderr, err)
		}
		os.Exit(exitcode.GeneralError)
	}
}

func run(ctx context.Context) error {
	rootCmd, opts := root.NewCmd()

	// Register all commands
	initcmd.Register(rootCmd, opts)
	configcmd.Register(rootCmd, opts)
	fields.Register(rootCmd, opts)
	issues.Register(rootCmd, opts)
	transitions.Register(rootCmd, opts)
	comments.Register(rootCmd, opts)
	links.Register(rootCmd, opts)
	remotelinks.Register(rootCmd, opts)
	attachments.Register(rootCmd, opts)
	automation.Register(rootCmd, opts)
	boards.Register(rootCmd, opts)
	dashboards.Register(rootCmd, opts)
	projects.Register(rootCmd, opts)
	sprints.Register(rootCmd, opts)
	users.Register(rootCmd, opts)
	me.Register(rootCmd, opts)
	refresh.Register(rootCmd, opts)
	setcredential.Register(rootCmd, opts)
	completion.Register(rootCmd, opts)

	return rootCmd.ExecuteContext(ctx)
}
