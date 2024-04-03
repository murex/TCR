# Using TCR with Elixir and Mix

## Prerequisites

- macOS, Linux or Windows
- [git](https://git-scm.com/) client
- [curl](https://curl.se/download.html) command line utility
- [elixir](https://elixir-lang.org/install.html)

## Instructions

### 1 - Open a terminal

> ***Note to Windows users***
>
> Use a **git bash** terminal for running the commands below.
> _Windows CMD and PowerShell are not supported_

### 2 - Install dependencies

> ***Reminder***: the command below should be run from
> [examples/elixir-mix](.)
> directory

```shell
mix deps.get
```

### 3 - Launch TCR

> ***Reminder***: the command below should be run from
> [examples/elixir-mix](.)
> directory

From the built-in terminal:

```shell
./tcrw
```

### Cheat Sheet

Here are the main shortcuts available once TCR utility is running:

| Shortcut  | Description                                  |
|-----------|----------------------------------------------|
| `o` / `O` | Open in browser (with `web` subcommand only) |
| `d` / `D` | Enter driver role                            |
| `n` / `N` | Enter navigator role                         |
| `t` / `T` | Query timer status                           |
| `p` / `P` | Toggle on/off git auto-push                  |
| `l` / `L` | Pull from remote                             |
| `s` / `S` | Push to remote                               |
| `q` / `Q` | Quit current role / Quit TCR                 |
| `?`       | List available options                       |

### Additional Details

Refer to [tcr.md](../../doc/tcr.md) page for additional details and explanations about TCR
available subcommands and options
