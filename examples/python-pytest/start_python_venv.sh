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

# ------------------------------------------------------------------------------
# For traces, warnings and errors
# ------------------------------------------------------------------------------

print_info() {
  message="$1"
  printf "%b" "${message}\n" | while IFS= read -r line; do printf "%b" "\e[1;34m>>> ${line} \e[0m\n"; done
}

print_warning() {
  message="$1"
  printf "%b" "${message}\n" | while IFS= read -r line; do printf "%b" "\e[1;33m>>> ${line} \e[0m\n"; done
}

print_error() {
  message="$1"
  printf "%b" "${message}\n" | while IFS= read -r line; do printf "%b" "\e[1;31m>>> ${line} \e[0m\n"; done
}

print_horizontal_line() {
  term_columns=$(tput cols)
  repeated=$((term_columns - 5))
  line=$(head -c "${repeated}" </dev/zero | tr '\0' '-')
  print_info "$line"
}

# ------------------------------------------------------------------------------
# Check if a command is available on the machine path
# ------------------------------------------------------------------------------

check_command_availability() {
  command_name="$1"
  if type "${command_name}" >/dev/null 2>/dev/null; then
    print_info "Found ${command_name} at $(get_command_path "${command_name}")"
    return 0
  else
    print_warning "Command ${command_name} not found"
    return 1
  fi
}

# ------------------------------------------------------------------------------
# Check if a command runs as expected
# ------------------------------------------------------------------------------

check_command_execution() {
  command_name="$1"
  if "${command_name}" --version >/dev/null 2>/dev/null; then
    print_info "Checking ${command_name} ==> ok"
    return 0
  else
    print_warning "Checking ${command_name} ==> failed"
    return 1
  fi
}

is_python3_command_available() {
  check_command_availability python3 && check_command_execution python3
}

is_python_command_available() {
  check_command_availability python && check_command_execution python
}

get_command_path() {
  command_name="$1"
  which "${command_name}"
}

get_python3_path() {
  get_command_path python3
}

get_python_path() {
  get_command_path python
}

get_python_version() {
  "${PYTHON_PATH}" --version | awk '{ print $2 }'
}

locate_python() {
  print_info "Locating Python executable..."
  if is_python3_command_available; then
    PYTHON_PATH=$(get_python3_path)
  elif is_python_command_available; then
    PYTHON_PATH=$(get_python_path)
  else
    python_url="https://www.python.org/downloads/"
    print_error "No usable Python executable found in path"
    print_error "You may update your path if Python is already installed on your machine"
    print_error "For a new installation, please refer to ${python_url}"
    exit 1
  fi

  print_info "Python version is $(get_python_version)"
}

# ------------------------------------------------------------------------------
# Entry point
# ------------------------------------------------------------------------------

venv_dir="venv"
base_dir="$(cd "$(dirname -- "$0")" && pwd)"
venv_path="${base_dir}/${venv_dir}"

if ! [ -d "${venv_path}" ]; then
  locate_python
  print_horizontal_line
  print_info "Creating python virtual environment..."
  "${PYTHON_PATH}" -m venv "${venv_path}"
fi

print_info "Starting python virtual environment..."
# instead of relying on venv's activate which sometimes screws up the path on windows,
# we append venv/Scripts to the path by ourselves
# source "${venv_dir}/Scripts/activate"
VIRTUAL_ENV="${venv_path}"
export VIRTUAL_ENV
PATH="$VIRTUAL_ENV/Scripts:$VIRTUAL_ENV/bin:$PATH"
export PATH

print_info "Upgrading pip..."
python -m pip install --upgrade pip
print_info "Adding kata module..."
pip install --editable .
print_info "Adding kata dependencies..."
pip install --use-pep517 -r ./requirements.txt

# Starting a new shell in order to keep the changes done on the path after the script ends
print_horizontal_line
print_info "Python virtual environment is ready (type 'exit' to leave)"
print_horizontal_line
exec bash --norc
