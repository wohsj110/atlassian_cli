# atk-jira Development Guide

This is the repo-local source for working on `atk-jira`, the Jira CLI inside the atlassian-cli monorepo.

## Repository

Binary: `atk-jira`
Module: `github.com/wohsj110/atlassian_cli/tools/atk-jira`
Entrypoint: `cmd/atk-jira/main.go`
Shared module replacement: `github.com/wohsj110/atlassian_cli/shared => ../../shared`

## Repo-Local Sources

### Monorepo Guide

Local source of truth: docs/development.md
Local convenience copy, if present: `../../../docs/development.md`

### Artifact Contract

Local source of truth: docs/ARTIFACT_CONTRACT.md
Local convenience copy, if present: `../../../docs/ARTIFACT_CONTRACT.md`

### Command Surface Guardrails

Local source of truth: tools/atk-jira/internal/cmd/GUARDRAILS.md
Local convenience copy, if present: `../internal/cmd/GUARDRAILS.md`

### Output Specification

Local source of truth: tools/atk-jira/internal/cmd/OUTPUT_SPEC.md
Local convenience copy, if present: `../internal/cmd/OUTPUT_SPEC.md`

## Shared Sources

### Shared Open CLI Standards

Source of truth: https://github.com/wohsj110/cli-common/tree/main/docs
Local convenience copy, if present: `../../../../cli-common/docs`

### Shared Automation

Source of truth: https://github.com/wohsj110/.github
Local convenience copy, if present: `../../../../.github`

## Quick Commands

```bash
make build
make test
make test-cover
make lint
make tidy
make check
make install
```

`make build` writes `bin/atk-jira` from `./cmd/atk-jira`. `make test` runs the test suite with the race detector. `make check` runs tidy, lint, test, and build for the tool module.

## Architecture

`api/` is the importable Jira API client package. `internal/cmd/` contains Cobra command packages by resource. `internal/present` owns output rendering by resource. `internal/config` owns configuration loading and migration. `internal/version` receives build-time version data through ldflags.

Commands use a root `Options` struct for global flags and dependencies. Command packages expose `Register(rootCmd *cobra.Command, opts *root.Options)` and add subcommands through unexported factories.

## Command and Output Contracts

Read `internal/cmd/GUARDRAILS.md` before changing command names, flags, arguments, pagination, mutation safety, or cache-backed identity resolution.

Read `internal/cmd/OUTPUT_SPEC.md` before changing default output, `--id`, `--extended`, `--fulltext`, or mutation output.

`atk-jira` is text-first. It has no global output-format flag. JSON is reserved for round-trip payloads such as automation export and control-plane envelopes documented by the shared standards.

## Auth and Config

`atk-jira` participates in the shared Atlassian credential/config model described by the monorepo guide. `ATLASSIAN_*` variables apply across both tools; `JIRA_*` variables override for atk-jira. The atk-jira-specific config section carries non-secret defaults such as `default_project`.

Basic auth uses an instance URL plus email and token. Bearer auth routes through `api.atlassian.com`, requires a cloud ID, and has Atlassian platform scope limitations for some Jira surfaces.

## Testing Notes

Use `shared/testutil` for assertions. Prefer table-driven tests and `httptest.NewServer` for API client behavior. Keep tests next to implementation and use presenter-focused tests for output behavior.
