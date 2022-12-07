# Using TCR with Kotlin and Maven

## Prerequisites

- macOS, Linux or Windows
- [git](https://git-scm.com/) client
- [curl](https://curl.se/download.html) command line utility
- [kotlin](https://kotlinlang.org/docs/home.html)
- [java SDK](https://www.oracle.com/java/technologies/downloads/) with `JAVA_HOME` environment variable properly set up
- [maven](https://maven.apache.org/) (previously installed and available from the command line)

## Instructions

### 1 - Open a terminal

> ***Note to Windows users***
>
> Use a **git bash** terminal for running the commands below.
> _Windows CMD and PowerShell are not supported_

### 2 - Launch TCR

> ***Reminder***: the command below should be run from
> [examples/kotlin-maven](.)
> directory

From the built-in terminal:

```shell
./tcrw
```

### Cheat Sheet

Here are the main shortcuts available once TCR utility is running:

| Shortcut  | Description                                   |
|-----------|-----------------------------------------------|
| `d` / `D` | Enter driver role (from main menu)            |
| `n` / `N` | Enter navigator role (from main menu)         |
| `p` / `P` | Toggle on/off git auto-push (from main menu)  |
| `l` / `L` | Pull from remote (from main menu)             |
| `s` / `S` | Push to remote (from main menu)               |
| `q` / `Q` | Quit current role - Quit TCR (from main menu) |
| `t` / `T` | Query timer status (from driver role only)    |
| `?`       | List available options                        |

### Additional Details

Refer to [tcr.md](../../doc/tcr.md) page for additional details and explanations about TCR
available subcommands and options
