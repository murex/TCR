[![Go](https://github.com/murex/tcr/actions/workflows/go.yml/badge.svg)](https://github.com/murex/tcr/actions/workflows/go.yml)
[![Go lint](https://github.com/murex/tcr/actions/workflows/golangci_lint.yml/badge.svg)](https://github.com/murex/tcr/actions/workflows/golangci_lint.yml)
[![NPM build and test](https://github.com/murex/TCR/actions/workflows/npm.yml/badge.svg)](https://github.com/murex/TCR/actions/workflows/npm.yml)
[![sonarcloud](https://sonarcloud.io/api/project_badges/measure?project=murex_TCR&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=murex_TCR)
[![Coveralls](https://coveralls.io/repos/github/murex/TCR/badge.svg?branch=main)](https://coveralls.io/github/murex/TCR?branch=main)
[![goreleaser](https://github.com/murex/tcr/actions/workflows/go_releaser.yml/badge.svg)](https://github.com/murex/tcr/actions/workflows/go_releaser.yml)
[![Add contributors](https://github.com/murex/tcr/actions/workflows/add_contributors.yml/badge.svg)](https://github.com/murex/tcr/actions/workflows/add_contributors.yml)

# TCR - Test && Commit || Revert

_A Go implementation of TCR, for practicing baby-steps development, and much more!_

## What is this?

TCR is a programming workflow, standing for **Test && Commit || Revert**.

Kent Beck, Oddmund Strømme, Lars Barlindhaug and Ole Johannessen came up with this concept
as described in [this post](https://medium.com/@kentbeck_7670/test-commit-revert-870bbd756864).

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

1. Have [git](https://git-scm.com/) or [Perforce client](https://www.perforce.com/downloads/helix-visual-client-p4v)
   installed on your machine
2. Have a clone of the git repository or a client view of the Perforce depot containing the code you intend to work on
3. TCR can run on macOS, Linux and Windows. Refer to [TCR releases page](https://github.com/murex/TCR/releases)
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

| Language   | Default Toolchain | Compatible Toolchains                                |
|------------|-------------------|------------------------------------------------------|
| cpp        | cmake             | cmake bazel make                                     |
| csharp     | dotnet            | dotnet bazel make                                    |
| elixir     | mix               | mix                                                  |
| go         | go-tools          | go-tools gotestsum bazel make                        |
| haskell    | stack             | stack                                                |
| java       | gradle-wrapper    | gradle gradle-wrapper maven maven-wrapper bazel make |
| javascript | yarn              | yarn bazel make                                      |
| kotlin     | gradle-wrapper    | gradle gradle-wrapper maven maven-wrapper bazel make |
| php        | phpunit           | phpunit                                              |
| python     | pytest            | pytest bazel make                                    |
| rust       | cargo             | cargo nextest                                        |
| scala      | sbt               | sbt                                                  |
| typescript | yarn              | yarn bazel make                                      |

### TCR Variants

TCR tool can run several variants of the TCR workflow, inspired by [this blog post](https://medium.com/@tdeniffel/tcr-variants-test-commit-revert-bf6bd84b17d3)
by Thomas Deniffel.

Some are great as pedagogic tools, some are better for day-to-day use.

The default variant is "The Relaxed". You can refer to [this page](./variants-doc/tcr_variants.md)
for further details on available variants.

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

  ```text
  <configuration directory>/
  └── .tcr/
      ├── config.yml             - contains all TCR configuration settings
      ├── language/              - subdirectory containing all language configurations
      │   ├── cpp.yml            - configuration for C++ language
      │   ├── java.yml           - configuration for java language
      │   └── etc.
      └── toolchain/             - subdirectory containing all toolchain configurations
          ├── cmake.yml          - configuration for cmake toolchain
          ├── gradle.yml         - configuration for gradle toolchain
          ├── gradle-wrapper.yml - configuration for gradle wrapper toolchain
          └── etc.
  ```

</details>

### Examples

Refer to the [examples](examples/README.md) directory on how to set up and run
TCR for various language/toolchain combinations.

### Running TCR

#### Operating Systems

<details>
  <summary>On MacOS</summary>

1. Download the latest version of TCR for Darwin from [here](https://github.com/murex/TCR/releases)

2. Extract TCR executable (replace with the appropriate version and architecture)

    ```shell
    tar zxf tcr_1.3.0_Darwin_x86_64.tar.gz
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
    tar zxf tcr_1.3.0_Linux_x86_64.tar.gz
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
    tar zxf tcr_1.3.0_Windows_x86_64.tar.gz
    ```

3. Launch TCR

    ```shell
    ./tcr.exe -b <base-directory> -w <work-directory> -l <language> -t <toolchain>
    ```

</details>

#### Version Control Systems

<details>
  <summary>Git</summary>

TCR uses git by default. There is no need to specify anything particular to use git.

> ***Note: TCR and git commits signing***
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

</details>

<details>
  <summary>Perforce</summary>

Before running TCR with Perforce, make sure that the P4 Client is properly configured.

To use TCR with Perforce, you'll need to add the `--vcs=p4` to the command line (you can also set this up in
the `.tcr/config.yml`, see the [Configuration Directory](#configuration-directory) section below).

> ***Note: Perforce limitations***
>
> At the moment, TCR over Perforce is still in the experimentation phase. It does not yet support all the options
> available with git.
> Here are the main limitations:
>
> - Option `--commit-failures` or `-f` is not supported
> - Option `--auto-push` or `-p` has no meaning with Perforce and is ignored
> - Sub-command `log` is not supported
> - Sub-command `stats` is not supported

</details>

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

### Using TCR's embedded web interface `experimental`

Since version `1.0.0`, TCR comes with an embedded web interface that can be used
to monitor the TCR cycle, control driver and navigator roles,
and show the countdown timer when the driver mode active.

To use it, start TCR using the `web` subcommand:

```shell
./tcr web
```

Once TCR is running, you can open the web interface in your browser
by typing the `O` shortcut in the terminal.

TCR runs its internal web server on port `8483` by default.
You can change the port through using the `-P` (or `--port-number`) command line option.

### Command line help (all platforms)

Refer to [here](./doc/tcr.md) for TCR command line help and additional options.

## Building, testing and releasing TCR

Refer to [TCR development documentation](./dev-doc/README.md) for details.

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
    <td align="center" style="word-wrap: break-word; width: 150.0; height: 150.0">
        <a href=https://github.com/philou>
            <img src=https://avatars.githubusercontent.com/u/23983?v=4 width="100;"  style="border-radius:50%;align-items:center;justify-content:center;overflow:hidden;padding-top:10px" alt=Philippe Bourgau/>
            <br />
            <sub style="font-size:14px"><b>Philippe Bourgau</b></sub>
        </a>
    </td>
    <td align="center" style="word-wrap: break-word; width: 150.0; height: 150.0">
        <a href=https://github.com/YifangDONG>
            <img src=https://avatars.githubusercontent.com/u/15253963?v=4 width="100;"  style="border-radius:50%;align-items:center;justify-content:center;overflow:hidden;padding-top:10px" alt=Yifang DONG/>
            <br />
            <sub style="font-size:14px"><b>Yifang DONG</b></sub>
        </a>
    </td>
    <td align="center" style="word-wrap: break-word; width: 150.0; height: 150.0">
        <a href=https://github.com/rwilsonmurex>
            <img src=https://avatars.githubusercontent.com/u/161576431?v=4 width="100;"  style="border-radius:50%;align-items:center;justify-content:center;overflow:hidden;padding-top:10px" alt=rwilsonmurex/>
            <br />
            <sub style="font-size:14px"><b>rwilsonmurex</b></sub>
        </a>
    </td>
</tr>
</table>

