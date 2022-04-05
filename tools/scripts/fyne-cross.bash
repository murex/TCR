#!/usr/bin/env bash
#
# Copyright (c) 2021 Murex
#
# Permission is hereby granted, free of charge, to any person obtaining a copy
# of this software and associated documentation files (the "Software"), to deal
# in the Software without restriction, including without limitation the rights
# to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
# copies of the Software, and to permit persons to whom the Software is
# furnished to do so, subject to the following conditions:
#
# The above copyright notice and this permission notice shall be included in all
# copies or substantial portions of the Software.
#
# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
# FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
# AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
# LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
# OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
# SOFTWARE.

set -eu

base_dir="$(cd "$(dirname -- "$0")" && pwd)/../tcr-gui"

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

#fyne-cross darwin --app-id ${app_id} --app-version "${app_version}" -arch=amd64
fyne-cross linux --app-id ${app_id} --app-version "${app_version}" -arch=amd64
fyne-cross windows --app-id ${app_id} --app-version "${app_version}" -arch=amd64

cd fyne-cross/dist

#tar -zcvf ${app_id}_"${app_version}"_Darwin_x86_64.tar.gz darwin-amd64
tar -zcvf ${app_id}_"${app_version}"_Linux_x86_64.tar.gz linux-amd64
tar -zcvf ${app_id}_"${app_version}"_Windows_x86_64.tar.gz windows-amd64

#cd ..

#mkdir -p release
#mv ./dist/*.tar.gz release


