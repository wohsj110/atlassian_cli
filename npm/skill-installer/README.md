# Atlassian Agent Skill

Install the `atk-jira` and `atk-cfl` agent skills for Codex and Claude Code.

## Quick Install

Preferred install through the open `skills` CLI.

Codex + Jira only:

```bash
npx skills add https://github.com/wohsj110/atlassian_cli --skill atk-jira --agent codex --global --yes
```

Codex + Confluence only:

```bash
npx skills add https://github.com/wohsj110/atlassian_cli --skill atk-cfl --agent codex --global --yes
```

Claude Code + Jira only:

```bash
npx skills add https://github.com/wohsj110/atlassian_cli --skill atk-jira --agent claude-code --global --yes
```

Claude Code + Confluence only:

```bash
npx skills add https://github.com/wohsj110/atlassian_cli --skill atk-cfl --agent claude-code --global --yes
```

Project-specific npm helper:

```bash
npx @wohsj110/atlassian-agent-skill add atlassian-agent
```

This helper installs:

- `atk-jira` into `~/.codex/skills` and `~/.claude/skills`
- `atk-cfl` into `~/.codex/skills` and `~/.claude/skills`

The open `npx skills` CLI installs Codex skills into `~/.agents/skills`.

## Install Skills and CLI

```bash
npx @wohsj110/atlassian-agent-skill add atlassian-agent --install-cli
```

The CLI installer uses Homebrew when available:

```bash
brew install --cask wohsj110/tap/atk-jira
brew install --cask wohsj110/tap/atk-cfl
```

## Install for One Agent

Codex only:

```bash
npx @wohsj110/atlassian-agent-skill add atlassian-agent --target codex
```

Claude Code only:

```bash
npx @wohsj110/atlassian-agent-skill add atlassian-agent --target claude
```

Custom directory:

```bash
npx @wohsj110/atlassian-agent-skill add atlassian-agent --dest /path/to/skills
```

## Doctor

```bash
npx @wohsj110/atlassian-agent-skill doctor
```

Checks whether `atk-jira`, `atk-cfl`, and the installed skill files are present.

## Configure Atlassian

After installing the CLI, run:

```bash
atk-jira init
atk-cfl init
```

Secrets are stored through the OS keyring flow. Do not paste API tokens into agent chat.
