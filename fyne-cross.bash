#!/usr/bin/env bash

set -eu

base_dir="$(cd "$(dirname -- "$0")" && pwd)"

cd "${base_dir}"

if ! type go >/dev/null 2>/dev/null; then
  echo "go is required to run this script"
  echo "Refer to https://golang.org/ for installation and setup"
  exit 1
fi

if ! type fyne-cross >/dev/null 2>/dev/null; then
  echo "fyne-cross command is required to run this script"
  echo "Refer to https://github.com/fyne-io/fyne-cross for installation and setup"
  exit 1
fi

app_id="tcr"
app_version=$(go run . --version | cut -f3 -d' ')
echo "Building TCR version ${app_version} for all targets"

fyne-cross darwin --app-id ${app_id} --app-version "${app_version}" -arch=*
fyne-cross linux --app-id ${app_id} --app-version "${app_version}" -arch=*
fyne-cross windows --app-id ${app_id} --app-version "${app_version}" -arch=*
