---
name: atk-jira
description: Use when working with Jira Cloud through the atk-jira CLI, including reading issues, searching JQL, creating/updating issues, comments, transitions, projects, sprints, boards, attachments, users, fields, dashboards, automation, and Jira credential setup.
---

# Jira CLI Skill

Use this skill when working with Jira issues through `atk-jira`.

## CLI Setup
- Before using Jira commands, check whether `atk-jira` is available with `command -v atk-jira`.
- If missing, install the CLI with `brew install --cask wohsj110/tap/atk-jira`.
- After installation, run `atk-jira --help` to verify the binary.
- If credentials are not configured, run `atk-jira init` or ask the user for the non-secret URL/email setup preference. Never ask the user to paste secrets into chat; use stdin, environment variables, or the OS keyring flow.

## Rules
- Read the current issue state before writing.
- Use `--id` when only the primary issue key is needed.
- Never include secrets in command output, logs, comments, or issue bodies.
- For write commands, use `--dry-run` first when the command supports it.

## Common Commands

```bash
atk-jira init
atk-jira me --id
atk-jira issues get PROJ-123
atk-jira issues search --jql "project = PROJ ORDER BY updated DESC"
atk-jira comments add PROJ-123 --body "Investigating this now."
```

## Output Expectations
- Default output is concise table text for humans.
- `--id` produces script-friendly identifiers when supported.
- Use `--extended` for admin/schema/audit detail and `--fulltext` to disable truncation where supported.

See `CliReference.md` for command details.
