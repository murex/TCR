# Using TCR with Go and gotestsum

## Prerequisites

- macOS, Linux or Windows
- [git](https://git-scm.com/) client
- [curl](https://curl.se/download.html) command line utility
- [Go SDK](https://go.dev/dl/)
- [gotestsum](https://github.com/gotestyourself/gotestsum) utility

## Instructions

### 1 - Open a terminal

> ***Note to Windows users***
>
> Use a **git bash** terminal for running the commands below.
> _Windows CMD and PowerShell are not supported_

### 2 - Launch TCR

> ***Reminder***: the command below should be run from
> [examples/go-gotestsum](.)
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
