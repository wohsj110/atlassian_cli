# 命令面兼容规格

[English](COMPATIBILITY.md)

本项目以本文档定义的命令面为兼容目标，同时保留本项目的公开二进制命名：

| 产品域 | 本项目 |
|---|---|
| Jira CLI | `atk-jira` |
| Confluence CLI | `atk-cfl` |

目标是在 `atk-*` 命令下提供稳定的操作面、安全模型、输出约束和测试策略。

## 兼容规则

1. resource 名、action 名、位置参数含义和 flag 名应保持稳定；如果变更，需要文档说明。
2. Jira 命令形态：`atk-jira [resource] [action] [KEY/ID] [flags]`。
3. Confluence 命令形态：`atk-cfl [resource] [action] [ID] [flags]`。
4. Jira 文本输出应收敛到 text-first、pipe-delimited、agent-friendly 契约。
5. Confluence 输出应收敛到本文档定义的 `-o table|plain`、`--full`、`--raw` 和 page-content 契约。
6. 破坏性操作默认要求确认，并使用 `--force` 跳过确认。
7. 非破坏性写操作在可逆或自然追加时可以不提示；本项目可以额外保留 `--dry-run` 作为安全能力。
8. API 行为应跟随当前 Atlassian Cloud API，避免依赖过时 endpoint。

## Jira 命令面目标

### 当前用户

- `atk-jira me`
- `atk-jira me --id`
- `atk-jira me --extended`

### Issues

- `atk-jira issues list --project KEY`
- `atk-jira issues get PROJ-123`
- `atk-jira issues create --project KEY --type TYPE --summary "TEXT"`
- `atk-jira issues update PROJ-123 --field FIELD=VALUE`
- `atk-jira issues search --jql "JQL"`
- `atk-jira issues assign PROJ-123 VALUE`
- `atk-jira issues assign PROJ-123 --unassign`
- `atk-jira issues move PROJ-1 [PROJ-2 ...] --to-project NEWPROJ`
- `atk-jira issues move-status TASK_ID`
- `atk-jira issues delete PROJ-123`
- `atk-jira issues types --project KEY`
- `atk-jira issues fields [PROJ-123]`
- `atk-jira issues field-options FIELD_NAME_OR_ID [--issue PROJ-123]`

迁移期兼容 alias：

- `atk-jira issue get KEY` 映射到 `atk-jira issues get KEY`。
- `atk-jira issue search QUERY` 映射到 `atk-jira issues search --jql QUERY`。
- `atk-jira issue comment KEY --body TEXT` 映射到 `atk-jira comments add KEY --body TEXT`。

### Transitions

- `atk-jira transitions list PROJ-123`
- `atk-jira transitions list PROJ-123 --extended`
- `atk-jira transitions do PROJ-123 "Transition Name"`
- `atk-jira transitions do PROJ-123 21`
- `atk-jira transitions do PROJ-123 "Done" --field NAME=VALUE`

### Projects

- `atk-jira projects list`
- `atk-jira projects get KEY`

### Sprints

- `atk-jira sprints list --board ID_OR_NAME`
- `atk-jira sprints current --board ID_OR_NAME`
- `atk-jira sprints issues SPRINT_ID_OR_NAME`
- `atk-jira sprints add SPRINT_ID_OR_NAME PROJ-1 PROJ-2 ...`

### Boards

- `atk-jira boards list`
- `atk-jira boards get ID_OR_NAME`

### Links

- `atk-jira links list PROJ-123`
- `atk-jira links create PROJ-123 PROJ-456 --type NAME`
- `atk-jira links delete LINK_ID`
- `atk-jira links types`

### Comments

- `atk-jira comments list PROJ-123`
- `atk-jira comments add PROJ-123 --body "TEXT"`
- `atk-jira comments delete PROJ-123 COMMENT_ID`

### Attachments

- `atk-jira attachments list PROJ-123`
- `atk-jira attachments add PROJ-123 --file PATH`
- `atk-jira attachments get ATTACHMENT_ID`
- `atk-jira attachments delete ATTACHMENT_ID`

### Users

- `atk-jira users search "QUERY"`
- `atk-jira users get ACCOUNT_ID`

### 管理类命令面

- `atk-jira fields ...`
- `atk-jira dashboards ...`
- `atk-jira automation ...`
- `atk-jira refresh ...`
- `atk-jira config ...`
- `atk-jira set-credential ...`
- `atk-jira completion ...`

