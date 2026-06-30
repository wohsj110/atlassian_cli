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

Install the bundled skills locally:

```bash
node npm/skill-installer/bin/install.js install
```

Install to a custom directory:

```bash
node npm/skill-installer/bin/install.js install --dest /path/to/skills
```

The installer copies:

- `Jira/SKILL.md`
- `Jira/CliReference.md`
- `Confluence/SKILL.md`
- `Confluence/CliReference.md`

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
go test ./...
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
