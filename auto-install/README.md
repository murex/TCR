# TCR Wrapper Auto-Installer

This directory contains a one-liner script for automatically setting up TCR (Test && Commit || Revert) in your project.

## Purpose

The auto-install directory provides a simple, one-command setup for TCR that:
- Downloads the `tcrw` wrapper script to your current directory
- Creates a `tcr` directory structure
- Downloads the version configuration file
- Provides clear usage instructions

## Quick Setup

To set up TCR in your current directory, run:

```bash
curl -sSL https://raw.githubusercontent.com/murex/TCR/main/auto-install/setup.sh | bash
```

This will:
- Download TCR wrapper `tcrw` to your current directory
- Create a `tcr/` subdirectory and set the latest available TCR version in `version.txt`
- Provide brief usage instructions

## After Setup

Once the setup is complete, you can start using TCR by running:

```bash
./tcrw
```

For help and available options:

```bash
./tcrw --help
```

## Requirements

The setup script requires `curl`to be installed on your system.

## More Information

For complete documentation and examples, visit the main TCR repository: https://github.com/murex/TCR
