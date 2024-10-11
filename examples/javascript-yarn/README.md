# Using TCR with Javascript and Yarn

## Prerequisites

- macOS, Linux or Windows
- [git](https://git-scm.com/) client
- [curl](https://curl.se/download.html) command line utility
- [Node](https://nodejs.org/en/download)
- [Yarn](https://classic.yarnpkg.com/lang/en/docs/install)
  <details><summary>Details</summary>

  You can install node directly or through nvm.

  </details>

## Instructions

### 1 - Open a terminal

> ***Note to Windows users***
>
> Use a **git bash** terminal for running the commands below.
> _Windows CMD and PowerShell are not supported_

### 2 - Install dependencies

> ***Reminder***: the command below should be run from
> [examples/javascript-yarn](.)
> directory

The example uses `corepack` when configuring `yarn` as the package manager to be used.
You may need to run the following command beforehand to enable it:

> ***Note***: depending on your environment you may need to run this command
> as an administrator (Windows) or with sudo (Linux and macOS)

```shell
corepack enable
```

To install the dependencies:

```shell
yarn install
```

### 3 - Launch TCR

> ***Reminder***: the command below should be run from
> [examples/javascript-yarn](.)
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
| `a` / `A` | Abort current command (when in driver role)  |
| `q` / `Q` | Quit current role / Quit TCR                 |
| `?`       | List available options                       |

### Additional Details

Refer to [tcr.md](../../doc/tcr.md) page for additional details and explanations about TCR
available subcommands and options
