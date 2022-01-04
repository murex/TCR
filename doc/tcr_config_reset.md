## tcr config reset

Reset TCR configuration

### Synopsis


config reset subcommand resets TCR configuration to default values.

This subcommand does not start TCR engine.

```
tcr config reset [flags]
```

### Options

```
  -h, --help   help for reset
```

### Options inherited from parent commands

```
  -p, --auto-push           enable git push after every commit
  -b, --base-dir string     indicate the base directory from which TCR is running (default: current directory)
  -c, --config-dir string   indicate the directory where TCR configuration is stored (default: current directory)
  -d, --duration duration   set the duration for role rotation countdown timer
  -l, --language string     indicate the programming language to be used by TCR
  -o, --polling duration    set git polling period when running as navigator
  -t, --toolchain string    indicate the toolchain to be used by TCR
```

### SEE ALSO

* [tcr config](tcr_config.md)	 - Manage TCR configuration

