#!/usr/bin/env bash
#
# Copyright (c) 2023 Murex
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

set -u

command_args=$*

repo_root_dir="$(git rev-parse --show-toplevel)"
src_dir="$(cd "${repo_root_dir}/src" && pwd)"
testdata_java_dir="/d/Perforce/tcr/dev/dxp/tcr/testdata/java/"

TCR_BASE_DIR="${testdata_java_dir}"
export TCR_BASE_DIR
TCR_WORK_DIR="${testdata_java_dir}"
export TCR_WORK_DIR
TCR_CONFIG_DIR="${testdata_java_dir}"
export TCR_CONFIG_DIR

# ------------------------------------------------------------------------------
# Trace messages
# ------------------------------------------------------------------------------

trace_info() {
  message="$1"
  echo >&2 "[TCR] ${message}"
}

trace_error() {
  message="$1"
  echo >&2 "[TCR] ERROR: ${message}"
}

# ------------------------------------------------------------------------------
# Main
# ------------------------------------------------------------------------------

cd "${src_dir}" && go run . $command_args \
  --base-dir="${TCR_BASE_DIR}" \
  --config-dir="${TCR_CONFIG_DIR}" \
  --work-dir="${TCR_WORK_DIR}" \
  --vcs=p4
exit $?
