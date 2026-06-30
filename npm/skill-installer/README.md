# Atlassian Agent Skill

Install the `atk-jira` and `atk-cfl` agent skills for Codex and Claude Code.

## Quick Install

Preferred install through the open skills CLI:

```bash
npx skills add https://github.com/wohsj110/atlassian_cli \
  --skill atk-jira \
  --skill atk-cfl \
  --agent codex \
  --agent claude-code \
  --global \
  --yes
```

Project-specific npm helper:

```bash
npx @wohsj110/atlassian-agent-skill add atlassian-agent
```

This installs:

- `atk-jira` into `~/.codex/skills` and `~/.claude/skills`
- `atk-cfl` into `~/.codex/skills` and `~/.claude/skills`

## Install Skills and CLI

```bash
npx @wohsj110/atlassian-agent-skill add atlassian-agent --install-cli
```

The CLI installer uses Homebrew when available:

```bash
brew tap wohsj110/tap
brew install --cask atk-jira
brew install --cask atk-cfl
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
