# 发布指南

[English](RELEASING.md)

这份文档说明如何把 Atlassian Agent CLI 发布到 GitHub Releases、Homebrew 和 npm。

## 命名建议

推荐公开命名：

- GitHub 仓库：`wohsj110/atlassian_cli`
- Homebrew cask 或 formula：`atk-jira`、`atk-cfl`
- npm 包：`atlassian-agent-skill` 或 `@<owner>/atlassian-agent-skill`
- 二进制命令：`atk-jira`、`atk-cfl`

`atk-*` 是当前计划使用的公开命令名。它比 `atlassian-jira` / `atlassian-confluence` 短，同时比泛用缩写更不容易冲突。

## 发布前检查

运行完整本地验证：

```bash
go test ./...
npm test --prefix npm/skill-installer
goreleaser check
goreleaser release --snapshot --clean
```

如果没有安装 `goreleaser`：

```bash
brew install goreleaser
```

## 创建 GitHub 仓库

```bash
git init
git add .
git commit -m "feat: initialize atlassian agent cli"
git branch -M main
git remote add origin git@github.com:wohsj110/atlassian_cli.git
git push -u origin main
```

把 `<owner>` 替换成你的 GitHub 用户名或组织名。

## GitHub Releases

当前 `.goreleaser.yml` 会构建：

- `atk-jira`
- `atk-cfl`

release artifact 按工具拆分，方便 Homebrew、Winget、Chocolatey 和直接下载按需安装单个二进制：

- `atk-jira_<version>_<os>_<arch>`
- `atk-cfl_<version>_<os>_<arch>`

当前包管理模板期望统一的版本 tag：

- `v<version>`

创建并推送发布 tag：

```bash
git tag v0.1.0
git push origin v0.1.0
```

如果要自动发布，添加一个 GitHub Actions workflow，让 GoReleaser 在 tag push 时运行。

触发条件示例：

```yaml
on:
  push:
    tags:
      - "v*"
```

workflow 需要 `contents: write` 权限，并使用可以创建 release 的 token。

## Homebrew

推荐安装形态：

```bash
brew install --cask <owner>/tap/atk-jira
brew install --cask <owner>/tap/atk-cfl
```

这两个包分别安装公开命令：

```bash
atk-jira
atk-cfl
```

先创建 tap 仓库：

```text
github.com/wohsj110/homebrew-tap
```

然后在主仓库 secrets 里添加：

```text
HOMEBREW_TAP_GITHUB_TOKEN
```

这个 token 必须能写入 `wohsj110/homebrew-tap`。

当前 `.goreleaser.yml` 使用 `skip_upload: true` 渲染 Homebrew cask。这样本地
GoReleaser check 不会写 tap。正式 release workflow 需要把渲染出来的 `atk-jira`
和 `atk-cfl` cask 提交到 tap；或者在 CI 里移除 `skip_upload: true` 并提供 tap token：

```yaml
homebrew_casks:
  - name: atk-jira
    skip_upload: false
    repository:
      owner: <owner>
      name: homebrew-tap
      branch: main
      token: "{{ .Env.HOMEBREW_TAP_GITHUB_TOKEN }}"
  - name: atk-cfl
    skip_upload: false
    repository:
      owner: <owner>
      name: homebrew-tap
      branch: main
      token: "{{ .Env.HOMEBREW_TAP_GITHUB_TOKEN }}"
```

每个 cask 安装一个二进制：

```yaml
    binaries:
      - atk-jira
```

先跑 snapshot：

```bash
goreleaser release --snapshot --clean
```

再从 tag 发布：

```bash
git tag v0.1.0
git push origin v0.1.0
```

发布后验证：

```bash
brew tap <owner>/tap
brew install --cask atk-jira
brew install --cask atk-cfl
atk-jira --help
atk-cfl --help
```

## npm Skill Installer

npm 包在 `npm/skill-installer`。

先检查发布内容：

```bash
cd npm/skill-installer
npm test
npm pack --dry-run
npm publish --dry-run
```

正式发布：

```bash
npm login
npm publish
```

如果发布 scoped package：

```bash
npm publish --access public
```

发布后安装：

```bash
npx atlassian-agent-skill install
```

如果是 scoped package：

```bash
npx @<owner>/atlassian-agent-skill install
```

## 版本策略

默认让 CLI 和 skill installer 保持同版本，除非后续有明确理由拆分。

每次发布：

1. 更新 `npm/skill-installer/package.json`。
2. 提交版本变更。
3. 打 tag，例如 `v0.1.0`。
4. 推送 tag。
5. 发布 npm 包。

## 发布 Checklist

- [ ] `go test ./...` 通过。
- [ ] `npm test --prefix npm/skill-installer` 通过。
- [ ] `goreleaser check` 通过。
- [ ] `goreleaser release --snapshot --clean` 通过。
- [ ] GitHub 主仓库已创建。
- [ ] Homebrew tap 仓库已创建。
- [ ] 如果发布 Homebrew，`HOMEBREW_TAP_GITHUB_TOKEN` secret 已配置。
- [ ] npm 包名可用。
- [ ] tag 已推送。
- [ ] GitHub Release 已创建。
- [ ] Homebrew 安装已验证。
- [ ] npm 安装已验证。
