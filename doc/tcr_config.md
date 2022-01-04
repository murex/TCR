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
  -p, --auto-push           enable git push after every commit
  -b, --base-dir string     indicate the base directory from which TCR is running (default: current directory)
  -c, --config-dir string   indicate the directory where TCR configuration is stored (default: current directory)
  -d, --duration duration   set the duration for role rotation countdown timer
  -l, --language string     indicate the programming language to be used by TCR
  -o, --polling duration    set git polling period when running as navigator
  -t, --toolchain string    indicate the toolchain to be used by TCR
```

### SEE ALSO

* [tcr](tcr.md)	 - TCR (Test && Commit || Revert)
* [tcr config reset](tcr_config_reset.md)	 - Reset TCR configuration
* [tcr config save](tcr_config_save.md)	 - Save TCR configuration
* [tcr config show](tcr_config_show.md)	 - Show TCR configuration

