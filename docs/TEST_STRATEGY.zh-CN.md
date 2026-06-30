# 测试策略

[English](TEST_STRATEGY.md)

测试策略以本项目文档化的命令面为目标，并适配 `atk-jira` 和 `atk-cfl` 二进制命令。

## 测试原则

1. 先测试命令行为，再实现。
2. API 行为优先使用 mocked HTTP server。
3. 真实 Atlassian 集成测试必须放在显式环境变量开关后面。
4. 默认 CI 不依赖真实凭据。
5. 写命令必须证明安全行为：
   - `--dry-run` 不触网。
   - 破坏性命令默认要求确认。
   - `--force` 跳过破坏性确认。
   - `--yes` 只用于本项目明确保留额外写确认的地方。
6. 校验 stdout/stderr 边界：
   - 数据写 stdout。
   - 诊断和 prompt 写 stderr。
7. 校验 config、错误、verbose output 和测试日志不泄露 secret。

## 测试层级

### Shared Packages

`shared/auth`

- Basic auth header 编码 `email:token`。
- Bearer auth header 输出 `Bearer <token>`。
- 非法 auth method 应在触网前失败。

`shared/config`

- tool-specific env 覆盖 `ATLASSIAN_*`。
- `ATLASSIAN_*` 覆盖 config。
- config 文件只加载非 secret 字段。
- token 必须进入 credential store，不能写入明文配置。

`shared/client`

- 相对和绝对 URL 正确解析。
- POST/PUT JSON body 正确 marshal。
- `Accept` 和 `Content-Type` header 正确设置。
- auth header 会被发送。
- 4xx/5xx 正确解析 Jira 和 Confluence 错误结构。
- verbose output 会 redact secret 并截断大 body。
- rate-limit 和 retry 行为实现后需要覆盖。

`shared/output`

- table rendering 稳定。
- plain rendering 稳定。
- JSON rendering 只用于明确的 JSON contract。
- 兼容契约要求时，空值输出为 `-`。

`shared/adf`

- plain text 转 ADF document。
- Markdown 转 ADF document。
- ADF 转 plain text，用于 comments、descriptions、page bodies。
- 空输入行为稳定。

### Jira 命令测试

命令解析：

- `atk-jira --help` 列出顶层 resource。
- 每个 resource help 列出目标命令兼容 action。
- 全局 flag 在目标命令允许的位置都能解析。
- 非法 flag 在触网前失败。

读能力：

- `me`
- `issues list/get/search`
- `comments list`
- `transitions list`
- `projects list/get`
- `users search/get`

写能力：

- `comments add` 发送 ADF body。
- `issues create/update/assign/transition` 发送预期 payload。
- 破坏性命令没有 `--force` 时提示确认。
- mutation 根据 flag 输出 post-state 或 identifier。

API 兼容：

- Jira search 使用 `/rest/api/3/search/jql`。
- Comments 使用 `/rest/api/3/issue/{key}/comment`。
- Issue get 使用 `/rest/api/3/issue/{key}`。
- Transition 使用 Jira transition API；必要时先读取可用 transition。

### Confluence 命令测试

命令解析：

- `atk-cfl --help` 列出顶层 resource。
- `page view`、`search`、`space`、`attachment`、`config`、`set-credential`、`completion` 命令面匹配兼容规格。

读能力：

- `me`
- `page list/view`
- `search`
- `space list/view`
- `attachment list`

写能力：

- `page create/edit/copy/delete`
- `space create/update/delete`
- `attachment upload/delete`

内容行为：

- `page view` 默认输出 markdown。
- `--raw` 返回 storage/source content。
- `--content-only` 省略 metadata。
- `--no-truncate` 关闭截断保护。

## 集成测试

集成测试默认关闭。

建议开关：

```bash
ATK_JIRA_INTEGRATION=1 go test ./tools/atk-jira/... -run Integration
ATK_CONFLUENCE_INTEGRATION=1 go test ./tools/atk-cfl/... -run Integration
```

所需环境变量：

Jira：

- `JIRA_URL` 或 `JIRA_BASE_URL`
- `JIRA_EMAIL`
- `JIRA_API_TOKEN`
- 可选 `JIRA_TEST_PROJECT`
- 可选 `JIRA_TEST_ISSUE`

Confluence：

- `CFL_URL` 或 `CONFLUENCE_URL`
- `CFL_EMAIL` 或 `CONFLUENCE_EMAIL`
- `CFL_API_TOKEN` 或 `CONFLUENCE_API_TOKEN`
- 可选 `CFL_TEST_SPACE`
- 可选 `CFL_TEST_PAGE`

集成测试必须：

- 默认避免破坏性行为。
- 创建资源时使用唯一前缀。
- 尽可能清理测试资源。
- 缺少环境变量时 skip，而不是 fail。
- 永不打印 token。

## 兼容覆盖跟踪

每个命令族维护覆盖表：

| Area | Parser tests | API tests | Output tests | Integration smoke |
|---|---:|---:|---:|---:|
| Jira me | yes | yes | partial | manual |
| Jira issues | partial | partial | partial | manual |
| Jira comments | partial | partial | partial | manual |
| Jira transitions | no | no | no | no |
| Confluence page | partial | partial | partial | no |
| Confluence search | partial | partial | partial | no |

实现或补测命令时同步更新这张表。
