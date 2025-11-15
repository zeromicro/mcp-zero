# Release Process Documentation

**Project**: mcp-zero
**Version**: 1.0.0
**Last Updated**: November 15, 2025

## Table of Contents

1. [Release Schedule](#release-schedule)
2. [Version Numbering](#version-numbering)
3. [Pre-Release Checklist](#pre-release-checklist)
4. [Release Steps](#release-steps)
5. [Post-Release Tasks](#post-release-tasks)
6. [Hotfix Process](#hotfix-process)
7. [Rollback Procedure](#rollback-procedure)

---

## Release Schedule

### Regular Releases

- **Major releases** (X.0.0): Quarterly (as needed for breaking changes)
- **Minor releases** (x.Y.0): Monthly (new features)
- **Patch releases** (x.y.Z): As needed (bug fixes)

### Release Branches

- **main**: Stable production releases
- **develop**: Active development
- **release/vX.Y.Z**: Release candidate branches

---

## Version Numbering

We follow [Semantic Versioning 2.0.0](https://semver.org/):

### Format: MAJOR.MINOR.PATCH

- **MAJOR**: Incompatible API changes
- **MINOR**: Backward-compatible functionality
- **PATCH**: Backward-compatible bug fixes

### Examples

- `1.0.0`: Initial release
- `1.1.0`: New tool added
- `1.1.1`: Bug fix in existing tool
- `2.0.0`: Breaking change (e.g., tool parameter renamed)

### Pre-release Tags

- `v1.0.0-alpha.1`: Alpha release
- `v1.0.0-beta.1`: Beta release
- `v1.0.0-rc.1`: Release candidate

---

## Pre-Release Checklist

### Code Quality

- [ ] All tests passing (unit + integration)
- [ ] Code coverage ≥ 80%
- [ ] No critical linter warnings
- [ ] Security scan passed
- [ ] Dependencies updated

### Documentation

- [ ] README.md updated
- [ ] CHANGELOG.md updated
- [ ] API documentation current
- [ ] Migration guide (if breaking changes)
- [ ] Quickstart examples tested

### Validation

- [ ] Manual testing on macOS
- [ ] Manual testing on Linux
- [ ] Manual testing on Windows (if applicable)
- [ ] Claude Desktop integration tested
- [ ] All quickstart scenarios validated

### Administrative

- [ ] Version number updated in code
- [ ] Release notes drafted
- [ ] Breaking changes documented
- [ ] Known issues documented

---

## Release Steps

### Step 1: Prepare Release Branch

```bash
# From develop branch
git checkout develop
git pull origin develop

# Create release branch
VERSION="1.2.0"
git checkout -b release/v${VERSION}
```

### Step 2: Update Version Information

```bash
# Update version in main.go
sed -i '' 's/Version = ".*"/Version = "'"${VERSION}"'"/' main.go

# Update README.md if needed
# Update CHANGELOG.md
```

Example `main.go`:

```go
var (
    Version = "1.2.0"
    BuildTime = ""
)
```

### Step 3: Update CHANGELOG.md

```markdown
## [1.2.0] - 2025-XX-XX

### Added
- New feature X
- New feature Y

### Changed
- Improved Z

### Fixed
- Bug #123
- Bug #456

### Deprecated
- Feature W (will be removed in 2.0.0)
```

### Step 4: Run Pre-Release Tests

```bash
# Run all tests
go test ./...

# Run integration tests with goctl
GOCTL_PATH=$(which goctl) go test ./tests/integration/...

# Build binary
go build -o mcp-zero

# Manual smoke test
./mcp-zero --version
```

### Step 5: Create Pull Request

```bash
# Push release branch
git add .
git commit -m "chore: prepare release v${VERSION}"
git push origin release/v${VERSION}

# Create PR from release/v${VERSION} to main
# Request review from maintainers
```

### Step 6: Merge to Main

```bash
# After approval, merge to main
git checkout main
git pull origin main
git merge --no-ff release/v${VERSION}
git push origin main
```

### Step 7: Create Git Tag

```bash
# Tag the release
git tag -a v${VERSION} -m "Release v${VERSION}"
git push origin v${VERSION}
```

### Step 8: Automated Build & Release

GitHub Actions will automatically:

1. Run CI/CD pipeline
2. Build binaries for all platforms
3. Create GitHub release
4. Upload release assets
5. Publish Docker image (if applicable)

### Step 9: Verify Release

```bash
# Check GitHub releases page
open https://github.com/zeromicro/mcp-zero/releases/tag/v${VERSION}

# Verify assets are present:
# - mcp-zero-v1.2.0-linux-amd64.tar.gz
# - mcp-zero-v1.2.0-darwin-amd64.tar.gz
# - mcp-zero-v1.2.0-darwin-arm64.tar.gz
# - mcp-zero-v1.2.0-windows-amd64.zip
# - Checksums for each
```

### Step 10: Merge Back to Develop

```bash
# Merge changes back to develop
git checkout develop
git pull origin develop
git merge --no-ff main
git push origin develop
```

---

## Post-Release Tasks

### Immediate (Same Day)

1. **Announce Release**
   - Post on GitHub Discussions
   - Update project website (if exists)
   - Post on go-zero community channels

2. **Monitor Issues**
   - Watch for bug reports
   - Respond to user questions
   - Track adoption metrics

3. **Update Documentation Sites**
   - Publish new documentation version
   - Update quickstart guides
   - Update example repositories

### Within 1 Week

1. **Gather Feedback**
   - Monitor GitHub issues
   - Check community channels
   - Review usage metrics

2. **Plan Next Release**
   - Create milestone for next version
   - Triage new issues
   - Assign priorities

3. **Update Roadmap**
   - Reflect completed features
   - Add new feature requests
   - Adjust timelines

---

## Hotfix Process

For critical bugs in production release:

### Step 1: Create Hotfix Branch

```bash
# From main branch
git checkout main
git pull origin main

# Create hotfix branch
git checkout -b hotfix/v1.2.1
```

### Step 2: Fix Bug

```bash
# Make minimal changes to fix critical bug
# Add tests to prevent regression
# Update CHANGELOG.md
```

### Step 3: Test Thoroughly

```bash
# Run all tests
go test ./...

# Run specific tests for the fix
go test -v ./tests/integration/affected_test.go

# Manual verification
```

### Step 4: Fast-Track Release

```bash
# Commit and push
git add .
git commit -m "hotfix: critical bug #789"
git push origin hotfix/v1.2.1

# Create PR to main (expedited review)
# After approval, merge immediately
```

### Step 5: Release Hotfix

```bash
# Tag and push
git tag -a v1.2.1 -m "Hotfix v1.2.1"
git push origin v1.2.1

# GitHub Actions will build and release
```

### Step 6: Merge Back

```bash
# Merge hotfix to develop
git checkout develop
git merge --no-ff hotfix/v1.2.1
git push origin develop
```

---

## Rollback Procedure

If a critical issue is discovered after release:

### Option 1: Immediate Hotfix (Preferred)

```bash
# Release hotfix version following hotfix process above
# Example: v1.2.0 has bug → release v1.2.1
```

### Option 2: Revert Tag (Last Resort)

```bash
# Delete problematic tag
git tag -d v1.2.0
git push origin :refs/tags/v1.2.0

# Delete GitHub release
# (Use GitHub UI or API)

# Communicate to users
# Create hotfix or new release
```

### Communication Template

```markdown
## Security Advisory / Critical Bug Notice

**Affected Version**: v1.2.0
**Severity**: Critical
**Status**: Fixed in v1.2.1

### Issue
[Description of the problem]

### Impact
[Who is affected and how]

### Action Required
Update to v1.2.1 immediately:
```bash
go install github.com/zeromicro/mcp-zero@v1.2.1
```

### Timeline
- Issue discovered: [date/time]
- Fix released: [date/time]
- All users notified: [date/time]
```

---

## Release Automation

### GitHub Actions Workflows

1. **`.github/workflows/ci.yml`**
   - Runs on every push/PR
   - Tests, linting, security scans
   - Builds for multiple platforms

2. **`.github/workflows/release.yml`**
   - Triggered by version tags
   - Builds release binaries
   - Creates GitHub release
   - Uploads assets

### Triggering a Release

```bash
# Push a version tag
git tag -a v1.2.0 -m "Release v1.2.0"
git push origin v1.2.0

# GitHub Actions will:
# 1. Build binaries
# 2. Run tests
# 3. Create release
# 4. Upload assets
```

---

## Release Checklist Template

Copy this for each release:

```markdown
## Release v1.X.Y Checklist

### Pre-Release
- [ ] All tests passing
- [ ] Documentation updated
- [ ] CHANGELOG.md updated
- [ ] Version bumped in code
- [ ] Release notes drafted
- [ ] Manual testing completed
- [ ] Security scan passed

### Release
- [ ] Release branch created
- [ ] PR reviewed and approved
- [ ] Merged to main
- [ ] Tag created and pushed
- [ ] GitHub Actions succeeded
- [ ] Assets uploaded
- [ ] Release published

### Post-Release
- [ ] Announcement posted
- [ ] Community notified
- [ ] Documentation site updated
- [ ] Merged back to develop
- [ ] Monitoring for issues
- [ ] Next milestone created
```

---

## Version History

| Version | Date | Type | Highlights |
|---------|------|------|------------|
| 1.0.0 | 2025-XX-XX | Major | Initial release with 10 tools |
| 1.0.1 | TBD | Patch | Bug fixes |
| 1.1.0 | TBD | Minor | New features |
| 2.0.0 | TBD | Major | Breaking changes |

---

## Release Notes Template

```markdown
# mcp-zero v1.X.Y

**Release Date**: YYYY-MM-DD
**Download**: [GitHub Releases](https://github.com/zeromicro/mcp-zero/releases/tag/v1.X.Y)

## 🎉 Highlights

- **New Feature**: Description
- **Improvement**: Description
- **Bug Fix**: Description

## 📦 Installation

```bash
go install github.com/zeromicro/mcp-zero@v1.X.Y
```

## ✨ What's New

### New Features

- Feature 1 (#PR)
- Feature 2 (#PR)

### Improvements

- Improvement 1 (#PR)
- Improvement 2 (#PR)

### Bug Fixes

- Fix 1 (#Issue)
- Fix 2 (#Issue)

## 🔄 Breaking Changes

None in this release.

## 📝 Upgrade Guide

No special steps required. Simply upgrade using:

```bash
go install github.com/zeromicro/mcp-zero@v1.X.Y
```

## 🙏 Contributors

Thank you to all contributors who made this release possible!

- @contributor1
- @contributor2

## 📚 Documentation

- [README](README.md)
- [Quickstart](specs/001-mcp-tool-go-zero/quickstart.md)
- [Contributing](CONTRIBUTING.md)

## 🐛 Known Issues

None at this time.

## ⬆️ What's Next

See our [roadmap](https://github.com/zeromicro/mcp-zero/milestones) for upcoming features.
```

---

## Contact

**Release Manager**: [Name]
**Repository**: https://github.com/zeromicro/mcp-zero
**Issues**: https://github.com/zeromicro/mcp-zero/issues
**Discussions**: https://github.com/zeromicro/mcp-zero/discussions
