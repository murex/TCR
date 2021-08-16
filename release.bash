#!/usr/bin/env bash

set -eu

base_dir="$(cd "$(dirname -- "$0")" && pwd)"

cd "${base_dir}"

if ! type goreleaser >/dev/null 2>/dev/null; then
  echo "goreleaser command is required to run this script"
  echo "Refer to https://goreleaser.com/ for installation and setup"
  exit 1
fi

goreleaser release --rm-dist
