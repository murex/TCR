# Using TCR with Python and Pytest

## Prerequisites

- macOS, Linux or Windows
- [git](https://git-scm.com/) client
- [curl](https://curl.se/download.html) command line utility
- [python](https://www.python.org/downloads/)

## Instructions

### 1 - Open a terminal

> ***Note to Windows users***
>
> Use a **git bash** terminal for running the commands below.
> _Windows CMD and PowerShell are not supported_

### 2 - Start python virtual environment

> ***Reminder***: the command below should be run from
> [examples/python-pytest](.)
> directory

Although not mandatory, this step allows to run in a virtual environment in order
to prevent interferences with other python projects that you may have on your machine.

From the built-in terminal:

```shell
./start_python_venv.sh
```

***Note:*** This command starts the python virtual environment in a sub-shell. You can end it
at any time by typing `exit`.

### 3 - Launch TCR

> ***Reminder***: the command below should be run from
> [examples/python-pytest](.)
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
