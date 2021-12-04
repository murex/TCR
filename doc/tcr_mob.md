## tcr mob

Run TCR in mob mode

### Synopsis


When used in "mob" mode, TCR ensures that any commit
is shared with other participants through calling git push-pull.

This subcommand runs directly in the terminal (no GUI).

```
tcr mob [flags]
```

### Options

```
  -h, --help   help for mob
```

### Options inherited from parent commands

```
  -p, --auto-push           enable git push after every commit
  -b, --base-dir string     indicate the base directory from which TCR is running
  -c, --config string       config file (default is $HOME/tcr.yaml)
  -d, --duration duration   set the duration for role rotation countdown timer
  -l, --language string     indicate the programming language to be used by TCR
  -o, --polling duration    set git polling period when running as navigator
  -t, --toolchain string    indicate the toolchain to be used by TCR
```

### SEE ALSO

* [tcr](tcr.md)	 - TCR (Test && Commit || Revert)

