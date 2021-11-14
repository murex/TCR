## tcr solo

Run TCR in solo mode

### Synopsis


When used in "solo" mode, TCR only commits changes locally.
It never pushes or pulls to a remote repository.

This subcommand runs directly in the terminal (no GUI).

```
tcr solo [flags]
```

### Options

```
  -h, --help   help for solo
```

### Options inherited from parent commands

```
  -b, --base-dir string     indicate the base directory from which TCR is running
  -d, --duration duration   set the duration for role rotation countdown timer (default 5m0s)
  -i, --info                show build information about TCR application
  -t, --toolchain string    indicate the toolchain to be used by TCR
```

### SEE ALSO

* [tcr](tcr.md)	 - TCR (Test && Commit || Revert)

