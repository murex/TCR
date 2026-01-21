# Building TCR on your machine

This section provides information related to TCR tool development environment setup for those who would like to build
TCR tool locally.

## Clone TCR repository - `Required`

```shell
git clone https://github.com/murex/TCR.git
cd TCR
```

## Install Go SDK - `Required`

TCR is written in Go. This implies having Go compiler and tools installed on your machine.

Simply follow the instructions provided [here](https://go.dev/). Make sure to install **Go version 1.25** or higher.

## Install build tools - `Required`

TCR uses GNU Make to build, test, and run various development tasks.

Please ensure that `make` is installed on your machine before proceeding.

- **Linux**: usually preinstalled or available via your package manager
- **macOS**:
  ```bash
  brew install make
  ```
- **Windows**: use WSL or install GNU Make using tools such as Chocolatey

## CGO configuration - Optional

If you do not have a C compiler (such as gcc) installed on your machine, you may need to disable CGO to avoid build or test errors.

You can do this by setting the following environment variable:

```
export CGO_ENABLED=0
```

On Windows (PowerShell):
```
setx CGO_ENABLED 0
```

This configuration is sufficient for building and testing TCR in environments where gcc is not available.

## Install additional Go tools and utility packages

### Install development dependencies - `Required for testing`

Before running tests, required Go tools must be installed.

From the repository root directory, run:

```
make install-tools
```

### Go IDE - `Optional`

You can check this [link](https://www.tabnine.com/blog/top-7-golang-ides-for-go-developers/)
for a list of recommended IDEs supporting Go language.

### Go test build tags configuration - Required for testing in IDEs

Some TCR tests rely on the Go build tag test_helper.

When running tests from an IDE, make sure this build tag is enabled; otherwise some tests may fail or not be detected correctly.

Examples:

- **VS Code:** Add the following setting in your workspace or user settings:
```
"go.testTags": "test_helper"
```

- **GoLand**: Go to Run / Debug Configurations → Go Test → Go tool arguments and add:
```
-tags=test_helper
```

### Cobra library and tools - `Optional`

TCR Go command line options and parameters are managed using [Cobra](https://github.com/spf13/cobra).

The Cobra library download is already dealt with through Go Module dependencies, which means that in most situations you
will not need to worry about installing it.

In case you need to add or modify subcommands, options or parameters, you may want to use the Cobra Generator. In this
situation you can refer to
[Cobra Generator documentation](https://github.com/spf13/cobra/blob/main/site/content/user_guide.md)

### GoReleaser utility - `Optional`

New versions of TCR Go are released through [GoReleaser](https://goreleaser.com/).

You should not need it as long as you don't plan to release a new TCR Go version.

If you do, you can refer to [GoReleaser Installation Instructions](https://goreleaser.com/install/)
for installing it locally on your machine.

In most cases you will not even have to install it locally as TCR-Go new releases are built through
a [GoReleaser GitHub action](../.github/workflows/go_releaser.yml).

### golangci-lint package - `Optional`

We use the Go Linter aggregator [golangci-lint](https://golangci-lint.run/) to perform various static checks on TCR Go
code.

A [dedicated GitHub action](../.github/workflows/go_releaser.yml) triggers execution of golangci-lint every time a new
TCR-Go version is being released.

Although not mandatory, we advise you to install it locally on your machine to check that your changes comply with
golangci-lint rules. Refer to [golangci-lint installation instructions](https://golangci-lint.run/welcome/install/)
for installation.

Once golangci-lint is installed, you can run it from the root directory:

```shell
make lint
```

Both local run and GitHub Action use [this configuration file](../.golangci.yml)

### gotestsum utility - `Optional`

We use [gotestsum](https://github.com/gotestyourself/gotestsum) for running tests
with the possibility to generate a xunit-compatible test report.

Although not mandatory, we advise you to install it locally on your machine as it greatly improves
readability of test results.
Refer to [gotestsum's Install section](https://github.com/gotestyourself/gotestsum#install)
for installation.

Once gotestsum is installed and development dependencies are installed, you can run make's test target from the root directory:

- For running all tests:

  ```shell
  make test
  ```
- For running short tests only:

  ```shell
  make test-short
  ```

- For listing slowest tests (default threshold: 500ms):

  ```shell
  make slow-tests
  ```

## Build TCR executable

To build TCR locally on your machine, simply type the following from the root directory:

```shell
make
```

This command generates the TCR (in [src](../src) directory) executable, as well as the command help pages (
in [doc](../doc) directory).

### To build TCR only

Either run the following command from the root directory:

```shell
make -C ./src
```

Or run make from [src](../src) directory:

```shell
cd src
make
```

### To generate TCR command markdown documentation

```shell
make doc
```

## Testing local TCR manually

The [tools/tcr/tcr-local.sh](../tools/tcr/tcr-local.sh) script runs the latest TCR built from local sources on
the [src/testdata/java](../src/testdata/java) sample.
The [src/tcr-local](../src/tcr-local) script is just a convenience wrapper for this.

To launch it:

```shell
cd src
./tcr-local
```

You can alternatively use the following `make` target from either the repository root directory
or [src](../src) directory:

```shell
make run
```

If you want to test with the Perforce client:

```shell
cd src
./tcr-local-p4
```
