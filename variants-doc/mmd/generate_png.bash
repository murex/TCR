#!/usr/bin/env bash

echo "Installing mermaid-cli"
npm update -g @mermaid-js/mermaid-cli || echo "Failed to install mermaid-cli. Aborting"

repo_root_dir="$(git rev-parse --show-toplevel)"
webapp_images_dir="${repo_root_dir}/webapp/src/assets/images"

echo "Generating files into ${webapp_images_dir}"
for mmd_file in ./*.mmd; do
    png_file="${mmd_file%.mmd}.png"
    echo "- Generating ${png_file} from ${mmd_file}"
    mmdc --input "${mmd_file}" --output "${webapp_images_dir}/${png_file}" --theme dark --backgroundColor transparent
done
