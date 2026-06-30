---
name: atk-cfl
description: Use when working with Atlassian Confluence Cloud through the atk-cfl CLI, including searching content, reading pages, creating/editing/copying/deleting pages, managing spaces, attachments, Confluence account checks, and Confluence credential setup.
---

# Confluence CLI Skill

Use this skill when working with Confluence pages through `atk-cfl`.

## CLI Setup
- Before using Confluence commands, check whether `atk-cfl` is available with `command -v atk-cfl`.
- If missing, install the CLI with `brew install --cask wohsj110/tap/atk-cfl` or `npx @wohsj110/atlassian-agent-skill install-cli`.
- After installation, run `atk-cfl --help` to verify the binary.
- If credentials are not configured, run `atk-cfl init` or ask the user for the non-secret URL/email setup preference. Never ask the user to paste secrets into chat; use stdin, environment variables, or the OS keyring flow.

## Rules
- Read the current page state before writing.
- Prefer `-o plain` when another tool or agent will parse list/search output.
- Use `--id` when only the authenticated account ID is needed.
- Use `--content-only` when only page content is needed; default output should stay concise.
- Never include secrets in command output, logs, comments, or page bodies.
- Read the current page with `page view` before write commands.

## Common Commands

```bash
atk-cfl init
atk-cfl me --id
atk-cfl search "release plan" -o plain
atk-cfl page view 123456 --content-only
atk-cfl page edit 123456 --file updated.md
```

## Output Expectations
- Default output is concise table text for humans.
- `-o plain` produces parseable text for list/search output.
- Long fields require explicit flags such as `--content-only`, `--no-truncate`, or `--full` as commands grow.

See `CliReference.md` for command details.
