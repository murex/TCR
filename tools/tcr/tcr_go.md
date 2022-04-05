# TCR - Go Flavor

# About

This is a [Go](https://golang.org/) implementation of TCR.

## Usage

### Prerequisites

- Supported Operating Systems: macOS, Linux, Windows
- Have a [git client](https://git-scm.com/downloads) installed
- Have [curl](https://curl.se/download.html) command line utility installed
- Have a fully operational Java or C++ development environment
- **Linux only** - have [kdialog](https://apps.kde.org/kdialog/) installed
  <details><summary>Details</summary>
  TCR leverages on the OS desktop notification framework to send timer reminders.
  On Linux, it relies on kdialog for that purpose. Make sure to have it installed
  in order to be able to see TCR's timer notifications.
  </details>

### Running TCR Go

<details><summary>C++ version of the kata</summary> 

```shell
cd cpp
./tcrw
```
</details>
<details><summary>Java version of the kata</summary> 

```shell
cd java
./tcrw
```

</details>

### Additional Options

TCR Go comes with various command line options.
Details related to these options can be accessed through command line help:

```shell
./tcrw help
```

### Main menu

After starting the script, you will see a menu that looks like the following:

```text
[TCR] Starting TCR version 0.8.0...
[TCR] Working directory is (...)
[TCR] Git auto-push is turned on
[TCR] Timer duration is 5m0s
[TCR] -------------------------------------------------------------------------
[TCR] Running in mob mode
[TCR] -------------------------------------------------------------------------
[TCR] Working Directory: (...)
[TCR] Language=java, Toolchain=gradle
[TCR] Running on git branch "xxxxx" with auto-push enabled
[TCR] -------------------------------------------------------------------------
[TCR] What shall we do?
[TCR]   D -> Driver role
[TCR]   N -> Navigator role
[TCR]   P -> Turn on/off git auto-push
[TCR]   Q -> Quit
[TCR]   ? -> List available options
```

If you're not familiar with the Driver and Navigator roles,
you can refer to [here](https://mobprogramming.org/mob-programming-basics/).

In TCR the driver is primarily the participant actively interacting with the keyboard,
while the navigators are all other participants.

### Driver role

- To take Driver role, hit `d` from the main menu.
- You remain with Driver role until you hit `q` to go back to the main menu.

When running with Driver role, the TCR process enters into action:
every time you save a file (either a source file or a test file), TCR triggers
a build, then runs all tests.

- If the build fails, TCR does not commit or revert anything. It goes back to the
  idle state and waits for you to fix compilation errors.
- If the build passes, TCR triggers the test execution.
- If all tests pass, TCR commits all changes to git.
- If one or more tests fail, TCR reverts all changes on source files, but leaves
  test files unchanged.

> ***Important***
> - __Make sure to turn off your IDE's auto-save mode while using TCR!!!__
    >   TCR constantly watches the file system, triggering builds,
    >     tests, commits and reverts as soon as it detects changes.
    >     For this reason, it does not get along well with your IDE's auto-save mode.
> - __There should not be more than one driver per branch at a time!__
    >     You will likely face occasional merge conflicts otherwise.

### Navigator role

- To take Navigator role, hit `n` from the main menu.
- You remain with Navigator role until you hit `q` to go back to the main menu.

When running with Navigator role, the script periodically pulls changes from the git repository
to your local clone. It does not push any change that you might make locally.

### Ending the script

- If you're running driver or navigator role, go back to the main menu by hitting `q`.
- Type `q` a second time to end the script.

## Contribution Workflow

TCR utility provides the basic mechanics to run TCR and synchronize files between contributors,
however it does not replace collaboration discipline.

So make sure that at the end of each driver rotation:

- The former driver switches back to Navigator role: `q` + `n`.
- The new driver switches to Driver role: `q` + `d`.

Other contributors have nothing to do as long as they remain navigators.

## Command Line Options

The `tcrw` utility provides the following options:

<details><summary>Command line help</summary>

In order to display available options when launching TCR:

```shell
./tcrw --help
```

Once TCR is running, you can hit `?` to list the available options and their shortcuts

</details>
<details><summary>TCR version query</summary>

To display the version of TCR utility running locally:

```shell
./tcrw --version
```

</details>
<details><summary>TCR build information query</summary>

To display build information related to the TCR binary running locally:

```shell
./tcrw --info
```

</details>
<details><summary>Git auto-push switch</summary>

### When using TCR on your own

By default, TCR runs on your local clone only:
it does not push any change to the `origin` git repository.

This is the preferred way of using it when you're running TCR on your own.

When you're done with it, it's up to you to decide what you want to do with it (squash, push, revert, etc.)

### When using TCR in pair or in mob

When using TCR together with others, sharing changes regularly becomes important.

For this situation, the script provides the command line option `-p` (or `--auto-push`).

With this option enabled, when in driver mode, the script performs a `git push` to origin
after every `git commit`.
This allows all participants running the script in Navigator mode to get the changes as soon as they
are committed by the Driver.

```shell
./tcrw --auto-push
```

Once TCR is running, you can toggle on and off git auto-push option by typing `p`

</details>
<details><summary>Toolchain selection</summary>

TCR can use different toolchains when running build and test.

Here are the toolchains currently supported for each language.

| Language | Toolchains                                         | Default        |
|----------|----------------------------------------------------|----------------|
| Java     | gradle, gradle-wrapper, maven, maven-wrapper, make | gradle-wrapper |
| C++      | cmake, cmake-kata, make                            | cmake-kata     |
| Go       | go-tools, make                                     | go-tools       |

Please note that you do not need to install any of these toolchains on your machine in order to use them.
We provide the wrappers allowing to download and run them in the context of the kata.

For example, if you prefer using Maven instead of Gradle wrapper when running the TCR script for a kata in Java:

```shell
./tcrw --toolchain maven
```

</details>
<details><summary>Mob timer</summary>

When running TCR as a driver in mob mode, TCR automatically starts a countdown timer.
Its purpose is to notify the driver when it's time to hand over the driver role to another
participant.

### Changing the timer duration

The default timer duration is 5 minutes.

You can change this value when starting TCR as follows:

For a 10-minute timer:
```shell
./tcrw --duration 10m
```

### Disabling the timer

If you do not want to use the timer, you can turn it off by setting its duration to 0m when starting TCR.

```shell
./tcrw --duration 0m
```

### Querying the timer status

You can check the timer status at any time after you started running in driver mode.

Simply type `t` in the terminal to display time already spent and time remaining.

> ***Notes***
>
> - This shortcut is only active when in mob mode with the driver role running
> - There is no timer in solo mode (`./tcrw solo`)

</details>

