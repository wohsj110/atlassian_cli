# Atlassian Agent CLI

[English](README.md) | [发布指南](docs/RELEASING.zh-CN.md) | [兼容规格](docs/COMPATIBILITY.zh-CN.md) | [测试策略](docs/TEST_STRATEGY.zh-CN.md)

Atlassian Agent CLI 是一套面向人类和 AI agent 的 Jira / Confluence 命令行工具。

当前项目提供两个二进制命令：

- `atk-jira`：Jira ticket CLI。
- `atk-cfl`：Confluence page CLI。

公开包名保持更清晰的长名字，比如 `atlassian-agent-cli`；二进制命令使用紧凑但带命名空间的 `atk-*` 形式，降低和用户本机已有命令冲突的概率。

目标命令面见 [docs/COMPATIBILITY.zh-CN.md](docs/COMPATIBILITY.zh-CN.md)。公开二进制命令是 Jira 的 `atk-jira` 和 Confluence 的 `atk-cfl`。

## 当前状态

当前代码实现了完整 Jira 和 Confluence 命令面。公开二进制命令是 `atk-jira` 和 `atk-cfl`。

已经实现：

- 完整 Jira 命令面：issues、comments、transitions、projects、sprints、boards、links、remote links、attachments、users、fields、dashboards、automation、config、completion、credential ingress。
- 完整 Confluence 命令面：page、search、space、attachment、config、completion、credential ingress。
- 共享 Atlassian client、keyring/credential store、输出 presenter，以及命令面兼容测试。
- 面向 `atk-jira` 和 `atk-cfl` 的 GoReleaser 构建配置。

## 使用 Homebrew 安装

正式发布后，可以用 Homebrew 安装 CLI 二进制命令：

```bash
brew install --cask wohsj110/tap/atk-jira
brew install --cask wohsj110/tap/atk-cfl
```

也可以先 tap，再安装：

```bash
brew tap wohsj110/tap
brew install --cask atk-jira
brew install --cask atk-cfl
```

验证命令：

```bash
atk-jira --help
atk-cfl --help
```

如果 Homebrew 还搜不到包，说明 release tag 还没有发布到 tap。发布步骤见 [docs/RELEASING.zh-CN.md](docs/RELEASING.zh-CN.md)。

## 本地开发安装

```bash
go install ./tools/atk-jira/cmd/atk-jira ./tools/atk-cfl/cmd/atk-cfl
export PATH="$(go env GOPATH)/bin:$PATH"
```

检查命令：

```bash
atk-jira --help
atk-cfl --help
```

## 配置

先创建 Atlassian API token，然后导出环境变量或直接传给 `init`。

```bash
export ATLASSIAN_API_TOKEN="your-api-token"
```

配置 Jira：

```bash
atk-jira init \
  --url https://example.atlassian.net \
  --email user@example.com \
  --token-stdin < <(printf %s "$ATLASSIAN_API_TOKEN")
```

配置 Confluence：

```bash
atk-cfl init \
  --url https://example.atlassian.net \
  --email user@example.com \
  --token-stdin < <(printf %s "$ATLASSIAN_API_TOKEN")
```

配置会写入共享 Atlassian command-line toolkit config 和系统 OS keyring。Jira 专用值也可以用 `JIRA_*` 覆盖；Confluence 专用值也可以用 `CFL_*` 覆盖。

可用环境变量覆盖配置：

- `ATLASSIAN_URL`
- `ATLASSIAN_EMAIL`
- `ATLASSIAN_AUTH_METHOD`
- `ATLASSIAN_API_TOKEN`

## Jira 示例

```bash
atk-jira me --json
atk-jira issues get PROJ-123
atk-jira issues get PROJ-123 --json
atk-jira issues get PROJ-123 --id
atk-jira issues search --jql "project = PROJ ORDER BY updated DESC" --json
```

## Confluence 示例

