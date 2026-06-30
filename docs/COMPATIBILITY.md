# Command Surface Compatibility

[简体中文](COMPATIBILITY.zh-CN.md)

This project targets the command surface documented here while keeping our public binary names:

| Product area | This project |
|---|---|
| `atk-jira` | `atk-jira` |
| `atk-cfl` | `atk-cfl` |

The goal is to provide a stable operator surface, safety model, output discipline, and testing strategy under the `atk-*` binaries.

## Compatibility Rules

1. Resource names, action names, positional argument roles, and flag names should stay stable unless there is a documented reason to change them.
2. Jira command shape follows `atk-jira [resource] [action] [KEY/ID] [flags]`.
3. Confluence command shape follows `atk-cfl [resource] [action] [ID] [flags]`.
4. Jira text output should converge on a text-first, pipe-delimited, agent-friendly contract.
5. Confluence output should converge on the documented `-o table|plain`, `--full`, `--raw`, and page-content contracts.
6. Destructive operations require confirmation by default and `--force` to skip.
7. Non-destructive writes can run without prompting when they are reversible or naturally additive, but this project may keep `--dry-run` as an additional safety feature.
8. API behavior should follow current Atlassian Cloud APIs, not obsolete endpoints.

## Jira Surface Target

### Current User

- `atk-jira me`
- `atk-jira me --id`
- `atk-jira me --extended`

### Issues

- `atk-jira issues list --project KEY`
- `atk-jira issues get PROJ-123`
- `atk-jira issues create --project KEY --type TYPE --summary "TEXT"`
- `atk-jira issues update PROJ-123 --field FIELD=VALUE`
- `atk-jira issues search --jql "JQL"`
- `atk-jira issues assign PROJ-123 VALUE`
- `atk-jira issues assign PROJ-123 --unassign`
- `atk-jira issues move PROJ-1 [PROJ-2 ...] --to-project NEWPROJ`
- `atk-jira issues move-status TASK_ID`
- `atk-jira issues delete PROJ-123`
- `atk-jira issues types --project KEY`
- `atk-jira issues fields [PROJ-123]`
- `atk-jira issues field-options FIELD_NAME_OR_ID [--issue PROJ-123]`

Compatibility aliases during migration:

- `atk-jira issue get KEY` maps to `atk-jira issues get KEY`.
- `atk-jira issue search QUERY` maps to `atk-jira issues search --jql QUERY`.
- `atk-jira issue comment KEY --body TEXT` maps to `atk-jira comments add KEY --body TEXT`.

### Transitions

- `atk-jira transitions list PROJ-123`
- `atk-jira transitions list PROJ-123 --extended`
- `atk-jira transitions do PROJ-123 "Transition Name"`
- `atk-jira transitions do PROJ-123 21`
- `atk-jira transitions do PROJ-123 "Done" --field NAME=VALUE`

### Projects

- `atk-jira projects list`
- `atk-jira projects get KEY`

### Sprints

- `atk-jira sprints list --board ID_OR_NAME`
- `atk-jira sprints current --board ID_OR_NAME`
- `atk-jira sprints issues SPRINT_ID_OR_NAME`
- `atk-jira sprints add SPRINT_ID_OR_NAME PROJ-1 PROJ-2 ...`

### Boards

- `atk-jira boards list`
- `atk-jira boards get ID_OR_NAME`

### Links

- `atk-jira links list PROJ-123`
- `atk-jira links create PROJ-123 PROJ-456 --type NAME`
- `atk-jira links delete LINK_ID`
- `atk-jira links types`

### Comments

- `atk-jira comments list PROJ-123`
- `atk-jira comments add PROJ-123 --body "TEXT"`
- `atk-jira comments delete PROJ-123 COMMENT_ID`

### Attachments

- `atk-jira attachments list PROJ-123`
- `atk-jira attachments add PROJ-123 --file PATH`
- `atk-jira attachments get ATTACHMENT_ID`
- `atk-jira attachments delete ATTACHMENT_ID`

### Users

- `atk-jira users search "QUERY"`
- `atk-jira users get ACCOUNT_ID`

### Administrative Surfaces

These can follow after daily-use parity:

- `atk-jira fields ...`
- `atk-jira dashboards ...`
- `atk-jira automation ...`
- `atk-jira refresh ...`
- `atk-jira config ...`
- `atk-jira set-credential ...`
- `atk-jira completion ...`

