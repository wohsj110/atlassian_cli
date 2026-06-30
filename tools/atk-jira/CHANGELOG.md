# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Changed

- **Breaking:** Short alias for `--file` renamed from `-f` to `-F` on `attachments add`, `automation create`, and `automation update`. `-f` continues to mean `--field` on field-setting commands (`issues create`/`update`, `transitions do`). No back-compat alias.
- **Breaking:** `--output` / `-o` flag removed entirely. Use `--id` for identifier-only output, `--extended` for admin/schema/audit detail, `--fulltext` to disable truncation. `automation export` is the only command that still emits JSON (writes directly to stdout, bypasses the global flag system). ([#216](https://github.com/wohsj110/atlassian_cli/issues/216), [#332](https://github.com/wohsj110/atlassian_cli/pull/332))
- Default page size for paginated commands converged to 50: `issues list` and `issues search` were 25, `users search` was 10. `users search` and `dashboards list` also gain the `-m` short alias for `--max`.
- Global output flags replaced with `--extended`, `--fulltext`, `--id` per the new output model. The `--full` flag is removed; use `--extended` for admin/schema/audit detail. `--id` takes precedence over `--extended` and `--fulltext`. Per-command `--no-truncate` flags remain as deprecated aliases for `--fulltext`. ([#231](https://github.com/wohsj110/atlassian_cli/issues/231), [#230](https://github.com/wohsj110/atlassian_cli/issues/230))
- `links types`, `issues types`, `boards list`, `sprints list`, and `users get` now serve from the local instance cache by default — removes per-command API calls in the most common paths. Run `atk-jira refresh` to update. ([#328](https://github.com/wohsj110/atlassian_cli/pull/328), [#329](https://github.com/wohsj110/atlassian_cli/pull/329), [#330](https://github.com/wohsj110/atlassian_cli/pull/330))

### Added

- `issues history <issue-key>` lists Jira changelog history with compact changed-field rows, `--id`, `--extended`, `--fields`, and offset pagination support.
- `issues get` now accepts multiple issue keys and renders a summary table for the batch.
- `issues check <issue-key>` subcommand to audit an issue for populated/missing field values, with `--require` (hard-fail) and `--warn` (advisory) flags. A curated default warn-list (Summary, Description, Assignee, Priority, Labels, Story Points, Sprint, Components, Fix Version/s) applies when no flags are passed. Useful as a transition guardrail or CI step.
- `users get <account-id>` subcommand to look up a user by account ID
- `--assignee none` on `issues update` and `--field assignee=null` to unassign issues
- `--field` flag accumulates repeated values for the same key, enabling multi-checkbox and multi-select custom fields
- `--fields` flag on `issues list` and `issues search` for explicit field selection
- Auto-pagination for `issues list` and `issues search` — `--max` returns up to N results across pages
- Automation rule builder module for constructing rule JSON programmatically with a fluent Go API
- Service account support with bearer auth (`--auth-method bearer`) for scoped API tokens
- Cursor-based pagination with `--next-page-token` and lightweight fields (`--full`) for `issues list` and `issues search`
- `dashboards` command group: `list`, `get`, `create`, `delete`, `gadgets list`, `gadgets remove`
- `links` command group for issue links: `list`, `create`, `delete`, `types`
- `--type` flag on `issues update` to change issue type via the bulk move API
- `\n`, `\t`, `\\` escape sequence handling in `--description` flag
- `--assignee` flag on `issues create` and `issues update` (accepts account ID, email, or `"me"`)
- `--parent` flag on `issues create` and `issues update`
- `fields` command group for custom field management: `create`, `delete` (trash), `restore`, `contexts` (list/create/delete), and `options` (list/add/update/delete)
- `projects create`, `update`, `delete`, `restore`, `types` commands for full project management
- `automation create` command to create rules from JSON files
- `automation enable`, `disable`, `update`, `export` commands for full automation rule management
- `--full` flag on `issues get` and `comments list` to show full content without truncation
- `init` command for guided setup wizard
- `issues move` command to move issues between projects
- `attachments` commands: list, add, get, delete
- Wiki markup detection and automatic conversion to ADF
- `issues field-options` command to list allowed values for select fields ([#36](https://github.com/wohsj110/atlassian_cli/tools/atk-jira/pull/36))
- `issues types` command to list valid issue types per project ([#22](https://github.com/wohsj110/atlassian_cli/tools/atk-jira/pull/22))
- `users search` command for finding account IDs by name/email ([#34](https://github.com/wohsj110/atlassian_cli/tools/atk-jira/pull/34))
- Show required fields for transitions in `transitions list --fields` ([#35](https://github.com/wohsj110/atlassian_cli/tools/atk-jira/pull/35))
- Include custom fields in issue JSON output ([#37](https://github.com/wohsj110/atlassian_cli/tools/atk-jira/pull/37))

### Changed

- Consolidated markdown-to-ADF conversion into shared package
- Improved init/config UX with huh forms and --force flag on clear
- **Binary renamed to `atk-jira`** - The CLI binary is now `atk-jira` (short for atk-jira). Install via `brew install atk-jira`, run with `atk-jira`. ([#41](https://github.com/wohsj110/atlassian_cli/tools/atk-jira/pull/41))
- Module path migrated to `github.com/wohsj110/atlassian_cli/tools/atk-jira` ([#39](https://github.com/wohsj110/atlassian_cli/tools/atk-jira/pull/39))

### Fixed

- `atk-jira issues create --field components=<id-or-name>` and `--field fixVersions=<id-or-name>` now work. Previously the array formatter only handled multi-checkbox (`option` items) and fell through to a plain string array for component and version items, which Jira rejects with `The list contains an invalid value`. Multi-value via repeated `--field` accumulates as expected. Thanks to @romiguelangel for the fix.
- `--field "key = value"` (whitespace around `=`) now parses correctly.
- `issues move` and `issues update --type` no longer rely on string-matching the API error message to detect 404s — uses the structured error code instead, eliminating false negatives if Jira reword the message.
- `--verbose` now logs the outbound request JSON body and any 4xx/5xx response body (each truncated at 4 KB), surfacing field-level Jira errors that previously appeared only as opaque codes like `INVALID_INPUT`.
- Empty fenced/indented code blocks and empty table cells no longer produce invalid ADF text nodes with empty content.
- `atk-jira issues move --no-wait` and `--no-notify` now parse correctly. Previously the help text mentioned them but pflag did not register the negations, so they failed with "unknown flag".
- `\n`, `\t`, `\\` escape sequences now work in `comments add --body`
- `issues search` and `issues list` with `-o json` now return all fields including custom fields by default
- Wiki markup conversion no longer mangles hyphens and tildes
- `--field parent=PROJ-123` and issuelink-type custom fields now format correctly instead of producing a `"data was not an object"` API error
- `config show -o json` no longer appends trailing plain text after JSON body
- `projects create` success message uses the input name instead of the empty API response name
- `ProjectDetail.ID` uses `json.Number` to handle numeric API responses
- Automation rule state endpoint uses correct payload format for Jira Cloud
- `automation create` strips server-assigned fields and parses `ruleUuid` correctly
- `--field` flag handles structured fields (e.g., `priority=High`) in create and update
- Validate file input before making network calls
- Automation API parsing aligned with Jira Cloud response format
- Show user display name instead of account ID in assign command output ([#33](https://github.com/wohsj110/atlassian_cli/tools/atk-jira/pull/33))
- Convert number and textarea fields to correct API format when updating issues ([#32](https://github.com/wohsj110/atlassian_cli/tools/atk-jira/pull/32))
