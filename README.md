[![Go](https://github.com/murex/tcr/actions/workflows/go.yml/badge.svg)](https://github.com/murex/tcr/actions/workflows/go.yml)
[![Go Lint](https://github.com/murex/tcr/actions/workflows/golangci_lint.yml/badge.svg)](https://github.com/murex/tcr/actions/workflows/golangci_lint.yml)
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

Although probably a bit challenging to use when working on real production code, we found TCR to be quite beneficial
when used as a learning and practicing tool on the katas that we use when doing software craftsmanship coaching.

TCR enforces developing in baby steps, with a strong focus on always keeping the green light on tests. Having a TCR tool
feels a bit like having a coaching assistant constantly enforcing such practices!

We initially came up with a small shell script implementing this workflow, and decided to embed it in each of our katas
so that people can use it if they like. This was a way for us to quickly provide a usable TCR solution for. However,
shell scripts are not the best in class when it comes to maintainability and changeability.

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
3. Supported platforms: MacOS (x86_64), Linux (x86_64) and Windows (x86_64)
4. Supported languages: TCR tool currently works with Java and C++ out of the box (more to come in the future)
5. Have [Java JDK](https://jdk.java.net/archive/) 11 or later installed for java, or a working C++ compiler for C++
6. Build toolchains: [Gradle](https://gradle.org/) and [Maven](https://maven.apache.org/index.html) for
   Java, [CMake](https://cmake.org/) for C++
7. If you're using an IDE, make sure that **your IDE's auto-save is turned off**. TCR is constantly watching for file
   changes in your code, and for this reason it does not get along well with IDE's auto-save features.

### Directory Layout

The TCR tool needs to know where both production code and test code are located.

For this reason it makes some assumptions on the code directory structure.

<details>
  <summary>Directory Layout for Java</summary>

- TCR expects that the root directory for java code is named `java`
- Under the root directory, TCR assumes that the code
  follows [Maven's Standard Directory Layout](https://maven.apache.org/guides/introduction/introduction-to-the-standard-directory-layout.html)
- If you intend to use Gradle as the build toolchain, make sure to
  install [Gradle Wrapper](https://docs.gradle.org/current/userguide/gradle_wrapper.html) under the root directory
- If you intend to use Maven as the build toolchain, install [Maven Wrapper](https://github.com/takari/maven-wrapper)
  under the root directory

In case of doubt you can have a look at [this example](./tcr-engine/testdata/java)

</details>

<details>
  <summary>Directory Layout for C++</summary>

- TCR expects that the root directory for C++ code is named `cpp`
- Under the root directory, TCR assumes that the code is organized into 4 subdirectories:
    - `src` - Source code
    - `include` - Header files
    - `test` - Test code
    - `build` - Build files and directories
- TCR expects to find a `CMakeLists.txt` under the root directory

We advise you to have a look at [this example](./tcr-engine/testdata/cpp) to get a better idea of what TCR tool is
expecting to find.

The provided script [cpp_easy_setup.sh](./tcr-engine/testdata/cpp/cpp_easy_setup.sh) should help you have everything
setup and running before running TCR. Among other things it downloads under the `build` directory a working version
of `CMake` that will then be used by the TCR tool.

</details>

### Running TCR

<details>
  <summary>On MacOS</summary>

1. Download the latest version of TCR for Darwin from [here](https://github.com/murex/tcr/releases)

2. Extract TCR executable

    ```shell
    $ # Replace "0.6.0" with the appropriate version
    $ tar zxf tcr_0.6.0_Darwin_x86_64.tar.gz
    ```

3. Launch TCR

    ```shell
    $ ./tcr -b <path to the code root directory>
    ```

</details>

<details>
  <summary>On Linux</summary>

1. Download the latest version of TCR for Linux from [here](https://github.com/murex/tcr/releases)

2. Extract TCR executable

    ```shell
    $ # Replace "0.6.0" with the appropriate version
    $ tar zxf tcr_0.6.0_Linux_x86_64.tar.gz
    ```

3. Launch TCR

    ```shell
    $ ./tcr -b <path to the code root directory>
    ```

</details>

<details>
  <summary>On Windows</summary>

1. Download the latest version of TCR for Windows from [here](https://github.com/murex/tcr/releases)

2. Extract TCR executable

    ```shell
    $ # Replace "0.6.0" with the appropriate version
    $ tar zxf tcr_0.6.0_Windows_x86_64.tar.gz
    ```

3. Launch TCR

    ```shell
    $ ./tcr.exe -b <path to the code root directory>
    ```

</details>

#### Command line help (all platforms)

Refer to [here](./doc/tcr.md) for TCR command line help and additional options.

## How to Contribute?

TCR tool is still at an early stage of development, and there are still plenty of features that we would like to add in
the future, such as additional languages and toolchains, collaboration utilities, etc.

Refer to [CONTRIBUTING.md](./CONTRIBUTING.md) for general contribution agreement and guidelines.

### Development environment and tools

`TODO - Add instructions related to development environment and tools`

## License

The `tcr` application and the accompanying materials are made available under the terms of the [MIT License](LICENSE.md)
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
</tr>
</table>

