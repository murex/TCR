# Using TCR with PHP and PHPUnit

## Prerequisites

- macOS, Linux or Windows
- [git](https://git-scm.com/) client
- [curl](https://curl.se/download.html) command line utility
- [PHP](https://www.php.net/manual/en/install.php)
- [Composer](https://getcomposer.org/) dependency manager

## Instructions

### 1 - Open a terminal

> ***Note to Windows users***
>
> Use a **git bash** terminal for running the commands below.
> _Windows CMD and PowerShell are not supported_

### 2 - Setup environment with Composer

> ***Reminder***: the command below should be run from
> [examples/php-phpunit](.)
> directory

```shell
composer update
```

### 3 - Launch TCR

> ***Reminder***: the command below should be run from
> [examples/php-phpunit](.)
> directory
>

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
