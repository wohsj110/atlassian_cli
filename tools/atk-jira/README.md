# atk-jira - Jira CLI

A command-line interface for managing Jira Cloud tickets.

## Features

- Manage Jira issues from the command line
- List, create, update, search, delete, and inspect history for issues
- Manage projects (create, update, delete, restore)
- Manage sprints and boards
- Add comments and perform transitions
- Manage attachments
- Manage custom fields (create, delete, restore, contexts, options)
- Manage automation rules (list, export, create, update, delete, enable/disable)
- Manage dashboards and gadgets
- Create and manage issue links
- Search and look up users
- Text-first output with `--id`, `--extended`, and `--fulltext` modifiers
- Shell completion for bash, zsh, fish, and PowerShell

## Installation

### macOS

**Homebrew (recommended)**

```bash
brew install open-cli-collective/tap/atk-jira
```

> Note: This installs from our third-party tap.

---

### Windows

**Chocolatey**

```powershell
choco install atk-jira
```

**Winget**

```powershell
winget install wohsj110.atk-jira
```

---

### Linux

**Snap**

```bash
sudo snap install ocli-jira
```

> Note: After installation, the command is available as `atk-jira`.

**APT (Debian/Ubuntu)**

```bash
# Add the GPG key
curl -fsSL https://open-cli-collective.github.io/linux-packages/keys/gpg.asc | sudo gpg --dearmor -o /usr/share/keyrings/open-cli-collective.gpg

# Add the repository
echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/open-cli-collective.gpg] https://open-cli-collective.github.io/linux-packages/apt stable main" | sudo tee /etc/apt/sources.list.d/open-cli-collective.list

# Install
sudo apt update
sudo apt install atk-jira
```

> Note: This is our third-party APT repository, not official Debian/Ubuntu repos.

**DNF/YUM (Fedora/RHEL/CentOS)**

```bash
# Add the repository
sudo tee /etc/yum.repos.d/open-cli-collective.repo << 'EOF'
[open-cli-collective]
name=wohsj110
baseurl=https://open-cli-collective.github.io/linux-packages/rpm
enabled=1
gpgcheck=1
gpgkey=https://open-cli-collective.github.io/linux-packages/keys/gpg.asc
EOF

# Install
sudo dnf install atk-jira
```

> Note: This is our third-party RPM repository, not official Fedora/RHEL repos.

**Binary download**

