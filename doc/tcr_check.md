## tcr check

Check TCR configuration and parameters and exit

### Synopsis


When used in "check" mode, TCR performs a series of verifications on provided parameters,
configuration and local environment, then exits.

Its purpose is to verify that TCR is ready to run. It does not trigger any TCR cycle execution.

The main checkpoints are organized in sections as follows:

- Configuration directory
- Base directory
- Work Directory
- Language settings
- Toolchain settings
- Git environment
- Auto-push settings
- Mob timer settings (for driver role)
- Polling period settings (for navigator role)

The return code of TCR "check" is one of the following:

| RC  | Meaning                                                                    |
|-----|----------------------------------------------------------------------------|
| 0   | All checks passed without any warning or error                             |
| 1   | One or more warnings were raised. This should not prevent TCR from running |
| 2   | One or more errors were raised. TCR will not be able to run properly       |

This subcommand runs directly in the terminal (no GUI).

```
tcr check [flags]
```

### Options

```
  -h, --help   help for check
```

### Options inherited from parent commands

```
  -p, --auto-push           enable git push after every commit
  -b, --base-dir string     indicate the directory from which TCR is looking for files (default: current directory)
  -f, --commit-failures     enable committing reverts on tests failure
  -c, --config-dir string   indicate the directory where TCR configuration is stored (default: current directory)
  -d, --duration duration   set the duration for role rotation countdown timer
  -l, --language string     indicate the programming language to be used by TCR
  -o, --polling duration    set git polling period when running as navigator
  -t, --toolchain string    indicate the toolchain to be used by TCR
  -w, --work-dir string     indicate the directory from which TCR is running (default: current directory)
```

### SEE ALSO

* [tcr](tcr.md)	 - TCR (Test && Commit || Revert)

