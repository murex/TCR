# TCR - Test && Commit || Revert - Copilot Instructions

## Project Overview

TCR (Test && Commit || Revert) is a programming workflow tool written in Go that enforces baby-step development practices. The concept was originally developed by Kent Beck, Oddmund Strømme, Lars Barlindhaug, and Ole Johannessen. This tool automatically commits code changes when tests pass, or reverts them when tests fail, encouraging developers to work in small, incremental steps.

## Project Structure

The project is organized into several key modules:

### Core Components

- **`src/`** - Main Go application source code
  - Entry point: `main.go`
  - CLI commands: `cmd/` directory using Cobra framework
  - Core engine: `engine/` directory
  - Language/toolchain support: `language/` and `toolchain/` directories
  - Version control integration: `vcs/` directory
  - Web interface: `http/` directory
  - Configuration: `config/` and `settings/` directories

- **`webapp/`** - Angular-based web interface (experimental)
  - TypeScript frontend
  - Communicates with Go backend via HTTP/WebSocket

- **`examples/`** - Language/toolchain examples
  - Supports multiple programming languages
  - Multiple build tools per language (Maven, Gradle, Cargo, etc.)
  - Each example includes README and sample code

- **`doc/`** - Generated command documentation (Markdown)
- **`dev-doc/`** - Development documentation
- **`tools/`** - Utility scripts and tools
- **`auto-install/`** - One-liner installation scripts

## Build System

### Root Makefile Targets

- `make help` - Show available targets
- `make prepare` - Full preparation (deps, tidy, lint, build, doc, test)
- `make build` - Build all production modules (webapp + src)
- `make test` - Run tests for all production modules
- `make lint` - Run linter on all modules
- `make doc` - Generate command line documentation
- `make release` - Create release using GoReleaser
- `make snapshot` - Create snapshot release

### Go Module (src/) Targets

- `make build` - Build TCR binary with version info
- `make test` - Run all tests with coverage
- `make test-short` - Run short tests only
- `make lint` - Run golangci-lint
- `make run` - Test locally built TCR on Java testdata
- `make tidy` - Run go mod tidy
- `make cov` - Show coverage report in browser

### Angular Module (webapp/) Targets

- `make setup` - Install npm dependencies
- `make build` - Build for production (output to Go's static assets)
- `make run` - Start development server
- `make test` - Run unit tests
- `make lint` - Run ESLint
- `make clean` - Clean build artifacts

## Key Technologies & Dependencies

### Backend (Go)
- **Go 1.25+** (required)
- **Cobra** - CLI framework
- **Viper** - Configuration management
- **Gin** - Web framework for HTTP API
- **go-git** - Git operations
- **fsnotify** - File system watching
- **WebSocket** - Real-time communication with frontend

### Frontend (Angular)
- **Node.js 22** (for development)
- **Angular 20+**
- **TypeScript**
- **Karma/Jasmine** - Testing

### Build Tools
- **GoReleaser** - Release automation
- **golangci-lint** - Go linting
- **gotestsum** - Enhanced test output

## Development Workflow

### Local Development Setup

1. **Clone repository**: `git clone https://github.com/murex/TCR.git`
2. **Install Go 1.25+**
3. **Install Node.js 22** (for webapp development)
4. **Optional tools**: golangci-lint, gotestsum, GoReleaser

### Building & Testing

```bash
# Build everything
make build

# Run tests
make test

# Run linter
make lint

# Test locally built TCR
cd src && ./tcr-local
```

### Web Development

```bash
# Start Angular dev server
cd webapp && npm start

# Start Go backend in development mode
cd src && ./tcr-local web -T=http
```

## Supported Languages & Toolchains

TCR supports multiple programming languages with multiple build tools:

- **C++**: cmake, bazel
- **C#**: dotnet
- **Elixir**: mix
- **Go**: go tools, gotestsum, make, bazel
- **Haskell**: stack
- **Java**: gradle, maven, make, bazel (with wrapper support)
- **JavaScript**: yarn
- **Kotlin**: gradle, maven (with wrapper support)
- **PHP**: phpunit
- **Python**: pytest, bazel
- **Rust**: cargo, nextest
- **Scala**: sbt
- **TypeScript**: yarn

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

## Configuration

TCR uses YAML configuration files with the following hierarchy:
1. Global config: `~/.tcr/config.yml`
2. Repository config: `<repo>/.tcr/config.yml`
3. Command-line flags

Key configuration sections:
- Language and toolchain detection
- VCS integration (Git, Perforce)
- Commit message templates
- Test timeout settings
- File filtering rules

## Version Control Support

- **Git**: Full support with automatic commit/revert
- **Perforce**: Limited support with specific workflows

## Testing Strategy

### Go Tests
- Unit tests with testify framework
- Integration tests with real VCS operations
- Coverage reporting to Coveralls and SonarCloud
- Test data in `src/testdata/`

### Angular Tests
- Unit tests with Karma/Jasmine
- ESLint for code quality
- Coverage reporting

## Release Process

- **Automated**: GitHub Actions with GoReleaser
- **Manual**: `make release` or `make snapshot`
- **Platforms**: Windows, macOS, Linux (multiple architectures)
- **Distribution**: GitHub releases + package managers

## Code Quality

- **SonarCloud**: Code quality metrics
- **Coveralls**: Test coverage tracking
- **golangci-lint**: Go code linting
- **ESLint**: TypeScript/JavaScript linting
- **Dependabot**: Dependency updates

## Development Guidelines

1. **Testing**: All changes should include tests
2. **Documentation**: Update relevant docs in `doc/` and `dev-doc/`
3. **Linting**: Code must pass all linters
4. **Baby Steps**: Follow TCR principles during development
5. **Cross-Platform**: Ensure compatibility across OS platforms

## Common Tasks for AI Agents

### Adding New Language Support
1. Create example in `examples/<language>-<toolchain>/`
2. Add language detection in `src/language/`
3. Add toolchain implementation in `src/toolchain/`
4. Update documentation and tests

### Bug Fixes
1. Reproduce issue with test case
2. Fix in appropriate module
3. Ensure tests pass on all platforms
4. Update documentation if needed

### Feature Development
1. Design with TCR principles in mind
2. Add CLI command in `src/cmd/` if needed
3. Implement core logic in appropriate module
4. Add comprehensive tests
5. Update web interface if applicable

## File Patterns to Understand

- `**/*_test.go` - Go test files
- `**/testdata/**` - Test data directories
- `**/*.spec.ts` - Angular test files
- `**/Makefile` - Build configuration
- `**/.tcr/config.yml` - TCR configuration files
- `**/README.md` - Documentation files

This project emphasizes clean architecture, comprehensive testing, and cross-platform compatibility while maintaining the core philosophy of encouraging small, incremental development steps.
