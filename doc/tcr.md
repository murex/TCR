## tcr

TCR (Test && Commit || Revert)

### Synopsis


This application is a tool for practicing TCR (Test && Commit || Revert).
It can be used either in solo, or as a group within a mob or pair session.

This application runs within a terminal.

```
tcr [flags]
```

### Options

```
  -p, --auto-push           enable git push after every commit
  -b, --base-dir string     indicate the directory from which TCR is looking for files (default: current directory)
  -f, --commit-failures     enable committing reverts on tests failure
  -c, --config-dir string   indicate the directory where TCR configuration is stored (default: current directory)
  -d, --duration duration   set the duration for role rotation countdown timer
  -h, --help                help for tcr
  -l, --language string     indicate the programming language to be used by TCR
  -o, --polling duration    set git polling period when running as navigator
  -t, --toolchain string    indicate the toolchain to be used by TCR
  -w, --work-dir string     indicate the directory from which TCR is running (default: current directory)
```

### SEE ALSO

* [tcr check](tcr_check.md)	 - Check TCR configuration and parameters and exit
* [tcr config](tcr_config.md)	 - Manage TCR configuration
* [tcr info](tcr_info.md)	 - Display TCR build information
* [tcr mob](tcr_mob.md)	 - Run TCR in mob mode
* [tcr one-shot](tcr_one-shot.md)	 - Run one TCR cycle and exit
* [tcr solo](tcr_solo.md)	 - Run TCR in solo mode

