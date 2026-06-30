# Test Strategy

[简体中文](TEST_STRATEGY.zh-CN.md)

The test strategy targets this repository's documented command surface for the `atk-jira` and `atk-cfl` binaries.

## Test Principles

1. Test command behavior before implementation.
2. Use mocked HTTP servers for API behavior.
3. Keep real Atlassian integration tests behind explicit environment variables.
4. Never require real credentials in default CI.
5. Prove write commands have safety behavior:
   - dry-run does not touch network
   - destructive commands require confirmation by default
   - `--force` skips destructive confirmation
   - `--yes` is allowed only where this project intentionally keeps an additional write confirmation path
6. Validate stdout/stderr discipline:
   - data to stdout
   - diagnostics and prompts to stderr
7. Validate no secret leakage in config, errors, verbose output, and test logs.

## Test Layers

### Shared Packages

`shared/auth`

- Basic auth header encodes `email:token`.
- Bearer auth header emits `Bearer <token>`.
- Invalid auth methods fail before network calls once validation is added.

`shared/config`

- Tool-specific env overrides `ATLASSIAN_*`.
- `ATLASSIAN_*` overrides config.
- Config file loads non-secret fields.
- Token storage must move to credential store before public release.

`shared/client`

- Relative and absolute URLs resolve correctly.
- JSON bodies marshal for POST/PUT.
- `Accept` and `Content-Type` headers are set correctly.
- Auth headers are sent.
- 4xx/5xx parse Jira and Confluence error shapes.
- Verbose output redacts secrets and truncates large bodies.
- Rate-limit and retry behavior is tested once implemented.

`shared/output`

- Table rendering is deterministic.
- Plain rendering is stable.
- JSON rendering is reserved for intentional JSON contracts.
- Empty values render as `-` where the compatibility contract requires it.

`shared/adf`

- Plain text to ADF document.
- Markdown to ADF document once markdown support is added.
- ADF to plain text for comments/descriptions/page bodies.
- Empty input behavior is deterministic.

### Jira Command Tests

Command parsing:

- `atk-jira --help` lists top-level resources.
- Every resource help lists documented actions.
- Global flags parse before and after positional arguments where the command contract permits it.
- Invalid flags fail before network calls.

Reads:

- `me`
- `issues list/get/search`
- `comments list`
- `transitions list`
- `projects list/get`
- `users search/get`

Writes:

- `comments add` posts ADF body.
- `issues create/update/assign/transition` post the expected payloads.
- destructive commands prompt unless `--force` is present.
- mutations print post-state or identifier according to flags.

API compatibility:

- Jira search uses `/rest/api/3/search/jql`.
- Comments use `/rest/api/3/issue/{key}/comment`.
- Issue get uses `/rest/api/3/issue/{key}`.
- Transition uses Jira transition APIs and reads available transitions before doing a transition where needed.

### Confluence Command Tests

Command parsing:

- `atk-cfl --help` lists top-level resources.
- `page view`, `search`, `space`, `attachment`, `config`, `set-credential`, and `completion` command surfaces match the compatibility spec.

Reads:

- `me`
- `page list/view`
- `search`
- `space list/view`
- `attachment list`

Writes:

- `page create/edit/copy/delete`
- `space create/update/delete`
- `attachment upload/delete`

Content behavior:

- `page view` defaults to markdown.
- `--raw` returns storage/source content.
- `--content-only` omits metadata.
- `--no-truncate` disables the truncation guard.

## Integration Tests

Integration tests are off by default.

Suggested gates:

```bash
ATK_JIRA_INTEGRATION=1 go test ./tools/atk-jira/... -run Integration
ATK_CONFLUENCE_INTEGRATION=1 go test ./tools/atk-cfl/... -run Integration
```

Required environment:

Jira:

- `JIRA_URL` or `JIRA_BASE_URL`
- `JIRA_EMAIL`
- `JIRA_API_TOKEN`
- optional `JIRA_TEST_PROJECT`
- optional `JIRA_TEST_ISSUE`

Confluence:

- `CFL_URL` or `CONFLUENCE_URL`
- `CFL_EMAIL` or `CONFLUENCE_EMAIL`
- `CFL_API_TOKEN` or `CONFLUENCE_API_TOKEN`
- optional `CFL_TEST_SPACE`
- optional `CFL_TEST_PAGE`

Integration tests must:

- avoid destructive defaults
- create test resources with unique prefixes
- clean up when possible
- skip rather than fail when required environment is absent
- never print tokens

## Parity Coverage Tracking

Use a coverage table per command family:

| Area | Parser tests | API tests | Output tests | Integration smoke |
|---|---:|---:|---:|---:|
| Jira me | yes | yes | partial | manual |
| Jira issues | partial | partial | partial | manual |
| Jira comments | partial | partial | partial | manual |
| Jira transitions | no | no | no | no |
| Confluence page | partial | partial | partial | no |
| Confluence search | partial | partial | partial | no |

Update this table as commands are implemented.
