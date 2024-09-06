#!/usr/bin/env bash

echo "Installing mermaid-cli"
npm update -g @mermaid-js/mermaid-cli || echo "Failed to install mermaid-cli. Aborting"

for mmd_file in ./*.mmd; do
    png_file="${mmd_file%.mmd}.png"
    echo "- Generating ${png_file} from ${mmd_file}"
    mmdc --input "${mmd_file}" --output "${png_file}" --theme dark --backgroundColor transparent
done