```bash
atk-cfl me
atk-cfl me --id
atk-cfl search "release plan"
atk-cfl search "release plan" -o plain
atk-cfl page view 123456
atk-cfl page view 123456 --content-only
```

## Agent Skills

最推荐使用开放的 `skills` CLI。你可以按 agent 和 skill 分开安装。

### Codex

只给 Codex 安装 Jira 支持：

```bash
npx skills add https://github.com/wohsj110/atlassian_cli \
  --skill atk-jira \
  --agent codex \
  --global \
  --yes
```

只给 Codex 安装 Confluence 支持：

```bash
npx skills add https://github.com/wohsj110/atlassian_cli \
  --skill atk-cfl \
  --agent codex \
  --global \
  --yes
```

### Claude Code

只给 Claude Code 安装 Jira 支持：

```bash
npx skills add https://github.com/wohsj110/atlassian_cli \
  --skill atk-jira \
  --agent claude-code \
  --global \
  --yes
```

只给 Claude Code 安装 Confluence 支持：

```bash
npx skills add https://github.com/wohsj110/atlassian_cli \
  --skill atk-cfl \
  --agent claude-code \
  --global \
  --yes
```

### 一次安装全部

同时给 Codex 和 Claude Code 安装 Jira / Confluence 两个 skills：

```bash
npx skills add https://github.com/wohsj110/atlassian_cli \
  --skill atk-jira \
  --skill atk-cfl \
  --agent codex \
  --agent claude-code \
  --global \
  --yes
```

安装前查看仓库内可用 skills：

```bash
npx skills add https://github.com/wohsj110/atlassian_cli --list
```

skills.sh 完成索引后，稳定 skill ID 是：

- `wohsj110/atlassian_cli/atk-jira`
- `wohsj110/atlassian_cli/atk-cfl`

开放 `skills` CLI 使用的安装目标：

- Codex：`~/.agents/skills`
- Claude Code：`~/.claude/skills`

安装后的 skill 自身会指导 agent 检查 `atk-jira` / `atk-cfl` 是否存在，并在缺失时通过 Homebrew 或 npm helper 安装 CLI。

### 可选 npm helper

本项目也提供 npm helper，用于安装两个 skills 并检查 CLI：

```bash
npx @wohsj110/atlassian-agent-skill add atlassian-agent
npx @wohsj110/atlassian-agent-skill add atlassian-agent --install-cli
npx @wohsj110/atlassian-agent-skill doctor
```

## 输出契约

命令输出需要同时适合人类阅读和 agent 解析：

- 默认输出是简洁 table 文本。
- `--json` 输出结构化 JSON。
- `--id` 在适用命令中只输出主标识符。
- `--verbose` 预留给请求诊断。
- `--no-color` 预留给纯文本终端输出。

长内容必须通过显式参数启用，比如 `--body`、`--comments` 或 `--full`。

## 开发

运行测试：

```bash
go test ./shared/... ./tools/atk-jira/... ./tools/atk-cfl/...
npm test --prefix npm/skill-installer
```

构建两个二进制：

```bash
make build
```

不安装直接运行：

```bash
go run ./tools/atk-jira/cmd/atk-jira --help
go run ./tools/atk-cfl/cmd/atk-cfl --help
```

## 项目结构

```text
shared/             可复用 auth、config、HTTP client、output、error、credential 边界。
tools/atk-jira/          Jira CLI command、API、query、view 代码。
tools/atk-cfl/          Confluence CLI command、API、query、view 代码。
skills/             面向 agent 的 Jira 和 Confluence skill。
npm/skill-installer 安装 agent skills 的 npm 包。
docs/               发布和项目文档。
```

## Attribution

本项目是一套独立的 Atlassian 命令行工具。见 [NOTICE.md](NOTICE.md)。

## 发布

见 [docs/RELEASING.zh-CN.md](docs/RELEASING.zh-CN.md)。
