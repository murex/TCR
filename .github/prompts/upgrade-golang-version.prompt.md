---
description: 'Upgrade Go version from 1.n to 1.n+1 across the entire project'
argument-hint: 'from=1.n to=1.n+1'
agent: 'agent'
tools: ['terminal', 'search/codebase', 'edit']
---

# Upgrade Go Version

This prompt guides you through upgrading the Go version across the entire TCR project, including all modules, examples, configuration files, workflows, and documentation.

## Input Parameters

- `${input:from}` - Current Go version (e.g., 1.25)
- `${input:to}` - Target Go version (e.g., 1.26)

If not provided in the chat input, ask the user for these values.

## Upgrade Process

Follow these steps in order to ensure a complete and verified upgrade:

### 1. Update Go Module Files

Search for and update all `go.mod` files in the project:

- Main application: `src/go.mod`
- Documentation generator: `tcr-doc/go.mod`
- All examples: `examples/*/go.mod` (go-bazel, go-go-tools, go-gotestsum, go-make, etc.)

Change the `go` directive from the old version to the new version.

### 2. Update Go Workspace Files

Search for and update all `go.work` files:

- Examples workspace: `examples/go.work`
- Any other workspace files in the project

Update the `go` directive to the new version.

### 3. Update Bazel Configuration

If the project uses Bazel, update:

- `examples/go-bazel/MODULE.bazel` - Change `go_sdk.download(version = "1.n.0")` to the new version

### 4. Update GitHub Workflow Files

Update all GitHub Actions workflow files in `.github/workflows/`:

- `go.yml` - Update `go-version: '1.n'` to the new version
- `go_releaser.yml` - Update `go-version: '1.n'` to the new version
- `golangci_lint.yml` - Update matrix `go: ["1.n"]` to the new version

### 5. Update Linter Configuration

Update `.golangci.yml`:

- Change `run.go: "1.n"` to the new version

### 6. Update Documentation

Search for and update all documentation files mentioning the Go version:

- `AGENTS.md` - Update "Go 1.n+" references to the new version
- `dev-doc/build.md` - Update "Go version 1.n" references
- `README.md` - Check for any version-specific mentions (usually not needed as users download binaries)
- `CONTRIBUTING.md` - Check for version requirements

### 7. Update golangci-lint Tool

After updating go.mod files, rebuild golangci-lint with the new Go version:

```bash
go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest
```

Verify the new version:

```bash
golangci-lint version
```

Expected output should show: `built with go1.n+1`

### 8. Verify Build Pipeline

Run the complete build pipeline to ensure everything works:

```bash
make prepare
```

This will:
- Update dependencies for all modules
- Run linting on all modules
- Build all modules (webapp + src)
- Generate documentation
- Run all tests (webapp + src)

Expected result: All stages should pass with 0 errors.

### 9. Verify Go Examples

For each Go example in `examples/go-*`, verify compatibility:

#### For each example directory:

1. **Run TCR check**:
   ```bash
   cd examples/go-<toolchain>
   ./tcrw check
   ```

2. **Extract and run build command** (from check output):
   - Look for "build command line:" in the output
   - Run the specified command (e.g., `bazel build ...`, `go test -count=0 ./...`, `make build`)

3. **Extract and run test command** (from check output):
   - Look for "test command line:" in the output
   - Run the specified command (e.g., `bazel test ...`, `gotestsum ...`, `go test -short ./...`, `make test`)

4. **Verify success**:
   - Build should complete without errors
   - Tests should pass

#### Examples to verify:
- `go-bazel` - Uses Bazel build system
- `go-gotestsum` - Uses gotestsum for testing
- `go-go-tools` - Uses standard Go tools
- `go-make` - Uses Makefile

### 10. Create Summary

After completing all steps, provide a summary including:

- Total number of files updated
- List of all updated files grouped by category (go.mod, workflows, configs, docs, etc.)
- Build verification results
- Test results (number of tests passed)
- Example verification results (for each example: build status, test status)

## Files to Search For

Use search patterns to find all files that need updating:

- `**/*.mod` - Go module files
- `**/*.work` - Go workspace files
- `**/*.yml` - Workflow and config files
- `**/*.yaml` - Workflow and config files
- `**/*.md` - Documentation files
- `**/MODULE.bazel` - Bazel configuration

## Common Issues and Solutions

### Issue: golangci-lint panics with "file requires newer Go version"

**Solution**: Rebuild golangci-lint with the new Go version using:
```bash
go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest
```

### Issue: go.work file version mismatch

**Solution**: Don't forget to update workspace files (`go.work`) in addition to module files (`go.mod`).

### Issue: Bazel fails to download Go SDK

**Solution**: Ensure `MODULE.bazel` has the correct Go SDK version specified in `go_sdk.download(version = "...")`.

### Issue: Tests fail after upgrade

**Solution**: Check for deprecated APIs or breaking changes in the Go release notes. Most minor version upgrades (1.n to 1.n+1) are backward compatible.

## Quality Gates (MANDATORY)

Before considering the upgrade complete, verify:

- [ ] All `go.mod` files updated
- [ ] All `go.work` files updated (if present)
- [ ] All GitHub workflow files updated
- [ ] Linter configuration updated
- [ ] Documentation updated
- [ ] golangci-lint rebuilt with new Go version
- [ ] `make prepare` passes with 0 errors
- [ ] All project tests pass
- [ ] All Go examples build successfully
- [ ] All Go examples tests pass

## Expected Outcome

A complete, verified upgrade of the Go version across the entire project with:
- Zero build errors
- Zero test failures
- All linting passing
- All examples working
- Documentation accurately reflecting the new version requirement