# AGENTS.md

This file provides guidance to AI coding agents (Claude Code, GitHub Copilot, Cursor, etc.) when working with code in this repository.

## Project Overview

TCR (Test && Commit || Revert) is a programming workflow tool written in Go that enforces baby-step development practices. The concept was originally developed by Kent Beck, Oddmund Strømme, Lars Barlindhaug, and Ole Johannessen. This tool automatically commits code changes when tests pass, or reverts them when tests fail, encouraging developers to work in small, incremental steps.

The project includes both a CLI tool written in Go and an experimental Angular web interface.

## Architecture

### Multi-Module Structure

The project is organized as a multi-module monorepo with dependencies:
- **src/** - Main Go application (depends on webapp)
  - Entry point: `main.go`
  - CLI commands: `cmd/` directory using Cobra framework
  - Core engine: `engine/` directory
  - Language/toolchain support: `language/` and `toolchain/` directories
  - Version control integration: `vcs/` directory
  - Web interface: `http/` directory
  - Configuration: `config/` and `settings/` directories
  - Events: `events/` directory for event-driven architecture
  - Filesystem: `filesystem/` directory with fsnotify for watching
  - Reporting: `report/`, `xunit/`, `stats/` directories

- **webapp/** - Angular-based web interface (experimental)
  - Angular 21+ TypeScript frontend
  - Real-time communication via WebSocket
  - Displays TCR cycle status, timer, and role management
  - Built output embedded in Go binary via static assets

- **tcr-doc/** - Documentation generator (depends on src)

- **examples/** - Language/toolchain examples
  - Supports multiple programming languages
  - Multiple build tools per language (Maven, Gradle, Cargo, etc.)
  - Each example includes README and sample code
  - Serve as integration tests

- **doc/** - Generated command documentation (Markdown)
- **dev-doc/** - Development documentation
- **tools/** - Utility scripts and tools
- **auto-install/** - One-liner installation scripts

## Build System

### Root Makefile Targets

```bash
make help           # Show available targets
make prepare        # Full pipeline: deps, tidy, lint, build, doc, test
make build          # Build all production modules (webapp + src)
make test           # Run tests for all production modules
make lint           # Run linter on all modules
make doc            # Generate command line documentation
make release        # Create release using GoReleaser
make snapshot       # Create snapshot release
```

### Go Module (src/) Targets

```bash
cd src
make build          # Build TCR binary with version info
make test           # Run all tests with coverage
make test-short     # Run short tests only
make lint           # Run golangci-lint
make run            # Test locally built TCR on testdata
make tidy           # Run go mod tidy
make cov            # Show coverage report in browser
make deps           # Download dependencies
make vet            # Run go vet
```

### Angular Module (webapp/) Targets

```bash
cd webapp
make setup          # Install npm dependencies
make build          # Build for production (output to Go's static assets)
make run            # Start development server
make test           # Run Vitest unit tests
make lint           # Run ESLint
make cov            # Run tests with coverage
make clean          # Clean build artifacts
make deps           # Update dependencies
```

### Running TCR Locally

Use `./src/tcr-local` script to test locally built binary on example projects.

## Key Technologies & Dependencies

### Backend (Go)
- **Go 1.26+** (required)
- **Cobra** - CLI framework
- **Viper** - Configuration management
- **Gin** - Web framework for HTTP API
- **go-git** - Git operations
- **fsnotify** - File system watching
- **WebSocket** - Real-time communication with frontend
- **testify** - Testing framework

### Frontend (Angular)
- **Node.js 22** (for development)
- **Angular 21+**
- **TypeScript**
- **Vitest** - Testing

### Build Tools
- **GoReleaser** - Release automation
- **golangci-lint** - Go linting
- **gotestsum** - Enhanced test output

## Testing

### Go Tests
- Use testify framework for assertions
- Tests tagged with `test_helper` build tag
- Coverage reports in `src/_test_results/`
- Run specific tests: `go test -tags=test_helper ./path/to/package`
- Integration tests with real VCS operations
- Test data in `src/testdata/`

### Angular Tests
- Vitest for unit tests
- Run with `npm test` or `make test` in webapp/
- ESLint for code quality
- Coverage reporting

### Test Data
- Examples in `examples/` directory serve as integration tests
- Go test fixtures in `src/testdata/`

## Configuration System

TCR uses hierarchical YAML configuration:
1. Built-in defaults (embedded in binary)
2. Global config: `~/.tcr/config.yml`
3. Repository config: `<repo>/.tcr/config.yml`
4. Command-line flags (highest priority)

### Configuration Files
- `.tcr/config.yml` - Main TCR settings
- `.tcr/language/*.yml` - Language definitions (source/test file patterns)
- `.tcr/toolchain/*.yml` - Toolchain definitions (build/test commands)

Key configuration sections:
- Language and toolchain detection
- VCS integration (Git, Perforce)
- Commit message templates
- Test timeout settings
- File filtering rules

## Supported Languages & Toolchains

TCR supports multiple programming languages with multiple build tools:

### Built-in Languages
- **C++**: cmake, bazel
- **C#**: dotnet
- **Elixir**: mix
- **Go**: go-tools, gotestsum, make, bazel
- **Haskell**: stack
- **Java**: gradle, gradle-wrapper, maven, maven-wrapper, make, bazel
- **JavaScript**: yarn
- **Kotlin**: gradle, gradle-wrapper, maven, maven-wrapper
- **PHP**: phpunit
- **Python**: pytest, bazel
- **Rust**: cargo, nextest
- **Scala**: sbt
- **TypeScript**: yarn

### Adding New Language
1. Create `.tcr/language/<name>.yml` with:
   - `toolchains.default` - Default toolchain
   - `toolchains.compatible-with` - List of compatible toolchains
   - `source-files.directories` - Where source files are
   - `source-files.patterns` - Regex patterns for source files (RE2 syntax)
   - `test-files.directories` - Where test files are
   - `test-files.patterns` - Regex patterns for test files

2. Add language implementation in `src/language/` if custom logic needed
3. **EXECUTE**: `make lint && make test`

### Adding New Toolchain
1. Create `.tcr/toolchain/<name>.yml` with build and test commands for each OS/arch
2. Add toolchain implementation in `src/toolchain/` if needed
3. **EXECUTE**: `make lint && make test`

## CLI Commands

### Core Commands
- `tcr solo` - Run TCR in solo mode
- `tcr mob` - Run TCR in mob programming mode
- `tcr one-shot` - Run TCR once and exit
- `tcr web` - Start web interface

### Utility Commands
- `tcr check` - Check TCR configuration and environment
- `tcr config` - Manage TCR configuration
- `tcr info` - Display TCR information
- `tcr log` - Show TCR activity log
- `tcr stats` - Display TCR statistics
- `tcr retro` - Generate retrospective report

## Version Control Integration

### Git (Default)
- Automatic commits on test pass
- Automatic reverts on test failure
- Full support with automatic commit/revert
- Note: TCR commits are deliberately unsigned (would be impractical with frequent auto-commits)

### Perforce
- Use `--vcs=p4` flag
- Requires P4 client configuration
- Limited support: no auto-push, log, or stats support
- Set `P4IGNORE=.gitignore` to avoid committing build artifacts

## Key Development Workflows

### Adding a New CLI Command
1. Create `src/cmd/<command>.go`
2. Define Cobra command with flags
3. Wire up to root command in `src/cmd/root.go`
4. Add implementation in appropriate package
5. Regenerate docs: `cd tcr-doc && make doc`
6. **EXECUTE**: `make lint && make test`

### Modifying TCR Engine
1. Core logic in `src/engine/`
2. Event-driven architecture via `src/events/`
3. Update tests in corresponding `_test.go` files
4. **EXECUTE**: `make lint && make test`

### Working on Web Interface
1. Start Go backend: `cd src && ./tcr-local web -T=http`
2. Start Angular dev: `cd webapp && npm start`
3. Make changes in `webapp/src/`
4. **EXECUTE**: `cd webapp && make lint && make test`

### Bug Fixes
1. Reproduce issue with test case
2. Fix in appropriate module
3. Ensure tests pass on all platforms
4. Update documentation if needed
5. **EXECUTE**: `make lint && make test`

### Feature Development
1. Design with TCR principles in mind
2. Add CLI command in `src/cmd/` if needed
3. Implement core logic in appropriate module
4. Add comprehensive tests
5. Update web interface if applicable
6. **EXECUTE**: `make lint && make test`

## Quality Gates (MANDATORY)

**CRITICAL**: After ANY code change (modification, addition, deletion, refactoring, formatting, import organization, etc.), you MUST ALWAYS execute these quality gates to prevent regressions:

### Required Steps (Execute in Order):

1. **Run Linter**: `make lint`
   - Must pass with 0 issues
   - Applies to both Go and TypeScript/Angular code
   - Non-negotiable - fix all linting errors before proceeding

2. **Run Tests**: `make test` or `make test-short` for Go changes
   - All tests must pass
   - No skipped tests due to failures
   - Verify test coverage is maintained

3. **Verify Build**: `make build` (if significant changes)
   - Ensure project still compiles and builds successfully
   - Check both webapp and Go modules

### Quality Gate Commands by Scope

```bash
# After any Go code changes
cd src && make lint && make test-short  # or make test for full suite

# After webapp changes  
cd webapp && make lint && make test

# After project-wide changes
make lint && make test
```

### Examples of Changes Requiring Quality Gates:
- Import organization/formatting (like `goimports`)
- Code refactoring or restructuring
- Adding/removing dependencies
- Modifying configuration files
- Updating documentation that affects code
- Any file modifications in `src/`, `webapp/`, or `examples/`

### Failure Handling:
- If linter fails: Fix all issues before proceeding
- If tests fail: Investigate and fix, don't ignore
- If build fails: Resolve compilation errors immediately
- Document any issues encountered and their resolution

**Remember**: The TCR philosophy applies to development too - small, verified steps prevent large regressions.

## AI Agent Response Pattern

**Every AI agent MUST follow this pattern:**

1. Make the requested changes
2. Execute quality gates: `make lint && make test`
3. Report results with ✅ or ❌ status
4. Fix any issues found
5. Re-run quality gates until all pass

**Never skip quality gates**, even for "simple" changes like formatting or imports.

## Code Organization Patterns

### Event-Driven Communication
- Events defined in `src/events/`
- Publishers emit events via event bus
- Subscribers listen for specific events
- Used for engine → UI communication

### Role Management
- `src/role/` handles driver/navigator roles
- Synchronized across participants in mob mode
- Timer in `src/timer/` for driver rotation

### Filesystem Watching
- `src/filesystem/` uses fsnotify
- Monitors source and test directories
- Triggers TCR cycle on changes

### Report Generation
- `src/report/` handles test output
- `src/xunit/` parses xUnit XML format
- `src/stats/` tracks commit/revert statistics

## Common Gotchas

1. **Build Order**: webapp must build before src (embedded assets)
2. **Test Tags**: Go tests require `-tags=test_helper` flag
3. **Coverage Files**: Multiple coverage files for SonarCloud vs local use
4. **Version Info**: Set via ldflags in build, defined in `src/settings/`
5. **Cross-Platform**: Commands defined per OS/arch in toolchain configs
6. **Web Port**: Default 8483, configurable via `--port-number`

## Development Tools Required

- Go 1.26+
- Node.js 22 (for webapp)
- golangci-lint (for linting)
- gotestsum (optional, better test output)
- GoReleaser (for releases)

## File Patterns to Understand

- `**/*_test.go` - Go test files
- `**/*.spec.ts` - Angular test files
- `**/testdata/**` - Test fixtures and data directories
- `**/.tcr/**` - TCR configuration
- `**/Makefile` - Build configuration
- `**/README.md` - Documentation files
- `src/_test_results/` - Go test outputs
- `webapp/dist/` - Angular build output

## Working with Examples

Examples in `examples/` demonstrate TCR usage with different language/toolchain combinations. Each example:
- Has a complete project setup
- Includes README with TCR usage
- Serves as integration test
- Can be run with locally built TCR using `./src/tcr-local`

## CI/CD

- GitHub Actions workflows in `.github/workflows/`
- Go tests, linting, and builds
- Angular builds and tests
- SonarCloud quality scanning
- Coveralls coverage reporting
- GoReleaser for multi-platform releases
- Dependabot for dependency updates

## Release Process

- **Automated**: GitHub Actions with GoReleaser
- **Manual**: `make release` or `make snapshot`
- **Platforms**: Windows, macOS, Linux (multiple architectures)
- **Distribution**: GitHub releases + package managers

## Development Guidelines

1. **Testing**: All changes should include tests
2. **Documentation**: Update relevant docs in `doc/` and `dev-doc/`
3. **Linting**: Code must pass all linters
4. **Baby Steps**: Follow TCR principles during development
5. **Cross-Platform**: Ensure compatibility across OS platforms
6. **Quality Gates**: Always run `make lint && make test` after changes

This project emphasizes clean architecture, comprehensive testing, and cross-platform compatibility while maintaining the core philosophy of encouraging small, incremental development steps.
