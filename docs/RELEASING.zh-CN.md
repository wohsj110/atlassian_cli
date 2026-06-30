# 发布指南

[English](RELEASING.md)

这份文档说明如何把 Atlassian Agent CLI 发布到 GitHub Releases 和 Homebrew。

## 命名建议

推荐公开命名：

- GitHub 仓库：`wohsj110/atlassian_cli`
- Homebrew cask 或 formula：`atk-jira`、`atk-cfl`
- 二进制命令：`atk-jira`、`atk-cfl`

`atk-*` 是当前计划使用的公开命令名。它比 `atlassian-jira` / `atlassian-confluence` 短，同时比泛用缩写更不容易冲突。

## 发布前检查

运行完整本地验证：

```bash
go test ./shared/... ./tools/atk-jira/... ./tools/atk-cfl/...
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
brew install --cask wohsj110/tap/atk-jira
brew install --cask wohsj110/tap/atk-cfl
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

当前 `.goreleaser.yml` 会把 Homebrew cask 发布到 tap。GoReleaser v2.16
已经废弃面向预编译二进制的旧 `brews` formula 生成器，所以支持的 Homebrew
路径是 cask。

```yaml
homebrew_casks:
  - name: atk-jira
    repository:
      owner: wohsj110
      name: homebrew-tap
      branch: main
      token: "{{ .Env.HOMEBREW_TAP_GITHUB_TOKEN }}"
  - name: atk-cfl
    repository:
      owner: wohsj110
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
brew tap wohsj110/tap
brew install --cask atk-jira
brew install --cask atk-cfl
atk-jira --help
atk-cfl --help
```

## 版本策略

两个 CLI 二进制使用同一个版本 tag。

每次发布：

1. 提交发布前变更。
2. 打 tag，例如 `v0.1.0`。
3. 推送 tag。
4. 验证 GitHub Release assets 和 Homebrew 安装。

## 发布 Checklist

- [ ] `go test ./shared/... ./tools/atk-jira/... ./tools/atk-cfl/...` 通过。
- [ ] `goreleaser check` 通过。
- [ ] `goreleaser release --snapshot --clean` 通过。
- [ ] GitHub 主仓库已创建。
- [ ] Homebrew tap 仓库已创建。
- [ ] 如果发布 Homebrew，`HOMEBREW_TAP_GITHUB_TOKEN` secret 已配置。
- [ ] tag 已推送。
- [ ] GitHub Release 已创建。
- [ ] Homebrew 安装已验证。
