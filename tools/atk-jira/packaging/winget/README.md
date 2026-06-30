# Winget Package for atk-jira

This directory contains the Winget manifest templates for distributing atk-jira on Windows via `winget install wohsj110.atk-jira`.

## Automated Publishing

Publishing to Winget is automated via GitHub Actions. When a new release tag is pushed, the release workflow copies these templates, substitutes version/checksum placeholders, and uses `wingetcreate submit` to submit a PR to [microsoft/winget-pkgs](https://github.com/microsoft/winget-pkgs).

**Required secret:** `WINGET_GITHUB_TOKEN` - A GitHub PAT with `public_repo` scope, needed to create PRs on microsoft/winget-pkgs.

**Note:** Unlike Chocolatey (direct publish), Winget submissions are PRs that go through Microsoft's automated validation before merging.

## Manifest Structure

```
packaging/winget/
├── wohsj110.atk-jira.yaml              # Version manifest
├── wohsj110.atk-jira.installer.yaml    # Installer manifest (URLs, checksums)
├── wohsj110.atk-jira.locale.en-US.yaml # Locale manifest (descriptions, tags)
└── README.md
```

## How Winget Works

Unlike Chocolatey (which hosts packages on their own feed), Winget manifests live in Microsoft's community repository [microsoft/winget-pkgs](https://github.com/microsoft/winget-pkgs). Publishing requires submitting a PR to that repo.

## Template Placeholders

The manifest templates use these placeholders that are replaced during CI:

| Placeholder | Replaced with |
|-------------|--------------|
| `0.0.0` | Release version (e.g., `0.1.18`) |
| `CHECKSUM_AMD64_PLACEHOLDER` | SHA256 of the x64 zip |
| `CHECKSUM_ARM64_PLACEHOLDER` | SHA256 of the arm64 zip |

URLs contain `0.0.0` in both the tag path and filename, so the version replacement handles them automatically.

## Publishing a New Version

### Option 1: Manual PR

1. **Get release info:**
   - Download URLs: `https://github.com/wohsj110/atlassian_cli/releases/download/v<VERSION>/atk-jira_<VERSION>_windows_amd64.zip`
   - SHA256 checksums from `checksums.txt` in the release

2. **Update manifests:**
   - Replace `0.0.0` with the actual version in all three YAML files
   - Replace checksum placeholders with real SHA256 values

3. **Validate manifests:**
   ```powershell
   winget validate --manifest packaging/winget/
   ```

4. **Fork and clone** [microsoft/winget-pkgs](https://github.com/microsoft/winget-pkgs)

5. **Create folder structure:**
   ```
   manifests/w/wohsj110/atk-jira/<VERSION>/
   ```

6. **Copy manifests** into the folder

7. **Submit PR** to microsoft/winget-pkgs

### Option 2: Using wingetcreate submit

```powershell
# Install wingetcreate
winget install Microsoft.WingetCreate

# Copy and update templates, then submit
wingetcreate submit --path <manifest-dir> --token <PAT>
```

## Manifest Schema

These manifests use schema version 1.10.0:
- [Version manifest schema](https://aka.ms/winget-manifest.version.1.10.0.schema.json)
- [Installer manifest schema](https://aka.ms/winget-manifest.installer.1.10.0.schema.json)
- [Locale manifest schema](https://aka.ms/winget-manifest.defaultLocale.1.10.0.schema.json)

## Installer Type

This package uses:
- `InstallerType: zip` - Our releases are zip archives
- `NestedInstallerType: portable` - Contains a standalone executable
- `PortableCommandAlias: atk-jira` - Command users type to invoke the tool

Winget extracts the zip, places `atk-jira.exe` in a managed location, and creates the command alias.

## Architecture Support

| Architecture | Installer URL Pattern |
|--------------|----------------------|
| x64 | `atk-jira_<VERSION>_windows_amd64.zip` |
| arm64 | `atk-jira_<VERSION>_windows_arm64.zip` |

## After Approval

Once the PR is merged to microsoft/winget-pkgs, users can install with:
```powershell
winget install wohsj110.atk-jira
```

## References

- [Winget Manifest Documentation](https://github.com/microsoft/winget-pkgs/tree/master/doc/manifest)
- [Submit packages to Windows Package Manager](https://learn.microsoft.com/en-us/windows/package-manager/package/repository)
- [wingetcreate tool](https://github.com/microsoft/winget-create)
