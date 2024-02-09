## tcr

TCR (Test && Commit || Revert)

### Synopsis


This application is a tool for practicing TCR (Test && Commit || Revert).
It can be used either in solo, or as a group within a mob or pair session.


```
tcr [flags]
```

### Options

```
  -p, --auto-push               enable VCS push after every commit
  -b, --base-dir string         indicate the directory from which TCR is looking for files (default: current directory)
  -f, --commit-failures         enable committing reverts on tests failure
  -c, --config-dir string       indicate the directory where TCR configuration is stored (default: current directory)
  -d, --duration duration       set the duration for role rotation countdown timer
  -h, --help                    help for tcr
  -l, --language string         indicate the programming language to be used by TCR
  -m, --message-suffix string   indicate text to append at the end of TCR commit messages (ex: "[#1234]")
  -o, --polling duration        set VCS polling period when running as navigator
  -P, --port-number int         indicate port number used by TCR HTTP server in web mode (experimental) (default: 8483)
  -t, --toolchain string        indicate the toolchain to be used by TCR
  -T, --trace string            indicate trace options. Recognized values: none or vcs
  -V, --vcs string              indicate the VCS (version control system) to be used by TCR: git (default) or p4
  -w, --work-dir string         indicate the directory from which TCR is running (default: current directory)
```

### SEE ALSO

* [tcr check](tcr_check.md)	 - Check TCR configuration and parameters and exit
* [tcr config](tcr_config.md)	 - Manage TCR configuration
* [tcr info](tcr_info.md)	 - Display TCR build information
* [tcr log](tcr_log.md)	 - Print the TCR commit history
* [tcr mob](tcr_mob.md)	 - Run TCR in mob mode
* [tcr one-shot](tcr_one-shot.md)	 - Run one TCR cycle and exit
* [tcr solo](tcr_solo.md)	 - Run TCR in solo mode
* [tcr stats](tcr_stats.md)	 - Print TCR stats
* [tcr web](tcr_web.md)	 - Run TCR with web user interface (experimental)

