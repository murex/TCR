## tcr config save

Save TCR configuration

### Synopsis


config save subcommand saves TCR configuration into a file.

This subcommand does not start TCR engine.

```
tcr config save [flags]
```

### Options

```
  -h, --help   help for save
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

* [tcr config](tcr_config.md)	 - Manage TCR configuration

