## tcr one-shot

Run one TCR cycle and exit

### Synopsis


When used in "one-shot" mode, TCR executes instantaneously one TCR cycle without waiting
for any file system change, then exits.

The return code of TCR "one-shot" is one of the following:

| RC  | Meaning                                                                        |
|-----|--------------------------------------------------------------------------------|
| 0   | Build and Test Passed and changes were successfully committed                  |
| 1   | Build failed                                                                   |
| 2   | Build passed, one or more tests failed, and changes were successfully reverted |
| 3   | Error in configuration or parameter values                                     |
| 4   | Error while interacting with git                                               |
| 5   | Any other error                                                                |

This subcommand runs directly in the terminal (no GUI).

```
tcr one-shot [flags]
```

### Options

```
  -h, --help   help for one-shot
```

### Options inherited from parent commands

```
  -p, --auto-push           enable git push after every commit
  -b, --base-dir string     indicate the directory from which TCR is looking for files (default: current directory)
  -c, --config-dir string   indicate the directory where TCR configuration is stored (default: current directory)
  -d, --duration duration   set the duration for role rotation countdown timer
  -l, --language string     indicate the programming language to be used by TCR
  -o, --polling duration    set git polling period when running as navigator
  -t, --toolchain string    indicate the toolchain to be used by TCR
  -w, --work-dir string     indicate the directory from which TCR is running (default: current directory)
```

### SEE ALSO

* [tcr](tcr.md)	 - TCR (Test && Commit || Revert)

