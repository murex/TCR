## tcr log

Print the TCR commit history

### Synopsis


TCR log subcommand prints out the TCR commit history.

The output format is similar to "git log" command's output format.

The commit history is retrieved for the repository containing
TCR base directory (cf. -b option). The branch is the current working
branch set for this repository.

Only TCR commits are printed. All other commits are filtered out.

This subcommand does not start TCR engine.

```
tcr log [flags]
```

### Options

```
  -h, --help   help for log
```

### Options inherited from parent commands

```
  -p, --auto-push           enable VCS push after every commit
  -b, --base-dir string     indicate the directory from which TCR is looking for files (default: current directory)
  -f, --commit-failures     enable committing reverts on tests failure
  -c, --config-dir string   indicate the directory where TCR configuration is stored (default: current directory)
  -d, --duration duration   set the duration for role rotation countdown timer
  -l, --language string     indicate the programming language to be used by TCR
  -o, --polling duration    set VCS polling period when running as navigator
  -t, --toolchain string    indicate the toolchain to be used by TCR
  -V, --vcs string          indicate the VCS (version control system) to be used by TCR: git (default) or p4
  -w, --work-dir string     indicate the directory from which TCR is running (default: current directory)
```

### SEE ALSO

* [tcr](tcr.md)	 - TCR (Test && Commit || Revert)

