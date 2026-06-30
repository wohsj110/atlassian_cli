# Jira CLI Skill

Use this skill when working with Jira issues through `atk-jira`.

## Rules
- Read the current issue state before writing.
- Prefer `--json` when another tool or agent will parse the result.
- Use `--id` when only the primary issue key is needed.
- Never include secrets in command output, logs, comments, or issue bodies.
- For write commands, use `--dry-run` first when the command supports it.

## Common Commands

```bash
atk-jira init
atk-jira me --json
atk-jira issues get PROJ-123 --json
atk-jira issues search "project = PROJ ORDER BY updated DESC" --json
atk-jira issues comment PROJ-123 --body "Investigating this now." --dry-run
```

## Output Expectations
- Default output is concise table text for humans.
- `--json` produces stable structured data.
- Long fields require explicit flags such as `--body`, `--comments`, or `--full` as commands grow.

See `CliReference.md` for command details.
