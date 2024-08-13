## tcr stats

Print TCR stats

### Synopsis


TCR stats subcommand prints out TCR usage stats based on commit history.

The commit history is retrieved for the repository containing
TCR base directory (cf. -b option). The branch is the current working
branch set for this repository.

The following stats are reported:

- First commit date and time
- Last commit date and time
- Number of commits
- Number of passing commits (absolute value and percentage) (*)
- Number of failing commits, (absolute value and percentage) (*)
- Time span between the first and last commit
- Time in green: total time where all tests passed (absolute value and percentage) (*)
- Time in red: total time where one or more tests failed (absolute value and percentage) (*)
- Time between commits (minimum, average and maximum values)
- Changes per commit (src): number of lines of source code changed per commit (minimum, average and maximum values)
- Changes per commit (test): number of lines of test code changed per commit (minimum, average and maximum values)
- Passing tests count evolution (values for first and last commit)
- Failing tests count evolution (values for first and last commit) (*)
- Skipped tests count evolution (values for first and last commit)
- Test execution duration cumulated for all tests (values for first and last commit)

> (*) These metrics are relevant only if TCR commit history was created while running TCR with "commit-failures" option.
> Without this option there is no record of test failures in TCR commit history, thus:
> - "Number of passing commits" and "time in green" will always be at 100%
> - "Number of failing commits" and "time in red" will always be at 0%
> - "Failing tests" will always be at 0

This subcommand does not start TCR engine.

```
tcr stats [flags]
```

### Options

```
  -h, --help   help for stats
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

