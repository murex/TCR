## tcr retro

Generate retrospective template with stats

### Synopsis


TCR retro subcommand generates a retrospective template in markdown format prefilled with TCR execution info.
The markdown file is saved into the TCR base directory with the name 'tcr-retro.md'. 

The following information is included in the markdown:

- Average size of changes per passing commit 
- Average size of changes per failing commit

These stats are extracted for the repository containing TCR base directory (cf. -b option). 
The branch is the current working branch set for this repository.

This subcommand does not start TCR engine.

```
tcr retro [flags]
```

### Options

```
  -h, --help   help for retro
```

### Options inherited from parent commands

```
  -p, --auto-push               enable VCS push after every commit
  -b, --base-dir string         indicate the directory from which TCR is looking for files (default: current directory)
  -c, --config-dir string       indicate the directory where TCR configuration is stored (default: current directory)
  -d, --duration duration       set the duration for role rotation countdown timer
  -g, --git-remote string       name of the git remote repository to sync with (default: "origin")
  -l, --language string         indicate the programming language to be used by TCR
  -m, --message-suffix string   indicate text to append at the end of TCR commit messages (ex: "[#1234]")
  -o, --polling duration        set VCS polling period when running as navigator
  -P, --port-number int         indicate port number used by TCR HTTP server in web mode (experimental) (default: 8483)
  -t, --toolchain string        indicate the toolchain to be used by TCR
  -T, --trace string            indicate trace options. Recognized values: none (default), vcs or http
  -r, --variant string          indicate the variant to be used by TCR: relaxed (default), btcr, or introspective
  -V, --vcs string              indicate the VCS (version control system) to be used by TCR: git (default) or p4
  -w, --work-dir string         indicate the directory from which TCR is running (default: current directory)
```

### SEE ALSO

* [tcr](tcr.md)	 - TCR (Test && Commit || Revert)

