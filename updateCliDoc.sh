#!/usr/bin/env bash

set -eu

BASE_DIR="$(cd "$(dirname -- "$0")" && pwd)"

cd "${BASE_DIR}"/doc
go run .