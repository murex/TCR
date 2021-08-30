#!/usr/bin/env bash

set -eu

base_dir="$(cd "$(dirname -- "$0")" && pwd)"

cd "${base_dir}"

if ! type fyne-cross >/dev/null 2>/dev/null; then
  echo "fyne-cross command is required to run this script"
  echo "Refer to https://github.com/fyne-io/fyne-cross for installation and setup"
  exit 1
fi

fyne-cross darwin --app-id tcr --app-version 0.3.1 -arch=*
fyne-cross linux --app-id tcr --app-version 0.3.1 -arch=*
fyne-cross windows --app-id tcr --app-version 0.3.1 -arch=*


