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
  -p, --auto-push               enable VCS push after every commit
  -b, --base-dir string         indicate the directory from which TCR is looking for files (default: current directory)
  -f, --commit-failures         enable committing reverts on tests failure
  -c, --config-dir string       indicate the directory where TCR configuration is stored (default: current directory)
  -d, --duration duration       set the duration for role rotation countdown timer
  -l, --language string         indicate the programming language to be used by TCR
  -m, --message-suffix string   indicate text to append at the end of TCR commit messages (ex: "[#1234]")
  -o, --polling duration        set VCS polling period when running as navigator
  -t, --toolchain string        indicate the toolchain to be used by TCR
  -V, --vcs string              indicate the VCS (version control system) to be used by TCR: git (default) or p4
  -w, --work-dir string         indicate the directory from which TCR is running (default: current directory)
```

### SEE ALSO

* [tcr config](tcr_config.md)	 - Manage TCR configuration

