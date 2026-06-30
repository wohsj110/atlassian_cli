# Atlassian Agent CLI

[简体中文](README.zh-CN.md) | [Release guide](docs/RELEASING.md) | [Compatibility](docs/COMPATIBILITY.md) | [Test strategy](docs/TEST_STRATEGY.md)

Atlassian Agent CLI is a Jira and Confluence command-line toolkit designed for both humans and AI agents.

The project currently ships two binaries:

- `atk-jira`: Jira ticket CLI.
- `atk-cfl`: Confluence page CLI.

The public package name stays longer and clearer, such as `atlassian-agent-cli`, while the binaries use the compact but namespaced `atk-*` form to reduce command-name collisions.

The target command surface is documented in [docs/COMPATIBILITY.md](docs/COMPATIBILITY.md). The public binaries are `atk-jira` for Jira and `atk-cfl` for Confluence.

## Status

This repository implements the full Jira and Confluence command surfaces for the public `atk-jira` and `atk-cfl` binaries.

Implemented:

- Full Jira command surface: issues, comments, transitions, projects, sprints, boards, links, remote links, attachments, users, fields, dashboards, automation, config, completion, and credential ingress.
- Full Confluence command surface: page, search, space, attachment, config, completion, and credential ingress.
- Shared Atlassian client, keyring/credential store, output presenters, and command-surface parity tests.
- GoReleaser build configuration for `atk-jira` and `atk-cfl`.

## Install with Homebrew

After a release is published, install the CLI binaries with Homebrew:

```bash
brew install --cask wohsj110/tap/atk-jira
brew install --cask wohsj110/tap/atk-cfl
```

Or tap first:

```bash
brew tap wohsj110/tap
brew install --cask atk-jira
brew install --cask atk-cfl
```

Verify the installed commands:

```bash
atk-jira --help
atk-cfl --help
```

If Homebrew cannot find the packages, the release tag has not been published to the tap yet. See [docs/RELEASING.md](docs/RELEASING.md).

## Install for Local Development

```bash
go install ./tools/atk-jira/cmd/atk-jira ./tools/atk-cfl/cmd/atk-cfl
export PATH="$(go env GOPATH)/bin:$PATH"
```

Check the commands:

```bash
atk-jira --help
atk-cfl --help
```

## Configure

Create an Atlassian API token first, then either export it or pass it directly.

```bash
export ATLASSIAN_API_TOKEN="your-api-token"
```

Configure Jira:

```bash
atk-jira init \
  --url https://example.atlassian.net \
  --email user@example.com \
  --token-stdin < <(printf %s "$ATLASSIAN_API_TOKEN")
```

Configure Confluence:

```bash
atk-cfl init \
  --url https://example.atlassian.net \
  --email user@example.com \
  --token-stdin < <(printf %s "$ATLASSIAN_API_TOKEN")
```

Configuration is stored in the shared Atlassian command-line toolkit config and the OS keyring. Jira-specific values can also be overridden with `JIRA_*`; Confluence-specific values can be overridden with `CFL_*`.

You can override values with environment variables:

- `ATLASSIAN_URL`
- `ATLASSIAN_EMAIL`
- `ATLASSIAN_AUTH_METHOD`
- `ATLASSIAN_API_TOKEN`

## Jira Examples

```bash
atk-jira me --json
atk-jira issues get PROJ-123
atk-jira issues get PROJ-123 --json
atk-jira issues get PROJ-123 --id
atk-jira issues search --jql "project = PROJ ORDER BY updated DESC" --json
```

## Confluence Examples

```bash
atk-cfl me
atk-cfl me --id
atk-cfl search "release plan"
atk-cfl search "release plan" -o plain
atk-cfl page view 123456
atk-cfl page view 123456 --content-only
```

## Agent Skills

The easiest install path is the open `skills` CLI. Pick the agent and skill you need.

### Codex

Install only Jira support for Codex:

```bash
npx skills add https://github.com/wohsj110/atlassian_cli \
  --skill atk-jira \
  --agent codex \
  --global \
  --yes
```

Install only Confluence support for Codex:

```bash
npx skills add https://github.com/wohsj110/atlassian_cli \
  --skill atk-cfl \
  --agent codex \
  --global \
  --yes
```

### Claude Code

Install only Jira support for Claude Code:

```bash
npx skills add https://github.com/wohsj110/atlassian_cli \
  --skill atk-jira \
  --agent claude-code \
  --global \
  --yes
```

Install only Confluence support for Claude Code:

```bash
npx skills add https://github.com/wohsj110/atlassian_cli \
  --skill atk-cfl \
  --agent claude-code \
  --global \
  --yes
```

### Install Everything

Install both skills for both Codex and Claude Code:

```bash
npx skills add https://github.com/wohsj110/atlassian_cli \
  --skill atk-jira \
  --skill atk-cfl \
  --agent codex \
  --agent claude-code \
  --global \
  --yes
```

List the available skills before installing:

```bash
npx skills add https://github.com/wohsj110/atlassian_cli --list
```

After skills.sh indexes this repository, the stable skill IDs are:

- `wohsj110/atlassian_cli/atk-jira`
- `wohsj110/atlassian_cli/atk-cfl`

Install targets used by the open `skills` CLI:

- Codex: `~/.agents/skills`
- Claude Code: `~/.claude/skills`

Each installed skill can help the agent check for `atk-jira` / `atk-cfl` and install the missing CLI with Homebrew or the npm helper.

### Optional Helper

This project also publishes an npm helper for installing both skills and checking the CLI:

```bash
npx @wohsj110/atlassian-agent-skill add atlassian-agent
npx @wohsj110/atlassian-agent-skill add atlassian-agent --install-cli
npx @wohsj110/atlassian-agent-skill doctor
```

## Output Contract

Commands are designed to be predictable for humans and agents:

- Default output is concise table text.
- `--json` returns structured output.
- `--id` returns the primary identifier when a command has one.
- `--verbose` is reserved for request diagnostics.
- `--no-color` is accepted by command surfaces that may eventually color output.

Long payloads should require explicit flags such as `--body`, `--comments`, or `--full`.

## Development

Run tests:

```bash
go test ./shared/... ./tools/atk-jira/... ./tools/atk-cfl/...
npm test --prefix npm/skill-installer
```

Build both binaries:

```bash
make build
```

Run without installing:

```bash
go run ./tools/atk-jira/cmd/atk-jira --help
go run ./tools/atk-cfl/cmd/atk-cfl --help
```

## Repository Layout

```text
shared/             Reusable auth, config, HTTP client, output, error, and credential boundaries.
tools/atk-jira/          Jira CLI command, API, query, and view code.
tools/atk-cfl/          Confluence CLI command, API, query, and view code.
skills/             Agent-facing Jira and Confluence skill instructions.
npm/skill-installer npm package that installs skills into agent skill directories.
docs/               Release and project documentation.
```

## Attribution

This project is a standalone Atlassian command-line toolkit. See [NOTICE.md](NOTICE.md).

## Release

See [docs/RELEASING.md](docs/RELEASING.md).
