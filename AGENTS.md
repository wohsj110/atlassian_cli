# AGENTS.md

This repository builds a general-purpose Atlassian command-line toolkit for both humans and AI agents.

## Language
- Think in English.
- Reply to the user in Simplified Chinese unless they request otherwise.

## Product Goal
Build two CLI tools with the documented command surface, using this repository's binary names:
- `atk-jira`: Jira ticket CLI.
- `atk-cfl`: Confluence CLI.

The tools must be pleasant for humans and predictable for AI agents. This is not project-specific automation and must not depend on any single workspace.

See `docs/COMPATIBILITY.md` and `docs/TEST_STRATEGY.md` before adding commands.

## Core Principles
- Rewrite-first. Borrow structure and small implementation ideas only when useful.
- If copying code from MIT projects, keep required license attribution.
- Prefer small, stable commands over broad REST API mirrors.
- Default output must be readable, stable, concise, and safe for agents to parse.
- Every write operation must support preview or dry-run behavior where practical.
- Do not add new third-party libraries without checking existing needs first.

## Architecture
Use Go for CLI binaries.

Shared packages live under `shared/`:
- `auth`: Basic auth, bearer token auth, auth context.
- `client`: HTTP client, retries, rate-limit handling, request diagnostics.
- `config`: config file loading, environment overrides.
- `credstore`: OS keychain integration where available.
- `output`: table/text/json rendering.
- `errors`: normalized API and CLI errors.

Tool-specific code lives under:
- `tools/atk-jira` (`atk-jira` binary)
- `tools/atk-cfl` (`atk-cfl` binary)

Do not put Jira-specific behavior into shared packages unless it is truly reusable by Confluence.

## Output Contract
All commands should support:
- Default text output for humans.
- `--json` for structured output.
- `--id` where a primary identifier exists.
- `--verbose` for request diagnostics.
- `--no-color` for plain terminal output.

Default output should avoid huge payloads. Use explicit flags like `--full`, `--comments`, or `--body` for long content.

## Safety
Write commands must be explicit:
- Dangerous commands require `--yes` or an interactive confirmation.
- AI-facing skills should instruct agents to read current state before writing.
- Prefer `--dry-run` on create/update/delete commands when practical.
- Never log secrets.

## MVP Scope
Jira:
- `atk-jira init`
- `atk-jira me`
- `atk-jira issue get KEY`
- `atk-jira issue search QUERY`
- `atk-jira issue create`
- `atk-jira issue update KEY`
- `atk-jira issue comment KEY`
- `atk-jira issue transition KEY`

Confluence:
- `atk-cfl init`
- `atk-cfl me`
- `atk-cfl page view ID`
- `atk-cfl search QUERY`
- `atk-cfl page create`
- `atk-cfl page edit ID`

Skill installer:
- `npx atlassian-agent-skill install`
- Installs Jira and Confluence skills into supported agent skill directories.

## Testing
- Unit-test query builders, output rendering, config loading, and error normalization.
- Use mocked HTTP clients for API behavior.
- Do not require real Atlassian credentials in CI.
- Add integration tests only behind explicit environment variables.

## Release
- Use GoReleaser for binaries.
- Homebrew installs CLI binaries.
- npm package installs skills only.
- Keep CLI and skill installer versioned together unless there is a reason to split.

## Inspiration
Useful references may inform:
- command grouping
- view/query separation
- GoReleaser/Homebrew release patterns
- completion generation
- stable agent-oriented output ideas

Avoid:
- blindly mirroring every REST endpoint
- huge default JSON output
- project-specific workflow assumptions
- copying large external source files without MIT attribution in `NOTICE.md`