Download `.deb`, `.rpm`, or `.tar.gz` from the [Releases page](https://github.com/wohsj110/atlassian_cli/releases) - available for x64 and ARM64.

```bash
# Direct .deb install
curl -LO https://github.com/wohsj110/atlassian_cli/releases/latest/download/atk-jira_VERSION_linux_amd64.deb
sudo dpkg -i atk-jira_VERSION_linux_amd64.deb

# Direct .rpm install
curl -LO https://github.com/wohsj110/atlassian_cli/releases/latest/download/atk-jira-VERSION.x86_64.rpm
sudo rpm -i atk-jira-VERSION.x86_64.rpm
```

---

### From Source

```bash
go install github.com/wohsj110/atlassian_cli/tools/atk-jira/cmd/atk-jira@latest
```

## Quick Start

### 1. Configure atk-jira

```bash
atk-jira init
```

This will prompt you for:
- Your Jira URL (e.g., `https://mycompany.atlassian.net`)
- Your email address
- An API token

Get your API token from: https://id.atlassian.com/manage-profile/security/api-tokens

### 2. List Issues

```bash
atk-jira issues list --project MYPROJECT
```

### 3. Get Issue Details

```bash
atk-jira issues get PROJ-123
```

---

## Command Reference

### Global Flags

These flags are available on all commands:

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--extended` | | `false` | Include admin/schema/audit fields in output |
| `--fulltext` | | `false` | Disable truncation of descriptions, comments, and history values |
| `--id` | | `false` | Emit only the primary identifier (takes precedence over `--extended` and `--fulltext`) |
| `--no-color` | | `false` | Disable colored output |
| `--verbose` | `-v` | `false` | Log each request's method/URL, JSON body, status, and any 4xx/5xx response body (each capped at 4 KB). Useful for diagnosing opaque Jira errors like `INVALID_INPUT`. |
| `--help` | `-h` | | Show help for command |
| `--version` | | | Show version (root command only) |

> `automation export` is the only command that emits JSON — it writes directly to stdout.

---

### `atk-jira init`

Initialize atk-jira with guided setup.

```bash
# Classic API token (Basic Auth — default)
atk-jira init
atk-jira init --url https://mycompany.atlassian.net --email user@example.com

# Service account with scoped token (Bearer Auth)
atk-jira init --auth-method bearer
atk-jira init --auth-method bearer --url https://mycompany.atlassian.net \
  --token YOUR_SCOPED_TOKEN --cloud-id YOUR_CLOUD_ID --no-verify
```

| Flag | Default | Description |
|------|---------|-------------|
| `--url` | | Jira URL (e.g., `https://mycompany.atlassian.net`) |
| `--email` | | Email address for authentication |
| `--token` | | API token |
| `--auth-method` | | Auth method: `basic` (default) or `bearer` |
| `--cloud-id` | | Cloud ID for bearer auth (find at `https://your-site.atlassian.net/_edge/tenant_info`) |
| `--no-verify` | `false` | Skip connection verification |

> **Bearer Auth:** For [Atlassian service accounts](https://support.atlassian.com/user-management/docs/manage-api-tokens-for-service-accounts/) with scoped API tokens. Email is not required. Requests route through the `api.atlassian.com` gateway.
>
> **Scope limitations:** Scoped tokens don't have scopes for Agile (boards/sprints), Automation, or Dashboards. These commands are unavailable with bearer auth — this is an Atlassian platform limitation.

---

### `atk-jira me`

Show information about the currently authenticated user.

```bash
atk-jira me
atk-jira me --id        # print just the account ID (for scripting)
atk-jira me --extended  # include timezone, locale, and group/application-role counts
```

Uses global flags `--id` and `--extended` — no command-specific flags.

---

### `atk-jira config`

Manage CLI configuration.

#### `atk-jira config show`

Display current configuration with masked credentials and source info.

```bash
atk-jira config show
```

#### `atk-jira config test`

Verify connection to Jira and test authentication.

```bash
atk-jira config test
```

#### `atk-jira config clear`

Remove stored configuration file.

```bash
atk-jira config clear
atk-jira config clear --force
```

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--force` | `-f` | `false` | Skip confirmation prompt |

---

### `atk-jira refresh`

Refresh the local instance cache (fields, projects, users, issue types, statuses, priorities, boards, link types, etc.).

With no arguments refreshes everything. With resource names, refreshes only those plus their dependencies. Use `--status` to check freshness without fetching.

Valid resource names: `fields`, `projects`, `users`, `issuetypes`, `statuses`, `priorities`, `resolutions`, `boards`, `sprints`, `linktypes`

```bash
# Refresh everything
atk-jira refresh

# Refresh specific resources (auto-expands dependencies)
atk-jira refresh statuses
atk-jira refresh users issuetypes

# Show cache freshness without fetching
atk-jira refresh --status
```

| Flag | Default | Description |
|------|---------|-------------|
| `--status` | `false` | Print cache freshness; no network calls |

---

### `atk-jira issues list`

List issues in a project.

**Aliases:** `atk-jira issue list`, `atk-jira i list`

```bash
atk-jira issues list --project MYPROJECT
atk-jira issues list --project MYPROJECT --sprint current
atk-jira issues list --project MYPROJECT --id

# Auto-pagination: fetch up to 200 results across multiple pages
atk-jira issues list --project MYPROJECT --max 200

# Explicit column projection
atk-jira issues list --project MYPROJECT --fields summary,status,customfield_10005
```

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--project` | `-p` | | Project key or name |
| `--sprint` | `-s` | | Filter by sprint: sprint name, numeric ID, or `current` |
| `--max` | `-m` | `50` | Maximum number of results to return |
| `--fields` | | | Comma-separated display columns (headers, Jira field IDs, or human names) |
| `--next-page-token` | | | Token for next page of results |

---

### `atk-jira issues get <issue-key> [issue-key...]`

Get details of one or more issues. A single key shows full detail; multiple keys show a summary table.

```bash
atk-jira issues get PROJ-123
atk-jira issues get PROJ-123 PROJ-456 PROJ-789
atk-jira issues get PROJ-123 --fulltext
atk-jira issues get PROJ-123 --id
atk-jira issues get PROJ-123 --fields Status,Assignee
atk-jira issues get PROJ-123 --custom-fields
```

| Flag | Default | Description |
|------|---------|-------------|
| `--fields` | | Comma-separated display fields (labels, Jira field IDs, or human names) |
| `--custom-fields` | `false` | Append custom fields section to output |
| `--fulltext` | `false` | Show full description without truncation (global) |
| `--id` | `false` | Emit only the issue key (global) |

**Arguments:**
- `<issue-key> [issue-key...]` - One or more issue keys (**required**)

---

### `atk-jira issues history <issue-key>`

List Jira changelog history for an issue as compact changed-field rows. Rows are chronological in Jira's changelog order.

```bash
atk-jira issues history PROJ-123
atk-jira issues history PROJ-123 --id
atk-jira issues history PROJ-123 --extended
atk-jira issues history PROJ-123 --fields CREATED,FIELD,TO
atk-jira issues history PROJ-123 --max 1
atk-jira issues history PROJ-123 --next-page-token 50
```

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--max` | `-m` | `50` | Maximum number of changelog groups to return |
| `--next-page-token` | | | Token for next page of results |
| `--fields` | | | Comma-separated display columns |
| `--extended` | | `false` | Include raw/audit history fields (global) |
| `--fulltext` | | `false` | Show full history values without truncation (global) |
| `--id` | | `false` | Emit changelog group IDs only (global) |

**Arguments:**
- `<issue-key>` - The issue key (**required**)

---

### `atk-jira issues create`

Create a new issue.

```bash
atk-jira issues create --project MYPROJECT --type Task --summary "Fix login bug"
atk-jira issues create -p MYPROJECT -t Story -s "Add new feature" --description "Details here"
atk-jira issues create -p MYPROJECT -s "Custom field issue" --field priority=High --field labels=backend

# Assign to yourself, by email, or by display name
atk-jira issues create -p MYPROJECT -t Task -s "My task" --assignee me
atk-jira issues create -p MYPROJECT -t Task -s "Their task" --assignee user@example.com
atk-jira issues create -p MYPROJECT -t Task -s "Their task" --assignee "Aaron Wong"
```

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--project` | `-p` | | Project key or name (**required**) |
| `--type` | `-t` | `Task` | Issue type: `Task`, `Bug`, `Story`, etc. |
| `--summary` | `-s` | | Issue summary (**required**) |
| `--description` | `-d` | | Issue description (supports `\n`, `\t`, `\\` escape sequences) |
| `--parent` | | | Parent issue key (epic or parent issue) |
| `--assignee` | `-a` | | Assignee (account ID, email, display name, or `"me"`) |
| `--field` | `-f` | | Additional field in `key=value` format (can be repeated) |

---

### `atk-jira issues update <issue-key>`

Update an existing issue.

```bash
atk-jira issues update PROJ-123 --summary "New summary"
atk-jira issues update PROJ-123 --field priority=High
atk-jira issues update PROJ-123 --description "Updated description" --field labels=urgent

# Unassign an issue
atk-jira issues update PROJ-123 --assignee none

# Change workflow status (routes to the transitions API under the hood).
# Quote multi-word status names: --status "In Progress"
atk-jira issues update PROJ-123 --status "Done"

# Multi-value fields: repeat --field with the same key to accumulate values
atk-jira issues update PROJ-123 --field customfield_10050=Option1 --field customfield_10050=Option2
```

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--summary` | `-s` | | New summary |
| `--description` | `-d` | | New description (supports `\n`, `\t`, `\\` escape sequences) |
| `--parent` | | | Parent issue key (epic or parent issue) |
| `--assignee` | `-a` | | Assignee (account ID, email, display name, `"me"`, or `"none"` to unassign) |
| `--type` | `-t` | | New issue type (uses Jira Cloud bulk move API) |
| `--status` | | | New workflow status (uses Jira transitions API; resolved before any writes) |
| `--field` | `-f` | | Field to update in `key=value` format (can be repeated; repeating the same key accumulates values for multi-select fields) |

**Arguments:**
- `<issue-key>` - The issue key (**required**)

---

### `atk-jira issues search`

Search issues using JQL.

```bash
atk-jira issues search --jql "project = MYPROJECT AND status = 'In Progress'"
atk-jira issues search --jql "assignee = currentUser()" --id

# Auto-pagination: fetch up to 200 results across multiple pages
atk-jira issues search --jql "project = MYPROJECT" --max 200

# Explicit column projection
atk-jira issues search --jql "project = MYPROJECT" --fields summary,status
```

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--jql` | | | JQL query string (**required**) |
| `--max` | `-m` | `50` | Maximum number of results to return |
| `--fields` | | | Comma-separated display columns (headers, Jira field IDs, or human names) |
| `--next-page-token` | | | Token for next page of results |

---

### `atk-jira issues assign <issue-key> [user]`

Assign an issue to a user, or unassign it. The `[user]` argument accepts an account ID, email, display name, or `"me"`.

```bash
atk-jira issues assign PROJ-123 5b10ac8d82e05b22cc7d4ef5
atk-jira issues assign PROJ-123 "Aaron Wong"
atk-jira issues assign PROJ-123 aaron@example.com
atk-jira issues assign PROJ-123 me
atk-jira issues assign PROJ-123 --unassign
```

| Flag | Default | Description |
|------|---------|-------------|
| `--unassign` | `false` | Remove current assignee |

**Arguments:**
- `<issue-key>` - The issue key (**required**)
- `[user]` - Account ID, email, display name, or `"me"` (required unless `--unassign`)

---

### `atk-jira issues delete <issue-key> [issue-key...]`

Delete one or more issues.

```bash
atk-jira issues delete PROJ-123
atk-jira issues delete PROJ-123 PROJ-124 PROJ-125
atk-jira issues delete PROJ-123 --force
```

| Flag | Default | Description |
|------|---------|-------------|
| `--force` | `false` | Skip confirmation prompt |

**Arguments:**
- `<issue-key> [issue-key...]` - One or more issue keys (**required**)

---

### `atk-jira issues archive <issue-key> [issue-key...]`

Archive one or more issues. Archived issues are hidden from boards and search by default but remain in Jira. There is no `issues restore` command — use the Jira UI to unarchive.

```bash
atk-jira issues archive PROJ-123
atk-jira issues archive PROJ-123 PROJ-124 PROJ-125
atk-jira issues archive PROJ-123 --id
```

**Arguments:**
- `<issue-key> [issue-key...]` - One or more issue keys (**required**)

---

### `atk-jira issues check <issue-key>`

Check whether an issue has values for expected fields. Useful as a guardrail
before transitions or as a CI step. Each field can be named by its display name
(e.g. `Story Points`), Jira field ID (e.g. `customfield_10035`), or property
key (e.g. `assignee`).

```bash
# Default warn list (Summary, Description, Assignee, Priority, Labels,
# Story Points, Sprint, Components, Fix Version/s) — fields not on the
# project's schema are silently skipped.
atk-jira issues check PROJ-123

# Hard-fail (non-zero exit) if Story Points or Sprint are missing.
atk-jira issues check PROJ-123 --require "Story Points" --require Sprint

# Mix required and warning fields, comma-separated.
atk-jira issues check PROJ-123 --require "Story Points,Sprint" --warn "Description,Assignee"

# Emit only the IDs of MISSING fields.
atk-jira issues check PROJ-123 --require Sprint --id
```

| Flag | Default | Description |
|------|---------|-------------|
| `--require` | (none) | Field must be populated; missing → non-zero exit (repeatable) |
| `--warn` | (curated list, only when neither flag is provided) | Field flagged if missing; never fails the check (repeatable) |

When `--require` is provided alone, the curated default warn-list is **not** applied — only the explicitly-named fields are checked.

Use `--id` to emit only the IDs of fields whose status is `MISSING`.

**Exit codes:** `0` if all `--require` fields populated; `1` if any are missing.

**Arguments:**
- `<issue-key>` - The issue key (**required**)

---

### `atk-jira issues fields [issue-key]`

List available fields, or show all fields with their current values for a specific issue.

```bash
atk-jira issues fields                    # All fields
atk-jira issues fields PROJ-123           # Field values for a specific issue
atk-jira issues fields --custom-fields    # Custom fields only
```

| Flag | Default | Description |
|------|---------|-------------|
| `--custom-fields` | `false` | Show only custom fields |

**Arguments:**
- `[issue-key]` - Optional issue key to show field values

---

### `atk-jira issues field-options [issue-key] <field-name-or-id>`

List allowed values for a field. Providing an issue key uses that issue's project context (recommended); omitting it uses the global field context.

The first positional argument is treated as an issue key if it matches the `PROJ-123` pattern (uppercase letters/digits, hyphen, digits); otherwise it is treated as the field name.

```bash
atk-jira issues field-options PROJ-123 priority
atk-jira issues field-options PROJ-123 customfield_10001
atk-jira issues field-options priority   # without issue context (single arg = field name)
```

**Arguments:**
- `[issue-key]` - Optional issue key for context-specific options (must match `KEY-NNN` pattern)
- `<field-name-or-id>` - Field name or ID (**required**)

---

### `atk-jira issues types`

List available issue types for a project.

```bash
atk-jira issues types --project MYPROJECT
```

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--project` | `-p` | | Project key (**required**) |

---

### `atk-jira issues move <issue-key>...`

Move one or more issues to a different project (Cloud only, max 1000 issues).

```bash
atk-jira issues move PROJ-123 --to-project OTHERPROJ
atk-jira issues move PROJ-123 PROJ-124 PROJ-125 --to-project OTHERPROJ --to-type Bug

# Move without waiting for completion
atk-jira issues move PROJ-123 --to-project OTHERPROJ --no-wait

# Move without notifications
atk-jira issues move PROJ-123 --to-project OTHERPROJ --no-notify
```

| Flag | Default | Description |
|------|---------|-------------|
| `--to-project` | | Target project key or name (**required**) |
| `--to-type` | (same as source) | Target issue type name |
| `--notify` | `true` | Send notifications; use `--no-notify` to disable |
| `--wait` | `true` | Wait for move to complete; use `--no-wait` to return immediately with the task ID |

**Arguments:**
- `<issue-key>...` - One or more issue keys (**required**)

---

### `atk-jira issues move-status <task-id>`

Check the status of an asynchronous move operation.

```bash
atk-jira issues move-status 12345
```

**Arguments:**
- `<task-id>` - The task ID returned by `issues move` (**required**)

---

### `atk-jira links list <issue-key>`

List all links on an issue.

**Aliases:** `atk-jira link list`, `atk-jira l list`

```bash
atk-jira links list PROJ-123
atk-jira links list PROJ-123 --id
atk-jira links list PROJ-123 --fields TYPE,ISSUE
```

| Flag | Default | Description |
|------|---------|-------------|
| `--fields` | | Comma-separated display columns |

**Arguments:**
- `<issue-key>` - The issue key (**required**)

---

### `atk-jira links create <issue-key> <target-issue-key>`

Create a link between two issues. The first issue is the outward issue and the second is the inward issue. `--type` accepts the canonical name, the outward verb, or the inward verb.

```bash
# A blocks B
atk-jira links create PROJ-123 PROJ-456 --type Blocks

# A is blocked by B (inward verb — issues are interpreted from user's perspective)
atk-jira links create PROJ-123 PROJ-456 --type "is blocked by"

# A relates to B
atk-jira links create PROJ-123 PROJ-456 --type Relates
```

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--type` | `-t` | | Link type: canonical name, outward verb, or inward verb (**required**) |

**Arguments:**
- `<issue-key>` - The outward issue key (**required**)
- `<target-issue-key>` - The inward issue key (**required**)

> Tip: Use `atk-jira links types` to see available link types.

---

### `atk-jira links delete <link-id>`

Delete an issue link.

```bash
atk-jira links delete 10001
```

**Arguments:**
- `<link-id>` - The link ID (**required**)

> Tip: Use `atk-jira links list PROJ-123` to find link IDs.

---

### `atk-jira links types`

List available issue link types.

```bash
atk-jira links types
atk-jira links types --id
```

| Flag | Default | Description |
|------|---------|-------------|
| `--fields` | | Comma-separated display columns |

---

### `atk-jira transitions list <issue-key>`

List available transitions for an issue.

**Aliases:** `atk-jira transition list`, `atk-jira tr list`

```bash
atk-jira transitions list PROJ-123
atk-jira transitions list PROJ-123 --extended
atk-jira transitions list PROJ-123 --id
```

**Arguments:**
- `<issue-key>` - The issue key (**required**)

---

### `atk-jira transitions do <issue-key> <transition>`

Perform a transition on an issue. For ordinary status changes, prefer
`atk-jira issues update <key> --status <name>` — it hides the Jira API split.
Reach for `transitions do` when you need to disambiguate multiple
transitions to the same target status, set fields-on-transition, or pick
a transition by ID.

**Aliases:** `atk-jira transition do`, `atk-jira tr do`

```bash
atk-jira transitions do PROJ-123 "In Progress"
atk-jira transitions do PROJ-123 "Done"
atk-jira transitions do PROJ-123 "Done" --field resolution=Fixed
```

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--field` | `-f` | | Field to set during transition in `key=value` format (can be repeated) |

**Arguments:**
- `<issue-key>` - The issue key (**required**)
- `<transition>` - Transition name or ID (**required**)

---

### `atk-jira comments list <issue-key>`

List comments on an issue.

**Aliases:** `atk-jira comment list`, `atk-jira c list`

```bash
atk-jira comments list PROJ-123
atk-jira comments list PROJ-123 --fulltext
atk-jira comments list PROJ-123 --fields ID,AUTHOR
```

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--max` | `-m` | `50` | Maximum number of comments |
| `--fulltext` | | `false` | Show full comment bodies without truncation (global) |
| `--fields` | | | Comma-separated display fields |

**Arguments:**
- `<issue-key>` - The issue key (**required**)

---

### `atk-jira comments add <issue-key>`

Add a comment to an issue.

**Aliases:** `atk-jira comment add`, `atk-jira c add`

```bash
atk-jira comments add PROJ-123 --body "This is my comment"
atk-jira comments add PROJ-123 --body "Line one\nLine two\n\tIndented line"
```

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--body` | `-b` | | Comment text (supports `\n`, `\t`, `\\` escape sequences) (**required**) |

**Arguments:**
- `<issue-key>` - The issue key (**required**)

---

### `atk-jira comments delete <issue-key> <comment-id>`

Delete a comment from an issue.

**Aliases:** `atk-jira comment delete`, `atk-jira c delete`

```bash
atk-jira comments delete PROJ-123 10042
```

**Arguments:**
- `<issue-key>` - The issue key (**required**)
- `<comment-id>` - The comment ID (**required**)

---

### `atk-jira attachments list <issue-key>`

List attachments on an issue.

**Aliases:** `atk-jira attachments ls`, `atk-jira attachment list`, `atk-jira att list`

```bash
atk-jira attachments list PROJ-123
atk-jira attachments list PROJ-123 --id
```

| Flag | Default | Description |
|------|---------|-------------|
| `--fields` | | Comma-separated display columns |

**Arguments:**
- `<issue-key>` - The issue key (**required**)

---

### `atk-jira attachments add <issue-key>`

Upload file(s) to an issue.

**Aliases:** `atk-jira attachment add`, `atk-jira att add`

```bash
atk-jira attachments add PROJ-123 --file screenshot.png
atk-jira attachments add PROJ-123 --file doc.pdf --file image.png
```

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--file` | `-F` | | File to attach (**required**, can be repeated) |

**Arguments:**
- `<issue-key>` - The issue key (**required**)

---

### `atk-jira attachments get <attachment-id>`

Download an attachment.

**Aliases:** `atk-jira attachments download`, `atk-jira attachment get`, `atk-jira att get`

```bash
atk-jira attachments get 12345
atk-jira attachments get 12345 --output ./downloads/
```

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--output` | `-o` | `.` | Output path (directory or filename) |

**Arguments:**
- `<attachment-id>` - The attachment ID (**required**)

---

### `atk-jira attachments delete <attachment-id>`

Delete an attachment.

**Aliases:** `atk-jira attachments rm`, `atk-jira attachment delete`, `atk-jira att delete`

```bash
atk-jira attachments delete 12345
```

**Arguments:**
- `<attachment-id>` - The attachment ID (**required**)

---

### `atk-jira sprints list`

List sprints for a board. `--board` accepts a board ID or name.

```bash
atk-jira sprints list --board 123
atk-jira sprints list --board "MON board" --state active
atk-jira sprints list --board 123 --id
```

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--board` | `-b` | | Board ID or name (**required**) |
| `--state` | `-s` | | Filter by state: `active`, `closed`, `future` |
| `--max` | `-m` | `50` | Maximum number of results |
| `--fields` | | | Comma-separated display columns |
| `--next-page-token` | | | Token for next page of results |

---

### `atk-jira sprints current`

Show the current active sprint.

```bash
atk-jira sprints current --board 123
atk-jira sprints current --board "MON board"
```

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--board` | `-b` | | Board ID or name (**required**) |
| `--fields` | | | Comma-separated display fields |

---

### `atk-jira sprints issues <sprint>`

List issues in a sprint. Accepts a sprint ID or name (resolved via cache).

```bash
atk-jira sprints issues 456
atk-jira sprints issues "MON Sprint 70"
atk-jira sprints issues 456 --id
atk-jira sprints issues 456 --fields KEY,STATUS,customfield_10005
```

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--max` | `-m` | `50` | Maximum number of results |
| `--fields` | | | Comma-separated display columns |
| `--next-page-token` | | | Token for next page of results |

**Arguments:**
- `<sprint>` - Sprint ID or name (**required**)

---

### `atk-jira sprints add <sprint> <issue-key>...`

Move one or more issues to a sprint. Accepts a sprint ID or name.

```bash
atk-jira sprints add 456 PROJ-123
atk-jira sprints add "MON Sprint 70" PROJ-123
atk-jira sprints add 456 PROJ-123 PROJ-124 PROJ-125
```

**Arguments:**
- `<sprint>` - Sprint ID or name (**required**)
- `<issue-key>...` - One or more issue keys (**required**)

---

### `atk-jira sprints remove <issue-key>...`

Move one or more issues from their current sprint to the backlog.

```bash
atk-jira sprints remove PROJ-456
atk-jira sprints remove PROJ-456 PROJ-789 PROJ-101
```

**Arguments:**
- `<issue-key>...` - One or more issue keys (**required**)

---

### `atk-jira boards list`

List boards.

```bash
atk-jira boards list
atk-jira boards list --project MYPROJECT
```

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--project` | `-p` | | Filter by project key or name |
| `--max` | `-m` | `50` | Maximum number of results |
| `--fields` | | | Comma-separated display columns |
| `--next-page-token` | | | Token for next page of results |

---

### `atk-jira boards get <board>`

Get board details. Accepts a board ID or name (resolved via cache).

```bash
atk-jira boards get 123
atk-jira boards get "MON board"
```

| Flag | Default | Description |
|------|---------|-------------|
| `--fields` | | Comma-separated display fields |

**Arguments:**
- `<board>` - Board ID or name (**required**)

---

### `atk-jira projects list`

List Jira projects.

**Aliases:** `atk-jira project list`, `atk-jira proj list`, `atk-jira p list`

```bash
atk-jira projects list
atk-jira projects list --query "my project"
atk-jira projects list --max 10
```

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--query` | `-q` | | Filter projects by name |
| `--max` | `-m` | `50` | Maximum number of results |
| `--fields` | | | Comma-separated display columns |
| `--next-page-token` | | | Token for next page of results |

---

### `atk-jira projects get <project-key>`

Get details for a specific project.

```bash
atk-jira projects get MYPROJECT
atk-jira projects get 10001
```

| Flag | Default | Description |
|------|---------|-------------|
| `--fields` | | Comma-separated display fields |

**Arguments:**
- `<project-key>` - Project key or numeric ID (**required**)

---

### `atk-jira projects create`

Create a new Jira project.

```bash
atk-jira projects create --key MYPROJ --name "My Project" --lead me
atk-jira projects create --key BIZ --name "Business" --type business --lead "Aaron Wong" --description "Business project"
```

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--key` | `-k` | | Project key (**required**) |
| `--name` | `-n` | | Project name (**required**) |
| `--type` | `-t` | `software` | Project type: `software`, `service_desk`, `business` |
| `--lead` | `-l` | | Lead: account ID, email, display name, or `"me"` (**required**) |
| `--description` | `-d` | | Project description |

> Tip: Use `atk-jira users search` to find account IDs, or `atk-jira me` to get your own.

---

### `atk-jira projects update <project-key>`

Update a project's metadata. Only specified fields are changed.

```bash
atk-jira projects update MYPROJ --name "New Name"
atk-jira projects update MYPROJ --description "Updated description"
atk-jira projects update MYPROJ --lead "Aaron Wong"
atk-jira projects update MYPROJ --lead aaron@example.com
```

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--name` | `-n` | | New project name |
| `--description` | `-d` | | New project description |
| `--lead` | `-l` | | New lead: account ID, email, display name, or `"me"` |

**Arguments:**
- `<project-key>` - Project key (**required**)

---

### `atk-jira projects delete <project-key>`

Soft-delete a project (moves it to trash). Can be restored with `atk-jira projects restore`.

```bash
atk-jira projects delete MYPROJ
atk-jira projects delete MYPROJ --force
```

| Flag | Default | Description |
|------|---------|-------------|
| `--force` | `false` | Skip confirmation prompt |

**Arguments:**
- `<project-key>` - Project key (**required**)

---

### `atk-jira projects restore <project-key>`

Restore a project from the trash.

```bash
atk-jira projects restore MYPROJ
```

**Arguments:**
- `<project-key>` - Project key (**required**)

---

### `atk-jira projects types`

List available project types for creating new projects.

```bash
atk-jira projects types
```

---

### `atk-jira users search <query>`

Search for Jira users.

**Aliases:** `atk-jira user search`, `atk-jira u search`

```bash
atk-jira users search "john"
atk-jira users search "john" --max 20
```

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--max` | `-m` | `50` | Maximum number of results |
| `--fields` | | | Comma-separated display columns |
| `--next-page-token` | | | Token for next page of results |

**Arguments:**
- `<query>` - Search query (matches display name, email, etc.) (**required**)

---

### `atk-jira users get <account-id>`

Get details for a specific user by account ID.

**Aliases:** `atk-jira user get`, `atk-jira u get`

```bash
atk-jira users get 5b10ac8d82e05b22cc7d4ef5
atk-jira users get 5b10ac8d82e05b22cc7d4ef5 --id     # global flag: emit only account ID
atk-jira users get 5b10ac8d82e05b22cc7d4ef5 --extended
```

| Flag | Default | Description |
|------|---------|-------------|
| `--fields` | | Comma-separated display fields |
| `--id` | `false` | Emit only the account ID (global flag) |

**Arguments:**
- `<account-id>` - The Atlassian account ID (**required**)

---

### `atk-jira automation list`

List automation rules.

**Aliases:** `atk-jira auto list`

```bash
atk-jira automation list
atk-jira automation list --state ENABLED
```

| Flag | Default | Description |
|------|---------|-------------|
| `--state` | | Filter by state: `ENABLED` or `DISABLED` |

---

### `atk-jira automation get <rule-id>`

Get details of an automation rule.

**Aliases:** `atk-jira auto get`

```bash
atk-jira automation get 123
atk-jira automation get 123 --show-components
```

| Flag | Default | Description |
|------|---------|-------------|
| `--show-components` | `false` | Show component type details |

**Arguments:**
- `<rule-id>` - The rule ID (**required**)

---

### `atk-jira automation export <rule-id>`

Export a rule definition as JSON.

**Aliases:** `atk-jira auto export`

```bash
atk-jira automation export 123
atk-jira automation export 123 --compact
atk-jira automation export 123 > rule-backup.json
```

| Flag | Default | Description |
|------|---------|-------------|
| `--compact` | `false` | Output minified JSON |

**Arguments:**
- `<rule-id>` - The rule ID (**required**)

> Note: Output is always JSON — this is the only atk-jira command that emits JSON directly.

---

### `atk-jira automation create`

Create an automation rule from a JSON file.

**Aliases:** `atk-jira auto create`

```bash
atk-jira automation create --file rule-definition.json
```

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--file` | `-F` | | Path to JSON file containing the rule definition (**required**) |

> Note: New rules are created in DISABLED state by default.

---

### `atk-jira automation update <rule-id>`

Update an automation rule from a JSON file.

**Aliases:** `atk-jira auto update`

```bash
atk-jira automation update 123 --file updated-rule.json
```

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--file` | `-F` | | Path to JSON file containing the rule definition (**required**) |

**Arguments:**
- `<rule-id>` - The rule ID (**required**)

> Tip: Use `atk-jira automation export` to get the current definition before editing.

---

### `atk-jira automation delete <rule-id>`

Permanently delete an automation rule. If the rule is currently ENABLED, it will be automatically disabled before deletion. This action cannot be undone.

**Aliases:** `atk-jira auto delete`

```bash
atk-jira automation delete 123
atk-jira automation delete 123 --force
```

| Flag | Default | Description |
|------|---------|-------------|
| `--force` | `false` | Skip confirmation prompt |

**Arguments:**
- `<rule-id>` - The rule ID (**required**)

---

### `atk-jira automation enable <rule-id>`

Enable a disabled automation rule.

**Aliases:** `atk-jira auto enable`

```bash
atk-jira automation enable 123
```

**Arguments:**
- `<rule-id>` - The rule ID (**required**)

---

### `atk-jira automation disable <rule-id>`

Disable an enabled automation rule.

**Aliases:** `atk-jira auto disable`

```bash
atk-jira automation disable 123
```

**Arguments:**
- `<rule-id>` - The rule ID (**required**)

---

### `atk-jira dashboards list`

List accessible dashboards.

**Aliases:** `atk-jira dashboard list`, `atk-jira dash list`

```bash
atk-jira dashboards list
atk-jira dashboards list --search "Sprint"
atk-jira dashboards list --max 10
```

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--search` | | | Search dashboards by name |
| `--max` | `-m` | `50` | Maximum number of results |

> Note: Dashboard commands are not available with bearer auth (scoped tokens lack the Dashboard scope).

---

### `atk-jira dashboards get <dashboard-id>`

Get dashboard details including gadgets.

```bash
atk-jira dashboards get 10001
```

**Arguments:**
- `<dashboard-id>` - The dashboard ID (**required**)

---

### `atk-jira dashboards create`

Create a new dashboard.

```bash
atk-jira dashboards create --name "My Dashboard"
atk-jira dashboards create --name "Sprint Board" --description "Sprint tracking"
```

| Flag | Default | Description |
|------|---------|-------------|
| `--name` | | Dashboard name (**required**) |
| `--description` | | Dashboard description |

---

### `atk-jira dashboards delete <dashboard-id>`

Delete a dashboard.

```bash
atk-jira dashboards delete 10001
```

**Arguments:**
- `<dashboard-id>` - The dashboard ID (**required**)

---

### `atk-jira dashboards gadgets list <dashboard-id>`

List gadgets on a dashboard.

```bash
atk-jira dashboards gadgets list 10001
atk-jira dashboards gadgets list 10001 --id
```

**Arguments:**
- `<dashboard-id>` - The dashboard ID (**required**)

---

### `atk-jira dashboards gadgets add <dashboard-id>`

Add a gadget to a dashboard by its module key.

```bash
atk-jira dashboards gadgets add 10001 --type com.atlassian.jira.gadgets:sprint-burndown-gadget
atk-jira dashboards gadgets add 10001 --type com.atlassian.jira.gadgets:filter-results-gadget --position 1,0 --title "My Filter"
```

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--type` | `-t` | | Gadget module key (**required**) |
| `--position` | `-p` | | Position as `row,column` (e.g. `1,0`) |
| `--title` | | | Gadget title |
| `--color` | | | Gadget color |

**Arguments:**
- `<dashboard-id>` - The dashboard ID (**required**)

---

### `atk-jira dashboards gadgets remove <dashboard-id> <gadget-id>`

Remove a gadget from a dashboard.

```bash
atk-jira dashboards gadgets remove 10001 42
```

**Arguments:**
- `<dashboard-id>` - The dashboard ID (**required**)
- `<gadget-id>` - The gadget ID (**required**)

---

### `atk-jira fields list`

List all fields (system and custom). Supports filtering by name with case-insensitive substring matching.

**Aliases:** `atk-jira field list`, `atk-jira f list`

```bash
atk-jira fields list
atk-jira fields list --custom-fields
atk-jira fields list --name "story point"
atk-jira fields list --id
```

| Flag | Default | Description |
|------|---------|-------------|
| `--custom-fields` | `false` | Show only custom fields |
| `--name` | | Filter fields by name (case-insensitive substring match) |

#### `atk-jira fields show <field-id>`

Show a flat denormalized view of a field's contexts, project mappings, and options.

```bash
atk-jira fields show customfield_10001
atk-jira fields show customfield_10001 --id   # emit context IDs only
```

**Arguments:**
- `<field-id>` - The field ID (**required**)

#### `atk-jira fields create`

Create a new custom field.

```bash
atk-jira fields create --name "My Select Field" --type com.atlassian.jira.plugin.system.customfieldtypes:select
atk-jira fields create --name "My Text Field" --type com.atlassian.jira.plugin.system.customfieldtypes:textarea --description "A text area field"
```

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--name` | `-n` | | Field name (**required**) |
| `--type` | `-t` | | Field type (**required**) |
| `--description` | `-d` | | Field description |

#### `atk-jira fields delete <field-id>`

Trash a custom field (can be restored).

```bash
atk-jira fields delete customfield_10001
atk-jira fields delete customfield_10001 --force
```

| Flag | Default | Description |
|------|---------|-------------|
| `--force` | `false` | Skip confirmation prompt |

**Arguments:**
- `<field-id>` - The field ID (**required**)

#### `atk-jira fields restore <field-id>`

Restore a trashed custom field.

```bash
atk-jira fields restore customfield_10001
```

**Arguments:**
- `<field-id>` - The field ID (**required**)

#### `atk-jira fields contexts list <field-id>`

List contexts for a custom field.

**Aliases:** `atk-jira fields context list`, `atk-jira fields ctx list`

```bash
atk-jira fields contexts list customfield_10001
atk-jira fields contexts list customfield_10001 --id
```

**Arguments:**
- `<field-id>` - The field ID (**required**)

#### `atk-jira fields contexts create <field-id>`

Create a context for a custom field.

**Aliases:** `atk-jira fields context create`, `atk-jira fields ctx create`

```bash
atk-jira fields contexts create customfield_10001 --name "My Context"
atk-jira fields contexts create customfield_10001 --name "Project Context" --project 10001
```

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--name` | `-n` | | Context name (**required**) |
| `--project` | `-p` | | Project ID to scope the context to |

**Arguments:**
- `<field-id>` - The field ID (**required**)

#### `atk-jira fields contexts delete <field-id> <context-id>`

Delete a context from a custom field.

**Aliases:** `atk-jira fields context delete`, `atk-jira fields ctx delete`

```bash
atk-jira fields contexts delete customfield_10001 10100
atk-jira fields contexts delete customfield_10001 10100 --force
```

| Flag | Default | Description |
|------|---------|-------------|
| `--force` | `false` | Skip confirmation prompt |

**Arguments:**
- `<field-id>` - The field ID (**required**)
- `<context-id>` - The context ID (**required**)

#### `atk-jira fields options list <field-id>`

List options for a select/multi-select custom field. Auto-detects the default context if `--context` is not specified.

**Aliases:** `atk-jira fields option list`, `atk-jira fields opt list`

```bash
atk-jira fields options list customfield_10001
atk-jira fields options list customfield_10001 --context 10001
```

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--context` | `-c` | | Context ID (auto-detected if omitted) |

**Arguments:**
- `<field-id>` - The field ID (**required**)

#### `atk-jira fields options add <field-id>`

Add an option to a select/multi-select custom field.

**Aliases:** `atk-jira fields option add`, `atk-jira fields opt add`

```bash
atk-jira fields options add customfield_10001 --value "New Option"
atk-jira fields options add customfield_10001 --value "Staging" --context 10001
```

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--value` | `-V` | | Option value (**required**) |
| `--context` | `-c` | | Context ID (auto-detected if omitted) |

**Arguments:**
- `<field-id>` - The field ID (**required**)

#### `atk-jira fields options update <field-id>`

Update an existing option value.

**Aliases:** `atk-jira fields option update`, `atk-jira fields opt update`

```bash
atk-jira fields options update customfield_10001 --option 10200 --value "Updated Value"
```

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--option` | | | Option ID to update (**required**) |
| `--value` | `-V` | | New option value (**required**) |
| `--context` | `-c` | | Context ID (auto-detected if omitted) |

**Arguments:**
- `<field-id>` - The field ID (**required**)

#### `atk-jira fields options delete <field-id>`

Delete an option from a select/multi-select custom field.

**Aliases:** `atk-jira fields option delete`, `atk-jira fields opt delete`

```bash
atk-jira fields options delete customfield_10001 --option 10200
atk-jira fields options delete customfield_10001 --option 10200 --force
```

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--option` | | | Option ID to delete (**required**) |
| `--force` | | `false` | Skip confirmation prompt |
| `--context` | `-c` | | Context ID (auto-detected if omitted) |

**Arguments:**
- `<field-id>` - The field ID (**required**)

---

## Configuration

`atk-jira init` stores the **API token in your OS keyring** (macOS Keychain /
Linux Secret Service / Windows Credential Manager, or an opt-in
encrypted-file backend) and writes only **non-secret** config to the
shared store at `OS-native atlassian-agent-cli/config.yml`:

```yaml
default:
  url: https://mycompany.atlassian.net
  email: user@example.com
  auth_method: basic                # or "bearer"
  cloud_id: ""                      # required for bearer
atk_jira:
  default_project: MYPROJECT        # atk-jira-only defaults; internal section key
```

There is **no `api_token:` field** — the secret never touches a
plaintext file. The same config file and keyring bundle are shared with
`atk-cfl` — one Atlassian token, both tools. Run `atk-jira init` after `atk-cfl init`
(or vice versa) and you'll be offered to reuse the credentials.

**Non-interactive token ingress (§1.5.2):** use `atk-jira set-credential` for
installer scripts, CI, and credential-manager driven setup. Required
flags:

- `--ref atlassian-agent-cli/default` (required on fresh installs; defaults to
  the canonical ref when a shared config already exists)
- `--key api_token` (always required)
- exactly one of `--stdin` or `--from-env VAR` (mutually exclusive; no
  `--value` — flag-passed secrets leak into process listings)
- `--overwrite` to replace an existing entry (default: fail loud)
- `--json` to emit the §1.5.2 control-plane envelope
  `{"ref","key","backend","written","error?"}` on stdout (stderr stays
  empty under `--json`)

```bash
# From a secrets manager
op read 'op://Vault/Atlassian/token' | atk-jira set-credential \
  --ref atlassian-agent-cli/default --key api_token --stdin

# From an environment variable
atk-jira set-credential --ref atlassian-agent-cli/default --key api_token \
  --from-env JIRA_API_TOKEN

# Replace an existing entry
op read 'op://Vault/Atlassian/token' | atk-jira set-credential \
  --ref atlassian-agent-cli/default --key api_token --stdin --overwrite

# Installer-script control-plane envelope
atk-jira set-credential --ref atlassian-agent-cli/default --key api_token \
  --from-env JIRA_API_TOKEN --json
```

Legacy per-tool config keeps working indefinitely (Linux: `~/.config/atk-jira/config.json`; macOS: `~/Library/Application Support/atk-jira/config.json`). The first command auto-migrates any pre-existing plaintext token into the keyring and scrubs the plaintext in place.

Run `atk-jira config show` to inspect the resolved values, including the keyring ref, backend, and whether a token is configured (the token value itself is never displayed). Token/keyring reporting is authoritative; the non-secret rows reflect env + the legacy per-tool file only, so a value set solely in the shared store appears as "-" there even though atk-jira uses it at runtime. `atk-jira config clear` removes the single shared `api_token` (warning that atk-cfl loses access too, since both tools resolve the same key); `atk-jira config clear --all` removes the whole bundle (including any deprecated per-tool keys) plus the non-secret config file.

### Environment Variables

Environment variables override file-based config. Variables are checked in order of precedence (first match wins):

| Setting | Precedence (highest to lowest) |
|---------|-------------------------------|
| URL | `JIRA_URL` → `ATLASSIAN_URL` → shared `default` → legacy → `JIRA_DOMAIN` |
| Email | `JIRA_EMAIL` → `ATLASSIAN_EMAIL` → shared `default` → legacy |
| API Token | `JIRA_API_TOKEN` → `ATLASSIAN_API_TOKEN` → keyring `api_token` (single shared key; OS keyring, never a plaintext file) |
| Default Project | `JIRA_DEFAULT_PROJECT` → shared internal `atk-jira.default_project` → legacy |
| Auth Method | `JIRA_AUTH_METHOD` → `ATLASSIAN_AUTH_METHOD` → shared `default` → legacy → `basic` |
| Cloud ID | `JIRA_CLOUD_ID` → `ATLASSIAN_CLOUD_ID` → shared `default` → legacy |

Per §2.2 connection config is single-sourced from the shared `default` section — internal per-tool `atk_cfl:`/`atk_jira:` sections carry only non-secret defaults and may not override `url`/`email`/`auth_method`/`cloud_id`.

**Shared credentials:** If you use both `atk-jira` and `atk-cfl` (Confluence CLI), set `ATLASSIAN_*` variables once:

```bash
export ATLASSIAN_URL=https://mycompany.atlassian.net
export ATLASSIAN_EMAIL=user@example.com
export ATLASSIAN_API_TOKEN=your-api-token
```

**Per-tool override:** Use `JIRA_*` to override for Jira specifically:

```bash
export ATLASSIAN_EMAIL=user@example.com
export ATLASSIAN_API_TOKEN=your-api-token
export JIRA_URL=https://jira.internal.corp.com  # Different URL for Jira
```

> **Note:** The legacy `JIRA_DOMAIN` environment variable is still supported for backwards compatibility but is deprecated.

---

## Shell Completion

atk-jira supports tab completion for bash, zsh, fish, and PowerShell.

### Bash

```bash
# Load in current session
source <(atk-jira completion bash)

# Install permanently (Linux)
atk-jira completion bash | sudo tee /etc/bash_completion.d/atk-jira > /dev/null

# Install permanently (macOS with Homebrew)
atk-jira completion bash > $(brew --prefix)/etc/bash_completion.d/atk-jira
```

### Zsh

```bash
# Load in current session
source <(atk-jira completion zsh)

# Install permanently
mkdir -p ~/.zsh/completions
atk-jira completion zsh > ~/.zsh/completions/_atk-jira

# Add to ~/.zshrc if not already present:
# fpath=(~/.zsh/completions $fpath)
# autoload -Uz compinit && compinit
```

### Fish

```bash
# Load in current session
atk-jira completion fish | source

# Install permanently
atk-jira completion fish > ~/.config/fish/completions/atk-jira.fish
```

### PowerShell

```powershell
# Load in current session
atk-jira completion powershell | Out-String | Invoke-Expression

# Install permanently (add to $PROFILE)
atk-jira completion powershell >> $PROFILE
```

---

## Development

### Prerequisites

- Go 1.24 or later
- golangci-lint (for linting)

### Build

```bash
make build
```

### Test

```bash
make test
```

### Lint

```bash
make lint
```

---

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

Adding a new command or flag? Read these specs first — they're the contract every command in this CLI is held to:

- [internal/cmd/GUARDRAILS.md](internal/cmd/GUARDRAILS.md) — verb language, flag aliases, pagination, mutation safety, boolean conventions, positional-vs-flag rule
- [internal/cmd/OUTPUT_SPEC.md](internal/cmd/OUTPUT_SPEC.md) — list/get/mutation output shapes, `--id` / `--extended` / `--fulltext` semantics, error conventions

## License

MIT License - see [LICENSE](LICENSE) for details.
