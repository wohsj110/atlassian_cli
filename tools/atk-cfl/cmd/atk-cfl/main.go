// Package main is the entry point for the atk-cfl Confluence CLI.
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

	"github.com/wohsj110/atlassian_cli/tools/atk-cfl/internal/cmd/attachment"
	"github.com/wohsj110/atlassian_cli/tools/atk-cfl/internal/cmd/completion"
	"github.com/wohsj110/atlassian_cli/tools/atk-cfl/internal/cmd/configcmd"
	initcmd "github.com/wohsj110/atlassian_cli/tools/atk-cfl/internal/cmd/init"
	"github.com/wohsj110/atlassian_cli/tools/atk-cfl/internal/cmd/me"
	"github.com/wohsj110/atlassian_cli/tools/atk-cfl/internal/cmd/page"
	"github.com/wohsj110/atlassian_cli/tools/atk-cfl/internal/cmd/root"
	"github.com/wohsj110/atlassian_cli/tools/atk-cfl/internal/cmd/search"
	"github.com/wohsj110/atlassian_cli/tools/atk-cfl/internal/cmd/setcredential"
	"github.com/wohsj110/atlassian_cli/tools/atk-cfl/internal/cmd/space"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cmd, opts := root.NewCmd()

	root.RegisterCommands(cmd, opts,
		initcmd.Register,
		configcmd.Register,
		me.Register,
		page.Register,
		space.Register,
		attachment.Register,
		search.Register,
		setcredential.Register,
		completion.Register,
	)

	err := cmd.ExecuteContext(ctx)
	// Emit the one-time §1.8 migration notice (if migration ran this
	// invocation) before exiting — flushed here, not in a defer, so it
	// still prints when a later command error triggers os.Exit.
	keyring.FlushMigrationNotice(os.Stderr)
	if err != nil {
		// set-credential --json may have already emitted its envelope on
		// stdout; in that case stderr stays empty per §1.5.2.
		if !errors.Is(err, keyring.ErrSetCredentialEnvelopeEmitted) {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		}
		os.Exit(exitcode.GeneralError)
	}
}