## Confluence 命令面目标

### 当前用户

- `atk-cfl me`
- `atk-cfl me --id`

### Pages

- `atk-cfl page list --space KEY`
- `atk-cfl page view PAGE_ID`
- `atk-cfl page view PAGE_ID --no-truncate`
- `atk-cfl page view PAGE_ID --content-only`
- `atk-cfl page view PAGE_ID --raw`
- `atk-cfl page view PAGE_ID --show-macros`
- `atk-cfl page view PAGE_ID --web`
- `atk-cfl page create --space KEY --title "TEXT"`
- `atk-cfl page create --space KEY --title "TEXT" --file content.md`
- `atk-cfl page edit PAGE_ID`
- `atk-cfl page edit PAGE_ID --file content.md`
- `atk-cfl page copy PAGE_ID --title "Copy Title"`
- `atk-cfl page delete PAGE_ID`
- `atk-cfl page delete PAGE_ID --force`

迁移期兼容 alias：

- `atk-cfl page get ID` 映射到 `atk-cfl page view ID`。
- `atk-cfl page search QUERY` 映射到 `atk-cfl search QUERY`。
- `atk-cfl page update ID` 映射到 `atk-cfl page edit ID`。

### Search

- `atk-cfl search "query"`
- `atk-cfl search "query" --space KEY`
- `atk-cfl search "query" --type page`
- `atk-cfl search --label TAG`
- `atk-cfl search --title "TEXT"`
- `atk-cfl search --cql "CQL_QUERY"`

### Spaces

- `atk-cfl space list`
- `atk-cfl space view KEY`
- `atk-cfl space create --key KEY --name "NAME"`
- `atk-cfl space update KEY --name "NAME"`
- `atk-cfl space delete KEY`

### Attachments

- `atk-cfl attachment list --page PAGE_ID`
- `atk-cfl attachment upload --page PAGE_ID --file PATH`
- `atk-cfl attachment download ATT_ID`
- `atk-cfl attachment delete ATT_ID`

### 控制面

- `atk-cfl init`
- `atk-cfl config show`
- `atk-cfl config test`
- `atk-cfl config clear`
- `atk-cfl set-credential`
- `atk-cfl completion`

## 环境变量兼容

优先级：

```text
tool-specific env -> ATLASSIAN_* env -> config file -> defaults
```

Jira：

- `JIRA_URL` / `JIRA_BASE_URL`
- `JIRA_EMAIL`
- `JIRA_API_TOKEN` / `JIRA_TOKEN`
- `JIRA_AUTH_METHOD` / `JIRA_AUTH_TYPE`
- `JIRA_CLOUD_ID`
- `JIRA_DEFAULT_PROJECT`

Confluence：

- `CFL_URL` / `CONFLUENCE_URL`
- `CFL_EMAIL` / `CONFLUENCE_EMAIL`
- `CFL_API_TOKEN` / `CONFLUENCE_API_TOKEN`
- `CFL_AUTH_METHOD` / `CONFLUENCE_AUTH_METHOD`
- `CFL_CLOUD_ID` / `CONFLUENCE_CLOUD_ID`

共享：

- `ATLASSIAN_URL` / `ATLASSIAN_SITE`
- `ATLASSIAN_EMAIL`
- `ATLASSIAN_API_TOKEN` / `ATLASSIAN_TOKEN`
- `ATLASSIAN_AUTH_METHOD` / `ATLASSIAN_AUTH_TYPE`
- `ATLASSIAN_CLOUD_ID`

## 输出兼容

Jira 输出目标：

- 默认：稳定文本。
- 列表：带大写 header 的 pipe-delimited rows。
- 详情：key-value blocks。
- `--id`：只输出主标识符。
- `--extended`：增加 admin/schema/audit 字段。
- `--fulltext`：关闭截断。
- JSON 不是通用 Jira 输出模式，只在 round-trip payload 需要时使用。

Confluence 输出目标：

- `-o table` 为默认。
- `-o plain` 用于脚本化文本处理。
- `-o json` 不是全局资源输出模式。
- `--full` 增加检查字段。
- `--raw` 是命令级能力，主要用于保真页面内容。

这是兼容目标；实现层仍可能在迁移期间保留早期命令的 `--json`。
