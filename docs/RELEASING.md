# Release Guide

[简体中文](RELEASING.zh-CN.md)

This guide describes how to publish Atlassian Agent CLI to GitHub Releases, Homebrew, and npm.

## Naming

Recommended public names:

- GitHub repository: `wohsj110/atlassian_cli`
- Homebrew casks or formulae: `atk-jira`, `atk-cfl`
- npm package: `atlassian-agent-skill` or `@<owner>/atlassian-agent-skill`
- Binaries: `atk-jira`, `atk-cfl`

The `atk-*` names are the intended public command names. They are shorter than `atlassian-jira` / `atlassian-confluence` while still being less collision-prone than generic abbreviations.

## Preflight

Run the full local verification:

```bash
go test ./shared/... ./tools/atk-jira/... ./tools/atk-cfl/...
npm test --prefix npm/skill-installer
goreleaser check
goreleaser release --snapshot --clean
```

If `goreleaser` is not installed:

```bash
brew install goreleaser
```

## Create the GitHub Repository

```bash
git init
git add .
git commit -m "feat: initialize atlassian agent cli"
git branch -M main
git remote add origin git@github.com:wohsj110/atlassian_cli.git
git push -u origin main
```

Replace `<owner>` with your GitHub user or organization.

## GitHub Releases

The current `.goreleaser.yml` builds release artifacts for:

- `atk-jira`
- `atk-cfl`

The release artifacts are intentionally split by tool so Homebrew, Winget,
Chocolatey, and direct downloads can install one binary at a time:

- `atk-jira_<version>_<os>_<arch>`
- `atk-cfl_<version>_<os>_<arch>`

The packaging templates expect one shared version tag:

- `v<version>`

Create and push the release tag:

```bash
git tag v0.1.0
git push origin v0.1.0
```

For automated releases, add a GitHub Actions workflow that runs GoReleaser on tags.

Example target behavior:

```yaml
on:
  push:
    tags:
      - "v*"
```

The workflow should give `contents: write` permission and run GoReleaser with a token that can create releases.

## Homebrew

Recommended install shape:

```bash
brew install --cask wohsj110/tap/atk-jira
brew install --cask wohsj110/tap/atk-cfl
```

Those packages install the public binaries:

```bash
atk-jira
atk-cfl
```

Create a tap repository:

```text
github.com/wohsj110/homebrew-tap
```

Then add a token to the main repository secrets:

```text
HOMEBREW_TAP_GITHUB_TOKEN
```

The token must be able to write to `wohsj110/homebrew-tap`.

The current `.goreleaser.yml` publishes Homebrew casks to the tap. GoReleaser
v2.16 deprecates the old `brews` formula generator for prebuilt binaries, so the
supported Homebrew path is cask-based.

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

Each cask installs one binary:

```yaml
    binaries:
      - atk-jira
```

Run a snapshot first:

```bash
goreleaser release --snapshot --clean
```

Then release from a tag:

```bash
git tag v0.1.0
git push origin v0.1.0
```

After release:

```bash
brew tap wohsj110/tap
brew install --cask atk-jira
brew install --cask atk-cfl
atk-jira --help
atk-cfl --help
```

## npm Skill Installer

The npm package lives in `npm/skill-installer`.

Check what will be published:

```bash
cd npm/skill-installer
npm test
npm pack --dry-run
npm publish --dry-run
```

Publish:

```bash
npm login
npm publish
```

If you publish a scoped package, use:

```bash
npm publish --access public
```

After publishing:

```bash
npx atlassian-agent-skill install
```

Or, for a scoped package:

```bash
npx @<owner>/atlassian-agent-skill install
```

## Versioning

Keep CLI and skill installer versions aligned unless there is a clear reason to split them.

For each release:

1. Update `npm/skill-installer/package.json`.
2. Commit the version change.
3. Tag the release, for example `v0.1.0`.
4. Push the tag.
5. Publish the npm package.

## Release Checklist

- [ ] `go test ./shared/... ./tools/atk-jira/... ./tools/atk-cfl/...` passes.
- [ ] `npm test --prefix npm/skill-installer` passes.
- [ ] `goreleaser check` passes.
- [ ] `goreleaser release --snapshot --clean` passes.
- [ ] GitHub repository exists.
- [ ] Homebrew tap repository exists.
- [ ] `HOMEBREW_TAP_GITHUB_TOKEN` secret exists if publishing Homebrew.
- [ ] npm package name is available.
- [ ] Tag pushed.
- [ ] GitHub Release created.
- [ ] Homebrew install verified.
- [ ] npm install verified.
