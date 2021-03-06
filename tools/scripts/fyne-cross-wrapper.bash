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
set -x

base_dir="$(cd "$(dirname -- "$0")" && pwd)/../tcr-gui"

cd "${base_dir}"

if ! type go >/dev/null 2>/dev/null; then
  echo >&2 "go is required to run this script"
  echo >&2 "Refer to https://golang.org/ for installation and setup"
  exit 1
fi

if ! type fyne-cross >/dev/null 2>/dev/null; then
  echo >&2 "fyne-cross command is required to run this script"
  echo >&2 "Refer to https://github.com/fyne-io/fyne-cross for installation and setup"
  exit 1
fi

# --------------------------------------------------------------------------------------
# Converts go command arguments (sent by goreleaser) to fyne-cross compatible ones,
# and calls fyne-cross with them.
# (goreleaser is assuming it's calling a go command)
# --------------------------------------------------------------------------------------

go_cmd_to_fyne_cross_wrapper() {
  REMAINDER=()
  LD_FLAGS=()
  BUILD_CALLED=0
  APP_NAME=""

  while [[ $# -gt 0 ]]; do
    key="$1"
    echo $key

    case $key in
    version)
      go version
      exit 0
      ;;
    build)
      BUILD_CALLED=1
      shift # past argument
      ;;
    -ldflags=-s)
      LD_FLAGS+=("-s")
      shift # past argument
      ;;
    -w)
      LD_FLAGS+=("$1")
      shift # past argument
      ;;
    -X)
      LD_FLAGS+=("$1")
      LD_FLAGS+=("$2")
      shift # past argument
      shift # past value
      ;;
    -o)
      APP_NAME=$(basename "$2")
      OUTPUT_DIR=$(dirname "$2")
      shift # past argument
      shift # past value
      ;;
    *) # remainder
      REMAINDER+=("$1") # save it in an array for later
      shift # past argument
      ;;
    esac
  done

  if [ $BUILD_CALLED -eq 1 ]; then
    app_id="${APP_NAME}"
    app_version=$(go run . --version | cut -f3 -d' ')

    case ${GOOS} in
    darwin)
      fyne-cross "${GOOS}" --app-id "${app_id}" --app-version "${app_version}" -arch="${GOARCH}" -ldflags "\"${LD_FLAGS[*]}\"" ${REMAINDER[@]}
      exe_file="fyne-cross/dist/${GOOS}-${GOARCH}/${app_id}.app/Contents/MacOS/${app_id}"
      ;;
    linux)
      fyne-cross "${GOOS}" --app-id "${app_id}" -arch="${GOARCH}" -ldflags "\"${LD_FLAGS[*]}\"" ${REMAINDER[@]}
      #fyne-cross "${GOOS}" --app-id "${app_id}" --app-version "${app_version}" -arch="${GOARCH}" -ldflags "\"${LD_FLAGS[*]}\"" ${REMAINDER[@]}
      (cd fyne-cross/dist/${GOOS}-${GOARCH} && tar xvf "${app_id}".tar.gz)
      exe_file="fyne-cross/dist/${GOOS}-${GOARCH}/usr/local/bin/${app_id}"
      ;;
    windows)
      fyne-cross "${GOOS}" --app-id "${app_id}" --app-version "${app_version}" -arch="${GOARCH}" -ldflags "\"${LD_FLAGS[*]}\"" ${REMAINDER[@]}
      #(cd fyne-cross/dist/${GOOS}-${GOARCH} && tar xvf "${app_id}".zip)
      (cd fyne-cross/dist/${GOOS}-${GOARCH} && unzip -q -o "${app_id}".zip)
      exe_file="fyne-cross/dist/${GOOS}-${GOARCH}/${app_id}"
      ;;
    *)
      echo "OS not supported: ${GOOS}"
      exit 1
      ;;
    esac

    # Move the built exe file from fyne-cross output dir to goreleaser directory
    mkdir -p "${OUTPUT_DIR}"
    mv "${exe_file}" "${OUTPUT_DIR}/${APP_NAME}"
  fi
}

#rm "$0".log
go_cmd_to_fyne_cross_wrapper $@ # >>"$0".log
exit $?
