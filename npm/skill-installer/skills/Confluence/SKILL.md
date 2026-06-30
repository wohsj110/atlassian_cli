# Confluence CLI Skill

Use this skill when working with Confluence pages through `atk-cfl`.

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
