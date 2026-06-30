package init

import (
	"errors"
	"path/filepath"

	"github.com/wohsj110/atlassian_cli/shared/credstore"
	"github.com/wohsj110/atlassian_cli/shared/view"

	"github.com/wohsj110/atlassian_cli/tools/atk-cfl/internal/config"
)

// reconcileResult captures everything finalizeInit needs after the
// detection phase: a *Config to seed the form, the shared store the user
// already had on disk (so save preserves unrelated fields like the jtk
// section), and the legacy files the user might want to clean up.
//
// Per §2.2 (MON-5328) connection config is single-sourced from the
// shared `default` section — there is no per-tool override and therefore
// no write-target choice; finalizeInit always writes connection to
// `default`.
type reconcileResult struct {
	prefill          *config.Config
	store            *credstore.Store
	consumedLegacies []string // legacy file paths folded into the result
	// affectsSibling is true when the save will mutate connection
	// credentials the sibling tool also reads (always, now that there is
	// one shared default) AND the store already held usable creds — so
	// finalizeInit confirms before overwriting a working shared config.
	affectsSibling bool
}

// detectAndReconcile decides what to do given whatever configs already
// exist on disk. Connection config is single-sourced (§2.2): it gathers
// every connection candidate (shared default, the pre-MON-5328 shared
// per-tool sections via the migration projection, and the legacy atk-cfl/jtk
// files), runs the pure divergence detector, and FAILS LOUD if they
// disagree (naming every source + field, never a value) rather than
// precedence-picking. Aligned → the unified connection is folded into
// the shared default; per-tool non-secret defaults are preserved.
//
// Path arguments are injected so tests can point them at a tempdir.
func detectAndReconcile(
	v *view.View,
	atkCFLLegacyPath, atkJiraLegacyPath, sharedPath string,
	prefillURL, prefillEmail, prefillAuthMethod, prefillCloudID string,
) (*reconcileResult, error) {
	// §3.2 PURE pre-token relocation gate: runs BEFORE per-tool
	// divergence detection and BEFORE any mutation. A corrupt old/new
	// file or a divergent old↔new shared config fails loud here naming
	// both paths, mutating nothing — consistent with the fail-loud,
	// mutate-nothing invariant.
	rel, relErr := credstore.DetectSharedRelocation(sharedPath)
	if relErr != nil {
		if errors.Is(relErr, credstore.ErrCorruptStore) {
			// A corrupt old OR new shared file: same contract/UX as the
			// Load path below — unreadable, refuse to overwrite.
			v.Error("Shared credential store at %s is unreadable: %v", sharedPath, relErr)
			v.Error("Refusing to overwrite. Fix or remove the file, then re-run atk-cfl init.")
		} else {
			v.Error("Shared credential store relocation check failed: %v", relErr)
			v.Error("Refusing to mutate anything. Reconcile the named file(s), then re-run atk-cfl init.")
		}
		return nil, relErr
	}

	store, err := credstore.Load(sharedPath)
	if err != nil {
		v.Error("Shared credential store at %s is unreadable: %v", sharedPath, err)
		v.Error("Refusing to overwrite. Fix or remove the file, then re-run atk-cfl init.")
		return nil, err
	}
	// Migration projection retains the pre-MON-5328 per-tool connection
	// fields the canonical Store dropped (EnsureMigrated's token-only
	// scrub preserves them, so they are still readable here).
	proj, err := credstore.LoadSharedLegacyProjection(sharedPath)
	if err != nil {
		v.Error("Shared credential store at %s is unreadable: %v", sharedPath, err)
		v.Error("Refusing to overwrite. Fix or remove the file, then re-run atk-cfl init.")
		return nil, err
	}
	if proj == nil {
		proj = &credstore.SharedLegacyProjection{Path: sharedPath}
	}

	atkCFLLegacy, atkCFLErr := credstore.LoadLegacyAtkCFL(atkCFLLegacyPath)
	if atkCFLErr != nil {
		if errors.Is(atkCFLErr, credstore.ErrCorruptStore) {
			v.Error("Legacy atk-cfl config at %s is unreadable: %v", atkCFLLegacyPath, atkCFLErr)
			v.Error("Refusing to overwrite. Fix or remove the file, then re-run atk-cfl init.")
		}
		return nil, atkCFLErr
	}
	atkJiraLegacy, atkJiraErr := credstore.LoadLegacyAtkJira(atkJiraLegacyPath)
	if atkJiraErr != nil {
		// Sibling-corrupt is a warning, not a hard stop.
		v.Info("Note: sibling atk-jira config at %s is unreadable; ignoring. (%v)", atkJiraLegacyPath, atkJiraErr)
		atkJiraLegacy = nil
	}
	agentLegacy, agentErr := loadAgentLegacyConfigs(sharedPath)
	if agentErr != nil {
		return nil, agentErr
	}

	// Build the full named connection candidate set and detect
	// divergence (pure, secret-free, no IO/keyring — shared with jtk).
	candidates := credstore.ConnCandidates(sharedPath, store.Default, proj, atkCFLLegacy, atkJiraLegacy)
	for _, l := range agentLegacy {
		candidates = append(candidates, credstore.NamedConn{
			Label: "legacy atlassian-agent-cli config",
			Path:  l.Path,
			Conn:  connFromSection(l.Section()),
		})
	}
	// Old-shared is ADDITIVE to the candidate set so a relocation copy
	// is gated on the per-tool divergence check passing (no copy while a
	// divergence is pending).
	candidates = append(candidates, credstore.OldSharedConnCandidates(rel)...)
	chosen, conflicts := credstore.DetectConnDivergence(candidates)
	if len(conflicts) > 0 {
		return nil, credstore.ConnConflictError(conflicts, candidates, "atk-cfl")
	}

	// Gated apply: every conflict gate (relocation + per-tool
	// divergence) has now passed — materialize the old shared file at
	// the new path (copy-leave-old), then reload so the remainder
	// reconciles the materialized file exactly as a returning user.
	if rel.CopyNeeded {
		if aerr := credstore.ApplySharedRelocation(rel); aerr != nil {
			v.Error("Could not relocate the shared credential store: %v", aerr)
			return nil, aerr
		}
		// Reload only the canonical store: the divergence candidates
		// (incl. old-shared) were already built and checked above, so
		// `proj` is not read again — the materialized store is what the
		// fold + preserveDefaults below operate on.
		store, err = credstore.Load(sharedPath)
		if err != nil {
			v.Error("Shared credential store at %s is unreadable: %v", sharedPath, err)
			return nil, err
		}
	}

	// affectsSibling must be judged on the ORIGINAL loaded store, BEFORE
	// folding `chosen` in — otherwise a first-time migration from only a
	// legacy file looks like it is overwriting an already-usable shared
	// default and the user gets a misleading "Save will affect sibling"
	// prompt. It is true only when the store already held usable creds
	// AND the resolved connection actually DIFFERS from what is on disk:
	// re-running `atk-cfl init` without changing the connection is a no-op
	// for jtk and must not nag (the prior per-tool model only prompted on
	// an explicit reuse choice; one shared default would otherwise prompt
	// on every re-init). Pure (HasUsableConfig + value compare only): NO
	// keyring I/O in reconcile (the B3 leak-regression rule).
	origDefault := store.Default
	affectsSibling := store.HasUsableConfig(credstore.ToolAtkCFL) &&
		!credstore.ConnEqualsSection(chosen, origDefault)

	// Aligned: fold the unified connection into the shared default and
	// preserve per-tool non-secret defaults (atk-cfl's space/output, jtk's
	// project) so neither tool loses them on next read. This in-place
	// write is intentionally redundant with applyResultToStore (which
	// finalizeInit calls after the form, with the final URL-normalized
	// values) — result.store.Default is never read between the two, so
	// the transient pre-normalization state is not observable.
	store.Default = credstore.Section{
		URL:        chosen.URL,
		Email:      chosen.Email,
		AuthMethod: chosen.AuthMethod,
		CloudID:    chosen.CloudID,
	}
	consumed := preserveDefaultsAndCollect(store, atkCFLLegacy, atkJiraLegacy, agentLegacy)

	cfg := configFromConn(chosen)
	if store.AtkCFL.DefaultSpace != "" {
		cfg.DefaultSpace = store.AtkCFL.DefaultSpace
	}
	if store.AtkCFL.OutputFormat != "" {
		cfg.OutputFormat = store.AtkCFL.OutputFormat
	}
	applyFlagOverrides(cfg, prefillURL, prefillEmail, prefillAuthMethod, prefillCloudID)

	return &reconcileResult{
		prefill:          cfg,
		store:            store,
		consumedLegacies: consumed,
		affectsSibling:   affectsSibling,
	}, nil
}

