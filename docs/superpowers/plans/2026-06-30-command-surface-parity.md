# Command Surface Parity Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Reimplement the `documented command surface` command surface and testing strategy under `atk-jira` and `atk-cfl`.

**Architecture:** Keep the existing monorepo shape, but grow it toward shared auth/config/client/error/output/adf packages plus tool-specific API, presenter, and command packages. Use documented behavior as the compatibility target, not as copied source.

**Tech Stack:** Go standard library first, existing project packages, mocked HTTP servers for unit tests, optional real Atlassian integration tests behind environment gates.

---

## File Map

- Modify: `AGENTS.md` to state command-surface parity as the product target.
- Modify: `README.md` and `README.zh-CN.md` to point to the compatibility spec.
- Modify: `skills/Jira/CliReference.md` and `skills/Confluence/CliReference.md` after each command family is implemented.
- Modify: `npm/skill-installer/skills/**` in lockstep with `skills/**`.
- Create: `docs/COMPATIBILITY.md` for the target command surface.
- Create: `docs/TEST_STRATEGY.md` for parity testing.
- Create: `NOTICE.md` for MIT attribution.
- Modify: `shared/auth`, `shared/config`, `shared/client`, `shared/errors`, `shared/output`.
- Create: `shared/adf` once comment/description/page body conversion grows beyond plain text.
- Modify: `tools/atk-jira/internal/jira`, `tools/atk-jira/internal/cmd`, `tools/atk-jira/internal/view`.
- Modify: `tools/atk-cfl/internal/confluence`, `tools/atk-cfl/internal/cmd`, `tools/atk-cfl/internal/view`.

## Phase 1: Shared Foundation

- [ ] Write failing tests for config precedence: `JIRA_*` and `CFL_*` override `ATLASSIAN_*`, which overrides config.
- [ ] Implement config precedence and auth method validation.
- [ ] Write failing tests for client `Get/Post/Put/Delete`, error parsing, redacted verbose output, and body truncation.
- [ ] Implement shared client parity.
- [ ] Write failing tests for API error parsing across Jira and Confluence shapes.
- [ ] Implement sentinel errors and helpers.
- [ ] Write failing tests for `shared/adf` plain text to ADF.
- [ ] Implement minimal ADF package.

## Phase 2: Jira Daily Surface

- [ ] Write parser tests for `atk-jira issues list/get/search/create/update/assign/delete/types/fields/field-options`.
- [ ] Implement command dispatch with contract-compatible plural `issues` resource while keeping current singular aliases.
- [ ] Write API tests for issue get/search/create/update using mocked HTTP.
- [ ] Implement Jira issue API methods.
- [ ] Write output tests for default, `--id`, `--extended`, `--fulltext`, and `--fields`.
- [ ] Implement presenter/view split for Jira issues.
- [ ] Write parser/API/output tests for `atk-jira comments list/add/delete`.
- [ ] Implement comments fully.
- [ ] Write parser/API/output tests for `atk-jira transitions list/do`.
- [ ] Implement transitions.

## Phase 3: Jira Supporting Surface

- [ ] Implement `projects list/get`.
- [ ] Implement `users search/get`.
- [ ] Implement `boards list/get`.
- [ ] Implement `sprints list/current/issues/add`.
- [ ] Implement `links list/create/delete/types`.
- [ ] Implement `attachments list/add/get/delete`.
- [ ] Add admin surfaces after daily-use parity: fields, dashboards, automation, refresh.

## Phase 4: Confluence Daily Surface

- [ ] Write parser tests for `atk-cfl page list/view/create/edit/copy/delete`.
- [ ] Implement contract-compatible `page view` while keeping current `page get` alias.
- [ ] Write API tests for page list/view/create/edit.
- [ ] Implement page API methods.
- [ ] Write output tests for `-o table`, `-o plain`, `--full`, `--raw`, `--content-only`, `--no-truncate`.
- [ ] Implement Confluence presenter/view split.
- [ ] Write parser/API/output tests for top-level `atk-cfl search`.
- [ ] Implement search with CQL construction.

## Phase 5: Confluence Supporting Surface

- [ ] Implement `space list/view/create/update/delete`.
- [ ] Implement `attachment list/upload/download/delete`.
- [ ] Implement config, set-credential, and completion surfaces.

## Phase 6: Integration and Release Safety

- [ ] Add optional Jira integration tests gated by `ATK_JIRA_INTEGRATION=1`.
- [ ] Add optional Confluence integration tests gated by `ATK_CONFLUENCE_INTEGRATION=1`.
- [ ] Ensure no default CI path requires credentials.
- [ ] Ensure no command logs secrets in verbose output, errors, config output, or tests.
- [ ] Run `go test ./shared/... ./tools/atk-jira/... ./tools/atk-cfl/...`.
- [ ] Run `npm test --prefix npm/skill-installer`.
- [ ] Run local build smoke for both binaries.

## Current Manual Verification Commands

```bash
go test ./shared/... ./tools/atk-jira/... ./tools/atk-cfl/...
npm test --prefix npm/skill-installer
tmpdir=$(mktemp -d)
go build -o "$tmpdir/atk-jira" ./tools/atk-jira/cmd/atk-jira
go build -o "$tmpdir/atk-cfl" ./tools/atk-cfl/cmd/atk-cfl
"$tmpdir/atk-jira" --help
"$tmpdir/atk-cfl" --help
```

## Compatibility Notes

- Keep MIT attribution in `NOTICE.md`.
- Do not copy large source files unless the copied component is explicitly listed in `NOTICE.md`.
- Compatibility aliases can preserve current early commands while the plural the command contract surface lands.
