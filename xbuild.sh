#!/usr/bin/env bash

set -eu

base_dir="$(cd "$(dirname -- "$0")" && pwd)"
cd "${base_dir}"

package=.
package_name="tcr"
platforms=("windows/amd64" "darwin/amd64" "linux/amd64")
output_dir="build"

mkdir -p ${output_dir}

for platform in "${platforms[@]}"; do
  platform_split=(${platform//\// })
  GOOS=${platform_split[0]}
  GOARCH=${platform_split[1]}
  output_name=${output_dir}/${package_name}'-'${GOOS}'-'${GOARCH}
  if [ "${GOOS}" = "windows" ]; then
    output_name+='.exe'
  fi

  echo "- Building ${package_name} for ${platform}"
  env GOOS="${GOOS}" GOARCH="${GOARCH}" go build -o ${output_name} ${package}
  if [ $? -ne 0 ]; then
    echo 'An error has occurred! Aborting the script execution...'
    exit 1
  fi
  chmod +x ${output_name}
done