func connFromSection(s credstore.Section) credstore.ConnProfile {
	return credstore.ConnProfile{
		URL:        s.URL,
		Email:      s.Email,
		AuthMethod: s.AuthMethod,
		CloudID:    s.CloudID,
	}
}

// preserveDefaultsAndCollect fills per-tool non-secret defaults from the
// legacy files, and returns the legacy file paths that contributed a
// connection (so init can offer to delete them). Legacy values only fill
// fields the shared store leaves EMPTY: a value the user already set in
// the shared store (e.g. via a prior init that changed default_space)
// must not be silently reverted to a stale legacy value just because the
// old file still exists. Shared store wins; legacy backfills absent.
func preserveDefaultsAndCollect(
	store *credstore.Store,
	atkCFLLegacy, atkJiraLegacy *credstore.LegacyCreds,
	agentLegacy []*credstore.LegacyCreds,
) []string {
	var consumed []string
	if atkCFLLegacy != nil {
		if store.AtkCFL.DefaultSpace == "" && atkCFLLegacy.DefaultSpace != "" {
			store.AtkCFL.DefaultSpace = atkCFLLegacy.DefaultSpace
		}
		if store.AtkCFL.OutputFormat == "" && atkCFLLegacy.OutputFormat != "" {
			store.AtkCFL.OutputFormat = atkCFLLegacy.OutputFormat
		}
		if legacyHasConn(atkCFLLegacy) {
			consumed = append(consumed, atkCFLLegacy.Path)
		}
	}
	if atkJiraLegacy != nil {
		if store.AtkJira.DefaultProject == "" && atkJiraLegacy.DefaultProject != "" {
			store.AtkJira.DefaultProject = atkJiraLegacy.DefaultProject
		}
		if legacyHasConn(atkJiraLegacy) {
			consumed = append(consumed, atkJiraLegacy.Path)
		}
	}
	for _, l := range agentLegacy {
		if l == nil {
			continue
		}
		if store.AtkCFL.DefaultSpace == "" && l.DefaultSpace != "" {
			store.AtkCFL.DefaultSpace = l.DefaultSpace
		}
		if store.AtkCFL.OutputFormat == "" && l.OutputFormat != "" {
			store.AtkCFL.OutputFormat = l.OutputFormat
		}
		if store.AtkJira.DefaultProject == "" && l.DefaultProject != "" {
			store.AtkJira.DefaultProject = l.DefaultProject
		}
		if legacyHasConn(l) {
			consumed = append(consumed, l.Path)
		}
	}
	return consumed
}

