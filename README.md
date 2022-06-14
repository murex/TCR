[![Go](https://github.com/murex/tcr/actions/workflows/go.yml/badge.svg)](https://github.com/murex/tcr/actions/workflows/go.yml)
[![Go lint](https://github.com/murex/tcr/actions/workflows/golangci_lint.yml/badge.svg)](https://github.com/murex/tcr/actions/workflows/golangci_lint.yml)
[![sonarcloud](https://sonarcloud.io/api/project_badges/measure?project=murex_TCR&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=murex_TCR)
[![Coveralls](https://coveralls.io/repos/github/murex/TCR/badge.svg?branch=main)](https://coveralls.io/github/murex/TCR?branch=main)
[![goreleaser](https://github.com/murex/tcr/actions/workflows/goreleaser.yml/badge.svg)](https://github.com/murex/tcr/actions/workflows/goreleaser.yml)
[![Add contributors](https://github.com/murex/tcr/actions/workflows/add_contributors.yml/badge.svg)](https://github.com/murex/tcr/actions/workflows/add_contributors.yml)

# TCR - Test && (Commit || Revert) application

_A Go implementation of TCR, for practicing baby-steps development, and much more!_

## What is this?

TCR is a programming workflow, standing for **Test && (Commit || Revert)**.

Kent Beck and Oddmund Strømme came up with this concept
in [this post](https://medium.com/@kentbeck_7670/test-commit-revert-870bbd756864).

Since then several people experimented with this idea.

This repository provides a tool allowing people to use the TCR workflow.

## Why TCR?

Although probably a bit challenging to use on real production code, we found TCR to be quite beneficial when used as a
learning and practicing tool on the katas that we use when doing software craftsmanship coaching.

TCR enforces developing in baby steps, with a strong focus on always keeping the green light on tests. Having a TCR tool
feels a bit like having a coaching assistant constantly enforcing such practices!

We initially came up with a small shell script implementing this workflow, and decided to embed it in each of our katas
so that people can use it if they like. This was a way for us to quickly provide a usable TCR solution. However, shell
scripts are not the best in class when it comes to maintainability and changeability.

For This reason we decided to work on a new implementation of TCR, written in [Go](https://golang.org/) this time.

## Who is this for?

If you are a developer willing to practice TCR workflow and baby-step development either alone, in pair or in a mob,
just download and run this tool with the piece of code on which you want to practice. You can also use it on real
production code if you feel like it!

If you are a technical coach, you can advise participants to your coaching sessions to use it during the sessions.

## Where to start?

### Prerequisites

1. Have [git](https://git-scm.com/) installed on your machine
2. Have a clone of the git repository containing the code you intend to work on
3. Supported platforms: macOS, Linux and Windows. Refer to [TCR releases page](https://github.com/murex/TCR/releases)
   for a complete list of supported platforms/architectures

### Languages and toolchains

TCR can potentially work with any programming language. The only things it needs to know for a language are the
following:

- Where to find, and how to recognize source files
- Where to find, and how to recognize test files
- How to build the code
- How to run the tests

The 2 first points are defined through a `language` setting

The 2 other points are defined through a `toolchain` setting

TCR comes with a few built-in languages and toolchains.

You can customize these built-in settings if needed. You can also add your own languages and toolchains if they are not
provided as built-in.

#### Built-in languages and toolchains

| Language | Toolchains                                         | Default        |
|----------|----------------------------------------------------|----------------|
| java     | gradle, gradle-wrapper, maven, maven-wrapper, make | gradle-wrapper |
| cpp      | cmake, make                                        | cmake          |
| go       | go-tools, gotestsum, make                          | go-tools       |
| csharp   | dotnet, make                                       | dotnet         |

### Base directory

In order to know which files and directories to watch, TCR needs to know on which part of the filesystem it should work.
This is what we call the `base directory`.

- The base directory can be specified when starting TCR using the `-b` (or `--base-dir`) command line option.
- When the base directory is not provided, TCR assumes that the current directory is the base directory.

### Work directory

TCR `work directory` is the directory from which build and test commands are launched. In most cases using the same
value as the `base directory` is sufficient.

In some situations (for instance on multi-component projects), it might be necessary to run build and test tools from a
different directory than the one where source and test files are located.

- The work directory can be specified when starting TCR using the `-w` (or `--work-dir`) command line option.
- When the work directory is not provided, TCR assumes that the current directory is the work directory.

### Configuration directory

If you want to save non-default TCR configuration options, customize a built-in language or toolchain, or add your own
language/toolchain, TCR needs to know where it should save them. This is the purpose of the `configuration directory`.

- The configuration directory can be specified when starting TCR using the `-c` (or `--config-dir`) command line option.
- When the configuration directory is not provided, TCR assumes that the current directory is the configuration
  directory.

<details>
  <summary>Configuration directory layout</summary>

- `configuration directory`/
    - `.tcr`/
        - `config.yml` - contains all TCR configuration settings
        - `language`/ - subdirectory containing all language configurations
            - `java.yml` - configuration for Java language
            - `cpp.yml` - configuration for C++ language
            - etc.
        - `toolchain`/ - subdirectory containing all toolchain configurations
            - `gradle.yml` - configuration for Gradle toolchain
            - `gradle-wrapper.yml` - configuration for Gradle wrapper toolchain
            - `cmake.yml` - configuration for CMake toolchain
            - etc.

</details>

### Running TCR

<details>
  <summary>On MacOS</summary>

1. Download the latest version of TCR for Darwin from [here](https://github.com/murex/TCR/releases)

2. Extract TCR executable (replace with the appropriate version and architecture)

    ```shell
    tar zxf tcr_0.12.0_Darwin_x86_64.tar.gz
    ```

3. Launch TCR

    ```shell
    ./tcr -b <base-directory> -w <work-directory> -l <language> -t <toolchain>
    ```

</details>

<details>
  <summary>On Linux</summary>

1. Download the latest version of TCR for Linux from [here](https://github.com/murex/TCR/releases)

2. Extract TCR executable (replace with the appropriate version and architecture)

    ```shell
    tar zxf tcr_0.12.0_Linux_x86_64.tar.gz
    ```

3. Launch TCR

    ```shell
    ./tcr -b <base-directory> -w <work-directory> -l <language> -t <toolchain>
    ```

</details>

<details>
  <summary>On Windows</summary>

1. Download the latest version of TCR for Windows from [here](https://github.com/murex/TCR/releases)

2. Extract TCR executable (replace with the appropriate version and architecture)

    ```shell
    tar zxf tcr_0.12.0_Windows_x86_64.tar.gz
    ```

3. Launch TCR

    ```shell
    ./tcr.exe -b <base-directory> -w <work-directory> -l <language> -t <toolchain>
    ```

</details>

> ***Note***
>
> <details><summary>TCR and git commits signing</summary>
>
> Some users prefer to set up their git configuration so that each of their commits is
> signed and verified through a GPG passphrase as described
> [here](https://git-scm.com/book/en/v2/Git-Tools-Signing-Your-Work).
>
> TCR automatically performs a significant number of commits.
> It would become unusable if the user had to enter a passphrase at each commit.
>
> For this reason, ***TCR commits are deliberately not signed***.
>
> If signing every commit is important to you, you can still do it when you're done
> working with TCR, when reworking git history and squashing TCR commits into meaningful ones.
>
> </details>

### Using TCR configuration

TCR runs by default without any local configuration, using either built-in settings or settings defined through command
line options.

If you prefer to use your own custom settings, you have the possibility to save them locally without the need to provide
them every time you run TCR. This can be achieved through the `tcr config` subcommand.

All configuration settings are saved in YAML format. Thus you can edit them later using a text editor.

<details><summary>Expand for usage examples</summary>

- To save TCR configuration in your HOME directory (using the default settings):

    ```shell
    ./tcr config save -c $HOME
    ```
- To save TCR configuration in the current directory, setting the timer duration to 10m, the language to java and the
  toolchain to maven:

    ```shell
    ./tcr config save -d 10m -l java -t maven
    ```

- To show the current TCR configuration settings (previously saved in the current directory)

    ```shell
    ./tcr config show
    ```

- To reset TCR configuration settings to their default values (in the current directory)

    ```shell
    ./tcr config reset
    ```

</details>

### Adding a new language and toolchain

New languages and toolchains can be added through adding related configuration files in the configuration directory.

<details><summary>Expand for details</summary>

We're assuming in this example that the configuration directory is your HOME directory. Replace `$HOME` in the examples
below if you prefer to use a different location.

Suppose you want to run TCR with Javascript language and yarn toolchain. Here is what you would do:

1. Create the TCR configuration directory structure (you can skip this step if you saved TCR configuration before)

    ```shell
    tcr config save -c $HOME
    ```

2. Create `yarn.yml` toolchain configuration file from an existing toolchain configuration

    ```shell
    cd $HOME/.tcr/toolchain
    cp gradle.yml yarn.yml
    ```

3. Adjust `yarn.yml` contents

   Edit `yarn.yml` with your favorite editor and adjust the contents as follows.

   We're assuming here that yarn is installed and that `yarn build` and `yarn test`
   are set up so that they run respectively the build and test.

    ```yaml
    build:
    - os: [darwin, linux, windows]
      arch: ["386", amd64, arm64]
      command: yarn
      arguments: [build]
    test:
    - os: [darwin, linux, windows]
      arch: ["386", amd64, arm64]
      command: yarn
      arguments: [test]
    ```

4. Create `javascript.yml` language configuration file from an existing language configuration

    ```shell
    cd $HOME/.tcr/language
    cp java.yml javascript.yml
    ```

5. Adjust `javascript.yml` contents

   Edit `javascript.yml` with your favorite editor and adjust the contents as follows:

   We're assuming here that all source files are under `src` subdirectory and are named `*.js`, and that all test files
   are under `test` subdirectory and are named
   `*.test.js`.

    ```yaml
    toolchains:
      default: yarn
      compatible-with: [yarn]
    source-files:
      directories: [src]
      patterns: ['(?i)^.*\.js$']
    test-files:
      directories: [test]
      patterns: ['(?i)^.*\.test\.js$']
    ```

   > ***Regex on filenames***
   >
   > TCR complies with [RE2](https://github.com/google/re2/wiki/Syntax)
   > for pattern matching on filenames.

6. Check TCR settings with the newly configured language and toolchain

   TCR's `check` subcommand performs a number of checks on configuration, parameters and local environment without
   triggering the TCR cycle. It helps you quickly tune your configuration and command line parameters. Make sure that
   there is no error checkpoint in the trace displayed before proceeding any further.

   ```shell
   cd <base-directory>
   tcr check -c $HOME -l javascript -t yarn
   ```

7. Try running TCR with the newly configured language and toolchain

   TCR's `one-shot` subcommand runs one single TCR cycle then exits. Through checking
   its [return code](doc/tcr_one-shot.md) you can quickly verify that everything works as expected.

   ```shell
   cd <base-directory>
   tcr one-shot -c $HOME -l javascript -t yarn
   echo $?
   ```

   > ***Language's Default Toolchain***
   >
   > Each language has a default toolchain, which is the one
   > that will be used with this language if no toolchain is specified on
   > the command line.
   > In the above example, `-t yarn` could actually be skipped.

</details>

### Command line help (all platforms)

Refer to [here](./doc/tcr.md) for TCR command line help and additional options.

## Building TCR on your machine

This section provides information related to TCR tool development environment setup for those who would like to build
TCR tool locally.

<details><summary>Expand for details</summary>

### Clone TCR repository - `Required`

```shell
git clone https://github.com/murex/TCR.git
cd TCR
```

### Install Go SDK - `Required`

TCR is written in Go. This implies having Go compiler and tools installed on your machine.

Simply follow the instructions provided [here](https://go.dev/). Make sure to install **Go version 1.18** or higher.

### Install additional Go tools and utility packages

#### Go IDE - `Optional`

You can check this [link](https://www.tabnine.com/blog/top-7-golang-ides-for-go-developers/)
for a list of recommended IDEs supporting Go language.

#### Cobra library and tools - `Optional`

TCR Go command line options and parameters are managed using [Cobra](https://github.com/spf13/cobra).

The Cobra library download is already dealt with through Go Module dependencies, which means that in most situations you
will not need to worry about installing it.

In case you need to add or modify subcommands, options or parameters, you may want to use the Cobra Generator. In this
situation you can refer to
[Cobra Generator documentation](https://github.com/spf13/cobra/blob/master/user_guide.md#using-the-cobra-generator)

#### GoReleaser utility - `Optional`

New versions of TCR Go are released through [GoReleaser](https://goreleaser.com/).

You should not need it as long as you don't plan to release a new TCR Go version.

If you do, you can refer to [GoReleaser Installation Instructions](https://goreleaser.com/install/)
for installing it locally on your machine.

In most cases you will not even have to install it locally as TCR-Go new releases are built through
a [GoReleaser GitHub action](.github/workflows/goreleaser.yml).

#### golangci-lint package - `Optional`

We use the Go Linter aggregator [golangci-lint](https://golangci-lint.run/) to perform various static checks on TCR Go
code.

A [dedicated GitHub action](.github/workflows/goreleaser.yml) triggers execution of golangci-lint every time a new
TCR-Go version is being released.

Although not mandatory, we advise you to install it locally on your machine to check that your changes comply with
golangci-lint rules. Refer to [golangci-lint installation instructions](https://golangci-lint.run/usage/install/)
for installation.

Once golangci-lint is installed, you can run it from the root directory:

```shell
make lint
```

Both local run and GitHub Action use [this configuration file](.golangci.yml)

#### gotestsum utility - `Optional`

We use [gotestsum](https://github.com/gotestyourself/gotestsum) for running tests
with the possibility to generate a xunit-compatible test report.

Although not mandatory, we advise you to install it locally on your machine as it greatly improves
readability of test results.
Refer to [gotestsum's Install section](https://github.com/gotestyourself/gotestsum#install)
for installation.

Once gotestsum is installed, you can run make's test target from the root directory:

- For running all tests:

  ```shell
  make test
  ```
- For running short tests only:

  ```shell
  make test-short
  ```

#### Fyne toolkit - `Optional`

The GUI version of TCR-Go is built on top of [Fyne toolkit](https://fyne.io/) for all GUI-related stuff.

Refer to [Fyne Develop](https://developer.fyne.io/) for installation and usage instructions.

You will not need it as long as you're working on the TCR Command Line implementation only.

> ***Note about TCR GUI version***
>
> So far we are only releasing the command line version of TCR.
> We have not reached yet the stage where we could fully automate through a goreleaser GitHub Action
> the cross-compilation and link of TCR with Fyne and its dependencies.
> In the meantime if you wish to give TCR GUI a try, you still can build and run it locally on your machine
> (refer to `Build TCR executable` section below)

### Build TCR executable

To build TCR locally on your machine, simply type the following from the root directory:

```shell
make
```

This command generates by default both TCR CLI (in [tcr-cli](./tcr-cli) directory)
and TCR GUI (in [tcr-gui](./tcr-gui) directory) executables, as well as the command help pages (in [doc](./doc)
directory).

<details><summary>To build TCR CLI only</summary>

Either run the following command from the root directory:

```shell
make -C ./tcr-cli
```

Or run make from [tcr-cli](./tcr-cli) directory:

```shell
cd tcr-cli
make
```

</details>

<details><summary>To build TCR GUI only</summary>

Either run the following command from the root directory:

```shell
make -C ./tcr-gui
```

Or run make from [tcr-gui](./tcr-gui) directory:

```shell
cd tcr-gui
make
```

</details>

<details><summary>To generate TCR command markdown documentation</summary>

```shell
make doc
```

</details>

</details>

## Releasing a new TCR version

We use [GoReleaser](https://goreleaser.com/) for releasing new TCR versions.

<details><summary>Expand for details</summary>

### Versioning Rules

TCR release versions comply with [Semantic Versioning rules](https://semver.org/).

### Release Branch

All TCR releases are published on GitHub's `main` branch.

### Release Preparation

- [ ] Cleanup Go module dependencies: `make tidy`
- [ ] Run static checks and fix any non-conformity: `make lint`
- [ ] Verify that the build works: `make build`
- [ ] Verify that all tests pass: `make test`
- [ ] Commit all changes on the `main` branch
- [ ] Push the changes to GitHub and [wait until all GitHub Actions are green](https://github.com/murex/TCR/actions)
- [ ] Create the release tag: `git tag -a vX.Y.Z`
- [ ] Verify that everything is ready for GoReleaser: `make snapshot`

### Releasing

The creation of the new release is triggered by pushing the newly created release tag to GitHub repository

- [ ] Push the release tag: `git push origin vX.Y.Z`
- [ ] [Wait until all GitHub Actions are green](https://github.com/murex/TCR/actions)
- [ ] Open [TCR Release page](https://github.com/murex/TCR/releases) and verify that the new release is there
- [ ] Edit the release notes document, and insert a `Summary` section at the top, listing the main changes included in
  this release. You may take a look at previous release notes if unsure what should go in there.

</details>

## How to Contribute?

TCR tool is still at an early stage of development, and there are still plenty of features that we would like to add in
the future, such as additional languages and toolchains, collaboration utilities, etc.

Refer to [CONTRIBUTING.md](./CONTRIBUTING.md) for general contribution agreement and guidelines.

## License

The `TCR` application and the accompanying materials are made available under the terms of the [MIT License](LICENSE.md)
which accompanies this distribution, and is available at the
[Open Source site](https://opensource.org/licenses/MIT).

## Contributors

<table>
<tr>
    <td align="center" style="word-wrap: break-word; width: 150.0; height: 150.0">
        <a href=https://github.com/mengdaming>
            <img src=https://avatars.githubusercontent.com/u/1313765?v=4 width="100;"  style="border-radius:50%;align-items:center;justify-content:center;overflow:hidden;padding-top:10px" alt=Damien Menanteau/>
            <br />
            <sub style="font-size:14px"><b>Damien Menanteau</b></sub>
        </a>
    </td>
    <td align="center" style="word-wrap: break-word; width: 150.0; height: 150.0">
        <a href=https://github.com/aatwi>
            <img src=https://avatars.githubusercontent.com/u/11088496?v=4 width="100;"  style="border-radius:50%;align-items:center;justify-content:center;overflow:hidden;padding-top:10px" alt=Ahmad Atwi/>
            <br />
            <sub style="font-size:14px"><b>Ahmad Atwi</b></sub>
        </a>
    </td>
</tr>
</table>