## Confluence Surface Target

### Current User

- `atk-cfl me`
- `atk-cfl me --id`

### Pages

- `atk-cfl page list --space KEY`
- `atk-cfl page view PAGE_ID`
- `atk-cfl page view PAGE_ID --no-truncate`
- `atk-cfl page view PAGE_ID --content-only`
- `atk-cfl page view PAGE_ID --raw`
- `atk-cfl page view PAGE_ID --show-macros`
- `atk-cfl page view PAGE_ID --web`
- `atk-cfl page create --space KEY --title "TEXT"`
- `atk-cfl page create --space KEY --title "TEXT" --file content.md`
- `atk-cfl page edit PAGE_ID`
- `atk-cfl page edit PAGE_ID --file content.md`
- `atk-cfl page copy PAGE_ID --title "Copy Title"`
- `atk-cfl page delete PAGE_ID`
- `atk-cfl page delete PAGE_ID --force`

Compatibility aliases during migration:

- `atk-cfl page get ID` maps to `atk-cfl page view ID`.
- `atk-cfl page search QUERY` maps to `atk-cfl search QUERY`.
- `atk-cfl page update ID` maps to `atk-cfl page edit ID`.

### Search

- `atk-cfl search "query"`
- `atk-cfl search "query" --space KEY`
- `atk-cfl search "query" --type page`
- `atk-cfl search --label TAG`
- `atk-cfl search --title "TEXT"`
- `atk-cfl search --cql "CQL_QUERY"`

### Spaces

- `atk-cfl space list`
- `atk-cfl space view KEY`
- `atk-cfl space create --key KEY --name "NAME"`
- `atk-cfl space update KEY --name "NAME"`
- `atk-cfl space delete KEY`

### Attachments

- `atk-cfl attachment list --page PAGE_ID`
- `atk-cfl attachment upload --page PAGE_ID --file PATH`
- `atk-cfl attachment download ATT_ID`
- `atk-cfl attachment delete ATT_ID`

### Control Plane

- `atk-cfl init`
- `atk-cfl config show`
- `atk-cfl config test`
- `atk-cfl config clear`
- `atk-cfl set-credential`
- `atk-cfl completion`

## Environment Compatibility

Precedence:

```text
tool-specific env -> ATLASSIAN_* env -> config file -> defaults
```

Jira:

- `JIRA_URL` / `JIRA_BASE_URL`
- `JIRA_EMAIL`
- `JIRA_API_TOKEN` / `JIRA_TOKEN`
- `JIRA_AUTH_METHOD` / `JIRA_AUTH_TYPE`
- `JIRA_CLOUD_ID`
- `JIRA_DEFAULT_PROJECT`

Confluence:

- `CFL_URL` / `CONFLUENCE_URL`
- `CFL_EMAIL` / `CONFLUENCE_EMAIL`
- `CFL_API_TOKEN` / `CONFLUENCE_API_TOKEN`
- `CFL_AUTH_METHOD` / `CONFLUENCE_AUTH_METHOD`
- `CFL_CLOUD_ID` / `CONFLUENCE_CLOUD_ID`

Shared:

- `ATLASSIAN_URL` / `ATLASSIAN_SITE`
- `ATLASSIAN_EMAIL`
- `ATLASSIAN_API_TOKEN` / `ATLASSIAN_TOKEN`
- `ATLASSIAN_AUTH_METHOD` / `ATLASSIAN_AUTH_TYPE`
- `ATLASSIAN_CLOUD_ID`

## Output Compatibility

Jira should converge on this text-first model:

- Default: stable text.
- Lists: pipe-delimited rows with uppercase headers.
- Details: key-value blocks.
- `--id`: primary identifier only.
- `--extended`: adds admin/schema/audit fields.
- `--fulltext`: disables truncation.
- JSON is not a general-purpose Jira output mode, except where a round-trip payload demands it.

Confluence should converge on these output flags:

- `-o table` default.
- `-o plain` for script-oriented text.
- `-o json` is not a global resource output mode.
- `--full` adds inspection fields.
- `--raw` is command-specific, mainly for source-faithful page content.

This is a target. Current implementation still keeps `--json` on early commands until the output layer is migrated.