func loadAgentLegacyConfigs(sharedPath string) ([]*credstore.LegacyCreds, error) {
	defaultPath, err := credstore.DefaultPath()
	if err != nil || filepath.Clean(sharedPath) != filepath.Clean(defaultPath) {
		return nil, nil
	}
	var out []*credstore.LegacyCreds
	for _, tool := range []string{credstore.ToolAtkCFL, credstore.ToolAtkJira} {
		for _, path := range credstore.LegacyAgentToolPaths(tool) {
			l, err := credstore.LoadLegacyAgentTool(path, tool)
			if err != nil {
				return nil, err
			}
			if l != nil {
				out = append(out, l)
				break
			}
		}
	}
	return out, nil
}

func legacyHasConn(l *credstore.LegacyCreds) bool {
	return l.URL != "" || l.Email != "" || l.AuthMethod != "" || l.CloudID != ""
}

func configFromConn(c credstore.ConnProfile) *config.Config {
	cfg := &config.Config{
		Email:      c.Email,
		AuthMethod: c.AuthMethod,
		CloudID:    c.CloudID,
	}
	if c.URL != "" {
		cfg.URL = credstore.URLForAtkCFL(c.URL)
	}
	return cfg
}

func applyFlagOverrides(cfg *config.Config, url, email, authMethod, cloudID string) {
	if url != "" {
		cfg.URL = url
	}
	if email != "" {
		cfg.Email = email
	}
	if authMethod != "" {
		cfg.AuthMethod = authMethod
	}
	if cloudID != "" {
		cfg.CloudID = cloudID
	}
}

// applyResultToStore writes the form's final cfg into the shared default
// (connection is single-sourced — §2.2) and preserves/sets the atk-cfl
// per-tool non-secret defaults. The jtk section and jtk defaults are
// left untouched.
func applyResultToStore(store *credstore.Store, cfg *config.Config) {
	store.Default = credstore.Section{
		URL:        credstore.NormalizeBaseURL(cfg.URL),
		Email:      cfg.Email,
		AuthMethod: cfg.AuthMethod,
		CloudID:    cfg.CloudID,
	}
	if cfg.DefaultSpace != "" {
		store.AtkCFL.DefaultSpace = cfg.DefaultSpace
	}
	if cfg.OutputFormat != "" {
		store.AtkCFL.OutputFormat = cfg.OutputFormat
	}
}
