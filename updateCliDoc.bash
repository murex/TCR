#!/usr/bin/env bash

set -eu

base_dir="$(cd "$(dirname -- "$0")" && pwd)"

cd "${base_dir}"/doc

go mod tidy
go run .
