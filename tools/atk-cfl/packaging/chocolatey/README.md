# Chocolatey Package for atk-cfl

This directory contains the Chocolatey package definition for distributing atk-cfl on Windows.

## Automated Publishing

Publishing to Chocolatey is automated via GitHub Actions. When a new release tag is pushed, the release workflow:

1. Builds binaries with GoReleaser
2. Packs the Chocolatey package with the release version
3. Pushes to the Chocolatey Community Repository

**Required secret:** `CHOCOLATEY_API_KEY` must be configured in repository settings.

## Package Structure

```
packaging/chocolatey/
├── atk-cfl.nuspec    # Package manifest
├── tools/
│   ├── chocolateyInstall.ps1    # Downloads and installs from GitHub Releases
│   └── chocolateyUninstall.ps1  # Cleanup script
└── README.md
```

## Local Testing

### Prerequisites

- Windows with [Chocolatey installed](https://chocolatey.org/install)
- PowerShell (Admin)

### Build the Package

```powershell
cd packaging/chocolatey

# Update version in nuspec to match a real release (e.g., 0.10.0)
# Then pack:
choco pack
```

This creates `atk-cfl.<version>.nupkg`.

### Install Locally

```powershell
# Install from local package
choco install atk-cfl -s . --force

# Verify
atk-cfl --version

# Uninstall
choco uninstall atk-cfl
```

## Publishing to Chocolatey Community Repository

### First-Time Setup

1. Create an account at https://community.chocolatey.org
2. Get your API key from https://community.chocolatey.org/account
3. Configure your API key:
   ```powershell
   choco apikey --key <your-api-key> --source https://push.chocolatey.org/
   ```

### Publishing a New Version

1. Update the `<version>` in `atk-cfl.nuspec` to match the GitHub release
2. Pack the package:
   ```powershell
   choco pack
   ```
3. Push to Chocolatey:
   ```powershell
   choco push atk-cfl.<version>.nupkg --source https://push.chocolatey.org/
   ```

### Moderation Process

- New packages go through moderation (typically 1-3 days)
- Automated checks verify the package downloads correctly
- Human moderators review the package
- Status updates are sent via email

## Architecture Support

The install script automatically detects Windows architecture:

| Architecture | Download |
|--------------|----------|
| ARM64 | `atk-cfl_<version>_windows_arm64.zip` |
| x64 | `atk-cfl_<version>_windows_amd64.zip` |
| x86 | Not supported (error) |

## Updating for New Releases

When a new version is released on GitHub:

1. Update `<version>` in `atk-cfl.nuspec`
2. Test locally with `choco pack && choco install atk-cfl -s . --force`
3. Push to Chocolatey with `choco push`

The install script dynamically fetches checksums from `checksums.txt` in the GitHub release, so no checksum updates are needed in the package files.
