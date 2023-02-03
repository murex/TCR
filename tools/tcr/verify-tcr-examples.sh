#!/usr/bin/env bash
#
# Copyright (c) 2022 Murex
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

repo_root_dir="$(git rev-parse --show-toplevel)"
src_dir="$(cd "${repo_root_dir}/src" && pwd)"
example_dir="$(cd "${repo_root_dir}/examples" && pwd)"


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

for dir in "${example_dir}"/*; do
  if [ -f "${dir}/tcrw" ]; then

    example_name=$(basename ${dir})
    echo "Checking Example ${example_name}"

    if [ -f "${dir}/cmake_easy_setup.sh" ]; then
      cd "${dir}" && ./cmake_easy_setup.sh
    fi

    cd "${src_dir}" && go run . "$@" \
      --base-dir="${dir}" \
      --config-dir="${dir}" \
      --work-dir="${dir}" \
      check

    if [ $? -ne 0 ]; then
      echo "TCR check on ${example_name} failed. \nAborting..."
      exit 1
    fi
  fi
done

exit 0
