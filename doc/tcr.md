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
  -b, --base-dir string     indicate the base directory from which TCR is running
  -d, --duration duration   set the duration for role rotation countdown timer (default 5m0s)
  -h, --help                help for tcr
  -i, --info                show build information about TCR application
  -t, --toolchain string    indicate the toolchain to be used by TCR
```

### SEE ALSO

* [tcr mob](tcr_mob.md)	 - Run TCR in mob mode
* [tcr solo](tcr_solo.md)	 - Run TCR in solo mode

