## tcr config

Manage TCR configuration

### Synopsis


TCR config subcommand provides management of TCR configuration.

This subcommand does not start TCR engine.

```
tcr config [flags]
```

### Options

```
  -h, --help   help for config
```

### Options inherited from parent commands

```
  -p, --auto-push               enable VCS push after every commit
  -b, --base-dir string         indicate the directory from which TCR is looking for files (default: current directory)
  -F, --commit-failures         enable committing reverts on tests failure
  -c, --config-dir string       indicate the directory where TCR configuration is stored (default: current directory)
  -d, --duration duration       set the duration for role rotation countdown timer
  -f, --flavor string           indicate the flavor to be used by TCR: nice (default) or original
  -l, --language string         indicate the programming language to be used by TCR
  -m, --message-suffix string   indicate text to append at the end of TCR commit messages (ex: "[#1234]")
  -o, --polling duration        set VCS polling period when running as navigator
  -P, --port-number int         indicate port number used by TCR HTTP server in web mode (experimental) (default: 8483)
  -t, --toolchain string        indicate the toolchain to be used by TCR
  -T, --trace string            indicate trace options. Recognized values: none (default), vcs or http
  -V, --vcs string              indicate the VCS (version control system) to be used by TCR: git (default) or p4
  -w, --work-dir string         indicate the directory from which TCR is running (default: current directory)
```

### SEE ALSO

* [tcr](tcr.md)	 - TCR (Test && Commit || Revert)
* [tcr config reset](tcr_config_reset.md)	 - Reset TCR configuration
* [tcr config save](tcr_config_save.md)	 - Save TCR configuration
* [tcr config show](tcr_config_show.md)	 - Show TCR configuration

