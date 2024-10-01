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
- VCS environment
- Auto-push settings
- Mob timer settings (for driver role)
- Polling period settings (for navigator role)

The return code of TCR "check" is one of the following:

| RC  | Meaning                                                                    |
|-----|----------------------------------------------------------------------------|
| 0   | All checks passed without any warning or error                             |
| 1   | One or more warnings were raised. This should not prevent TCR from running |
| 2   | One or more errors were raised. TCR will not be able to run properly       |


```
tcr check [flags]
```

### Options

```
  -h, --help   help for check
```

### Options inherited from parent commands

```
  -p, --auto-push               enable VCS push after every commit
  -b, --base-dir string         indicate the directory from which TCR is looking for files (default: current directory)
  -c, --config-dir string       indicate the directory where TCR configuration is stored (default: current directory)
  -d, --duration duration       set the duration for role rotation countdown timer
  -l, --language string         indicate the programming language to be used by TCR
  -m, --message-suffix string   indicate text to append at the end of TCR commit messages (ex: "[#1234]")
  -o, --polling duration        set VCS polling period when running as navigator
  -P, --port-number int         indicate port number used by TCR HTTP server in web mode (experimental) (default: 8483)
  -t, --toolchain string        indicate the toolchain to be used by TCR
  -T, --trace string            indicate trace options. Recognized values: none (default), vcs or http
  -r, --variant string          indicate the variant to be used by TCR: relaxed (default), btcr, or introspective
  -V, --vcs string              indicate the VCS (version control system) to be used by TCR: git (default) or p4
  -w, --work-dir string         indicate the directory from which TCR is running (default: current directory)
```

### SEE ALSO

* [tcr](tcr.md)	 - TCR (Test && Commit || Revert)

